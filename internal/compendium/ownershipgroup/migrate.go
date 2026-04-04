package ownershipgroup

import (
	"log"

	"gesitr/internal/compendium/ownershipgroup/models"

	"gorm.io/gorm"
)

var topLevelTables = []string{
	"exercises",
	"equipment",
	"workouts",
	"workout_groups",
	"localities",
}

type subEntityTable struct {
	table       string
	parentTable string
	parentFK    string
}

var subEntityTables = []subEntityTable{
	{table: "exercise_relationships", parentTable: "exercises", parentFK: "from_exercise_id"},
	{table: "equipment_relationships", parentTable: "equipment", parentFK: "from_equipment_id"},
	{table: "fulfillments", parentTable: "equipment", parentFK: "equipment_id"},
	{table: "exercise_groups", parentTable: "exercises", parentFK: "id"}, // standalone — treated as top-level
	{table: "workout_relationships", parentTable: "workouts", parentFK: "from_workout_id"},
	{table: "locality_availabilities", parentTable: "localities", parentFK: "locality_id"},
}

// allAffectedTables is every table that formerly had an owner column.
var allAffectedTables = []string{
	"exercises", "equipment", "workouts", "workout_groups", "localities",
	"exercise_relationships", "equipment_relationships", "fulfillments",
	"exercise_groups", "exercise_group_members",
	"workout_relationships", "locality_availabilities",
}

// MigrateExistingOwners creates ownership groups for existing entities that still
// have an owner column but no ownership_group_id. Idempotent. Also drops old
// unique indexes and the owner column once migration is complete.
func MigrateExistingOwners(db *gorm.DB) error {
	// Drop old unique indexes that included 'owner'.
	for _, idx := range []struct {
		table string
		index string
	}{
		{"exercise_relationships", "idx_exercise_relationship"},
		{"equipment_relationships", "idx_equipment_relationship"},
		{"workout_relationships", "idx_workout_relationship"},
	} {
		if db.Migrator().HasIndex(idx.table, idx.index) {
			if err := db.Migrator().DropIndex(idx.table, idx.index); err != nil {
				log.Printf("warning: could not drop index %s on %s: %v", idx.index, idx.table, err)
			}
		}
	}

	// Only migrate if the owner column still exists on a representative table.
	if !db.Migrator().HasColumn("exercises", "owner") {
		return nil
	}

	// Migrate top-level entities.
	for _, table := range topLevelTables {
		if err := migrateTopLevel(db, table); err != nil {
			return err
		}
	}

	// Migrate sub-entities.
	for _, se := range subEntityTables {
		if se.table == "exercise_groups" {
			if err := migrateTopLevel(db, se.table); err != nil {
				return err
			}
			if err := migrateSubEntity(db, "exercise_group_members", "exercise_groups", "group_id"); err != nil {
				return err
			}
			continue
		}
		if err := migrateSubEntity(db, se.table, se.parentTable, se.parentFK); err != nil {
			return err
		}
	}

	// Drop the owner column from all affected tables.
	// Use raw SQL because GORM's SQLite DropColumn requires a registered model
	// (struct with the column), but the Owner field has already been removed.
	for _, table := range allAffectedTables {
		if !db.Migrator().HasColumn(table, "owner") {
			continue
		}
		// Drop any index referencing owner before dropping the column.
		idx := "idx_" + table + "_owner"
		if db.Migrator().HasIndex(table, idx) {
			if err := db.Exec("DROP INDEX IF EXISTS " + idx).Error; err != nil {
				log.Printf("warning: could not drop index %s: %v", idx, err)
			}
		}
		if err := db.Exec("ALTER TABLE " + table + " DROP COLUMN owner").Error; err != nil {
			log.Printf("warning: could not drop owner column from %s: %v", table, err)
		} else {
			log.Printf("migrate: dropped owner column from %s", table)
		}
	}

	return nil
}

// migrateTopLevel creates an ownership group for each row that has owner set but no ownership_group_id.
func migrateTopLevel(db *gorm.DB, table string) error {
	type row struct {
		ID    uint
		Owner string
	}

	var rows []row
	if err := db.Table(table).
		Select("id, owner").
		Where("owner != '' AND (ownership_group_id IS NULL OR ownership_group_id = 0)").
		Where("deleted_at IS NULL").
		Find(&rows).Error; err != nil {
		return err
	}

	if len(rows) == 0 {
		return nil
	}

	log.Printf("migrate: creating ownership groups for %d rows in %s", len(rows), table)

	return db.Transaction(func(tx *gorm.DB) error {
		for _, r := range rows {
			group := models.OwnershipGroupEntity{}
			if err := tx.Create(&group).Error; err != nil {
				return err
			}
			membership := models.OwnershipGroupMembershipEntity{
				GroupID: group.ID,
				UserID:  r.Owner,
				Role:    models.RoleOwner,
			}
			if err := tx.Create(&membership).Error; err != nil {
				return err
			}
			if err := tx.Table(table).Where("id = ?", r.ID).Update("ownership_group_id", group.ID).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// migrateSubEntity copies the parent entity's ownership_group_id to sub-entity rows.
func migrateSubEntity(db *gorm.DB, table, parentTable, parentFK string) error {
	result := db.Exec(`
		UPDATE ` + table + ` SET ownership_group_id = (
			SELECT ` + parentTable + `.ownership_group_id FROM ` + parentTable + `
			WHERE ` + parentTable + `.id = ` + table + `.` + parentFK + `
		)
		WHERE (ownership_group_id IS NULL OR ownership_group_id = 0)
		AND deleted_at IS NULL
	`)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		log.Printf("migrate: copied ownership_group_id for %d rows in %s from %s", result.RowsAffected, table, parentTable)
	}
	return nil
}
