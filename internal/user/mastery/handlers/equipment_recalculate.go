package handlers

import (
	equipmentmodels "gesitr/internal/equipment/models"
	"gesitr/internal/user/mastery/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// RecalculateEquipmentContributions recomputes equipment_mastery_contributions rows
// for the given owner and affected equipment IDs.
// Called after equipment relationship or fulfillment create/delete.
func RecalculateEquipmentContributions(db *gorm.DB, owner string, equipmentIDs ...uint) error {
	// 1. Fetch equipment relationships for this owner involving affected IDs.
	var relationships []equipmentmodels.EquipmentRelationshipEntity
	if err := db.Where("owner = ? AND (from_equipment_id IN ? OR to_equipment_id IN ?)", owner, equipmentIDs, equipmentIDs).
		Find(&relationships).Error; err != nil {
		return err
	}

	// 2. Fetch fulfillments for this owner involving affected IDs.
	var fulfillments []equipmentmodels.FulfillmentEntity
	if err := db.Where("owner = ? AND (equipment_id IN ? OR fulfills_equipment_id IN ?)", owner, equipmentIDs, equipmentIDs).
		Find(&fulfillments).Error; err != nil {
		return err
	}

	// 3. Build contribution rows, keeping the best multiplier per pair.
	type contribKey struct {
		equipmentID     uint
		contributesFrom uint
	}
	best := make(map[contribKey]models.EquipmentMasteryContributionEntity)

	// From equipment relationships (bidirectional).
	for _, rel := range relationships {
		mult, ok := models.ComputeEquipmentContributionMultiplier(rel.Strength, string(rel.RelationshipType))
		if !ok {
			continue
		}

		pairs := [][2]uint{
			{rel.FromEquipmentID, rel.ToEquipmentID},
			{rel.ToEquipmentID, rel.FromEquipmentID},
		}
		for _, pair := range pairs {
			key := contribKey{pair[0], pair[1]}
			if existing, exists := best[key]; !exists || mult > existing.Multiplier {
				best[key] = models.EquipmentMasteryContributionEntity{
					Owner:             owner,
					EquipmentID:       pair[0],
					ContributesFromID: pair[1],
					Multiplier:        mult,
					RelationshipType:  string(rel.RelationshipType),
				}
			}
		}
	}

	// From fulfillments (bidirectional).
	fulfillmentMult := models.FulfillmentContributionMultiplier()
	for _, ful := range fulfillments {
		pairs := [][2]uint{
			{ful.EquipmentID, ful.FulfillsEquipmentID},
			{ful.FulfillsEquipmentID, ful.EquipmentID},
		}
		for _, pair := range pairs {
			key := contribKey{pair[0], pair[1]}
			if existing, exists := best[key]; !exists || fulfillmentMult > existing.Multiplier {
				best[key] = models.EquipmentMasteryContributionEntity{
					Owner:             owner,
					EquipmentID:       pair[0],
					ContributesFromID: pair[1],
					Multiplier:        fulfillmentMult,
					RelationshipType:  "fulfillment",
				}
			}
		}
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// Delete existing contributions for affected equipment.
		if err := tx.Where("owner = ? AND (equipment_id IN ? OR contributes_from_id IN ?)", owner, equipmentIDs, equipmentIDs).
			Delete(&models.EquipmentMasteryContributionEntity{}).Error; err != nil {
			return err
		}

		if len(best) == 0 {
			return nil
		}
		rows := make([]models.EquipmentMasteryContributionEntity, 0, len(best))
		for _, row := range best {
			rows = append(rows, row)
		}
		return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&rows).Error
	})
}

// BackfillEquipmentContributions computes equipment_mastery_contributions for all existing
// equipment relationships and fulfillments. Intended to be called once during migration.
func BackfillEquipmentContributions(db *gorm.DB) error {
	// Collect distinct owners from equipment relationships and fulfillments.
	ownerSet := make(map[string]bool)

	var relOwners []string
	db.Model(&equipmentmodels.EquipmentRelationshipEntity{}).
		Distinct("owner").Pluck("owner", &relOwners)
	for _, o := range relOwners {
		ownerSet[o] = true
	}

	var fulOwners []string
	db.Model(&equipmentmodels.FulfillmentEntity{}).
		Distinct("owner").Pluck("owner", &fulOwners)
	for _, o := range fulOwners {
		ownerSet[o] = true
	}

	for owner := range ownerSet {
		// Collect all equipment IDs involved for this owner.
		var eqIDs []uint
		seen := make(map[uint]bool)

		var fromIDs, toIDs []uint
		db.Model(&equipmentmodels.EquipmentRelationshipEntity{}).
			Where("owner = ?", owner).Pluck("from_equipment_id", &fromIDs)
		db.Model(&equipmentmodels.EquipmentRelationshipEntity{}).
			Where("owner = ?", owner).Pluck("to_equipment_id", &toIDs)

		var fulEqIDs, fulFulIDs []uint
		db.Model(&equipmentmodels.FulfillmentEntity{}).
			Where("owner = ?", owner).Pluck("equipment_id", &fulEqIDs)
		db.Model(&equipmentmodels.FulfillmentEntity{}).
			Where("owner = ?", owner).Pluck("fulfills_equipment_id", &fulFulIDs)

		for _, id := range append(append(append(fromIDs, toIDs...), fulEqIDs...), fulFulIDs...) {
			if !seen[id] {
				seen[id] = true
				eqIDs = append(eqIDs, id)
			}
		}

		if err := RecalculateEquipmentContributions(db, owner, eqIDs...); err != nil {
			return err
		}
	}
	return nil
}
