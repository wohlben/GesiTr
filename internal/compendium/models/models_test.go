package models

import (
	"testing"
	"time"
)

// --- Equipment ---

func TestEquipmentEntityTableName(t *testing.T) {
	if got := (EquipmentEntity{}).TableName(); got != "equipment" {
		t.Errorf("TableName() = %q, want %q", got, "equipment")
	}
}

func TestEquipmentEntityToDTO(t *testing.T) {
	now := time.Now()
	imgUrl := "http://example.com/img.png"
	e := &EquipmentEntity{
		BaseModel:   BaseModel{ID: 1, CreatedAt: now, UpdatedAt: now},
		Name:        "barbell",
		DisplayName: "Barbell",
		Description: "A long bar",
		Category:    EquipmentCategoryFreeWeights,
		ImageUrl:    &imgUrl,
		TemplateID:  "barbell",
		CreatedBy:   "system",
		Version:     2,
	}
	dto := e.ToDTO()
	if dto.ID != 1 {
		t.Errorf("ID = %d, want 1", dto.ID)
	}
	if dto.Name != "barbell" {
		t.Errorf("Name = %q", dto.Name)
	}
	if dto.DisplayName != "Barbell" {
		t.Errorf("DisplayName = %q", dto.DisplayName)
	}
	if dto.Description != "A long bar" {
		t.Errorf("Description = %q", dto.Description)
	}
	if dto.Category != EquipmentCategoryFreeWeights {
		t.Errorf("Category = %q", dto.Category)
	}
	if *dto.ImageUrl != imgUrl {
		t.Errorf("ImageUrl = %q", *dto.ImageUrl)
	}
	if dto.TemplateID != "barbell" {
		t.Errorf("TemplateID = %q", dto.TemplateID)
	}
	if dto.CreatedBy != "system" {
		t.Errorf("CreatedBy = %q", dto.CreatedBy)
	}
	if dto.Version != 2 {
		t.Errorf("Version = %d", dto.Version)
	}
}

func TestEquipmentFromDTO(t *testing.T) {
	imgUrl := "http://example.com/img.png"
	dto := Equipment{
		BaseModel:   BaseModel{ID: 3},
		Name:        "bench",
		DisplayName: "Bench",
		Description: "A flat bench",
		Category:    EquipmentCategoryBenches,
		ImageUrl:    &imgUrl,
		TemplateID:  "bench",
		CreatedBy:   "user",
		Version:     1,
	}
	e := EquipmentFromDTO(dto)
	if e.ID != 3 || e.Name != "bench" || e.Category != EquipmentCategoryBenches || e.Version != 1 {
		t.Error("EquipmentFromDTO field mismatch")
	}
	if e.DisplayName != "Bench" || e.Description != "A flat bench" || e.TemplateID != "bench" || e.CreatedBy != "user" {
		t.Error("EquipmentFromDTO field mismatch")
	}
	if *e.ImageUrl != imgUrl {
		t.Error("EquipmentFromDTO ImageUrl mismatch")
	}
}

// --- Fulfillment ---

func TestFulfillmentEntityTableName(t *testing.T) {
	if got := (FulfillmentEntity{}).TableName(); got != "fulfillments" {
		t.Errorf("TableName() = %q, want %q", got, "fulfillments")
	}
}

func TestFulfillmentEntityToDTO(t *testing.T) {
	e := &FulfillmentEntity{
		BaseModel:                   BaseModel{ID: 1},
		EquipmentTemplateID:         "eq1",
		FulfillsEquipmentTemplateID: "eq2",
		CreatedBy:                   "system",
	}
	dto := e.ToDTO()
	if dto.ID != 1 || dto.EquipmentTemplateID != "eq1" || dto.FulfillsEquipmentTemplateID != "eq2" || dto.CreatedBy != "system" {
		t.Error("ToDTO field mismatch")
	}
}

func TestFulfillmentFromDTO(t *testing.T) {
	dto := Fulfillment{
		BaseModel:                   BaseModel{ID: 2},
		EquipmentTemplateID:         "a",
		FulfillsEquipmentTemplateID: "b",
		CreatedBy:                   "user",
	}
	e := FulfillmentFromDTO(dto)
	if e.ID != 2 || e.EquipmentTemplateID != "a" || e.FulfillsEquipmentTemplateID != "b" || e.CreatedBy != "user" {
		t.Error("FulfillmentFromDTO field mismatch")
	}
}

// --- ExerciseRelationship ---

func TestExerciseRelationshipEntityTableName(t *testing.T) {
	if got := (ExerciseRelationshipEntity{}).TableName(); got != "exercise_relationships" {
		t.Errorf("TableName() = %q, want %q", got, "exercise_relationships")
	}
}

func TestExerciseRelationshipEntityToDTO(t *testing.T) {
	desc := "test desc"
	e := &ExerciseRelationshipEntity{
		BaseModel:              BaseModel{ID: 5},
		RelationshipType:       ExerciseRelationshipTypeSimilar,
		Strength:               0.8,
		Description:            &desc,
		CreatedBy:              "system",
		FromExerciseTemplateID: "ex1",
		ToExerciseTemplateID:   "ex2",
	}
	dto := e.ToDTO()
	if dto.ID != 5 || dto.RelationshipType != ExerciseRelationshipTypeSimilar || dto.Strength != 0.8 {
		t.Error("ToDTO field mismatch")
	}
	if *dto.Description != "test desc" || dto.CreatedBy != "system" {
		t.Error("ToDTO field mismatch")
	}
	if dto.FromExerciseTemplateID != "ex1" || dto.ToExerciseTemplateID != "ex2" {
		t.Error("ToDTO field mismatch")
	}
}

func TestExerciseRelationshipFromDTO(t *testing.T) {
	desc := "desc"
	dto := ExerciseRelationship{
		BaseModel:              BaseModel{ID: 3},
		RelationshipType:       ExerciseRelationshipTypeVariation,
		Strength:               0.5,
		Description:            &desc,
		CreatedBy:              "user",
		FromExerciseTemplateID: "a",
		ToExerciseTemplateID:   "b",
	}
	e := ExerciseRelationshipFromDTO(dto)
	if e.ID != 3 || e.RelationshipType != ExerciseRelationshipTypeVariation || e.Strength != 0.5 {
		t.Error("FromDTO field mismatch")
	}
	if *e.Description != "desc" || e.CreatedBy != "user" || e.FromExerciseTemplateID != "a" || e.ToExerciseTemplateID != "b" {
		t.Error("FromDTO field mismatch")
	}
}

// --- ExerciseGroup ---

func TestExerciseGroupEntityTableName(t *testing.T) {
	if got := (ExerciseGroupEntity{}).TableName(); got != "exercise_groups" {
		t.Errorf("TableName() = %q, want %q", got, "exercise_groups")
	}
}

func TestExerciseGroupEntityToDTO(t *testing.T) {
	desc := "group desc"
	e := &ExerciseGroupEntity{
		BaseModel:   BaseModel{ID: 10},
		TemplateID:  "g1",
		Name:        "Group One",
		Description: &desc,
		CreatedBy:   "system",
	}
	dto := e.ToDTO()
	if dto.ID != 10 || dto.TemplateID != "g1" || dto.Name != "Group One" || *dto.Description != "group desc" || dto.CreatedBy != "system" {
		t.Error("ToDTO field mismatch")
	}
}

func TestExerciseGroupFromDTO(t *testing.T) {
	desc := "d"
	dto := ExerciseGroup{
		BaseModel:   BaseModel{ID: 7},
		TemplateID:  "g2",
		Name:        "Group Two",
		Description: &desc,
		CreatedBy:   "user",
	}
	e := ExerciseGroupFromDTO(dto)
	if e.ID != 7 || e.TemplateID != "g2" || e.Name != "Group Two" || *e.Description != "d" || e.CreatedBy != "user" {
		t.Error("FromDTO field mismatch")
	}
}

// --- ExerciseGroupMember ---

func TestExerciseGroupMemberEntityTableName(t *testing.T) {
	if got := (ExerciseGroupMemberEntity{}).TableName(); got != "exercise_group_members" {
		t.Errorf("TableName() = %q, want %q", got, "exercise_group_members")
	}
}

func TestExerciseGroupMemberEntityToDTO(t *testing.T) {
	e := &ExerciseGroupMemberEntity{
		BaseModel:          BaseModel{ID: 20},
		GroupTemplateID:    "g1",
		ExerciseTemplateID: "ex1",
		AddedBy:            "user",
	}
	dto := e.ToDTO()
	if dto.ID != 20 || dto.GroupTemplateID != "g1" || dto.ExerciseTemplateID != "ex1" || dto.AddedBy != "user" {
		t.Error("ToDTO field mismatch")
	}
}

func TestExerciseGroupMemberFromDTO(t *testing.T) {
	dto := ExerciseGroupMember{
		BaseModel:          BaseModel{ID: 15},
		GroupTemplateID:    "g3",
		ExerciseTemplateID: "ex5",
		AddedBy:            "admin",
	}
	e := ExerciseGroupMemberFromDTO(dto)
	if e.ID != 15 || e.GroupTemplateID != "g3" || e.ExerciseTemplateID != "ex5" || e.AddedBy != "admin" {
		t.Error("FromDTO field mismatch")
	}
}

// --- Exercise ---

func TestExerciseEntityTableName(t *testing.T) {
	if got := (ExerciseEntity{}).TableName(); got != "exercises" {
		t.Errorf("TableName() = %q, want %q", got, "exercises")
	}
}

func TestExerciseEntityToDTOFull(t *testing.T) {
	authorName := "author"
	authorUrl := "http://author.com"
	parentID := uint(99)
	e := &ExerciseEntity{
		BaseModel:           BaseModel{ID: 1},
		Name:                "Bench Press",
		Slug:                "bench-press",
		Type:                ExerciseTypeStrength,
		TechnicalDifficulty: DifficultyIntermediate,
		BodyWeightScaling:   0.0,
		Description:         "A chest exercise",
		AuthorName:          &authorName,
		AuthorUrl:           &authorUrl,
		CreatedBy:           "system",
		Version:             3,
		ParentExerciseID:    &parentID,
		TemplateID:          "test-exercise",
		Forces:              []ExerciseForce{{ExerciseID: 1, Force: ForcePush}},
		Muscles: []ExerciseMuscle{
			{ExerciseID: 1, Muscle: MuscleChest, IsPrimary: true},
			{ExerciseID: 1, Muscle: MuscleTriceps, IsPrimary: false},
		},
		Paradigms:        []ExerciseMeasurementParadigm{{ExerciseID: 1, Paradigm: MeasurementRepBased}},
		Instructions:     []ExerciseInstruction{{ExerciseID: 1, Position: 0, Text: "Step 1"}},
		Images:           []ExerciseImage{{ExerciseID: 1, Position: 0, Path: "/img/bench.jpg"}},
		AlternativeNames: []ExerciseAlternativeName{{ExerciseID: 1, Name: "Flat Bench"}},
		Equipment:        []ExerciseEquipment{{ExerciseID: 1, EquipmentTemplateID: "barbell"}},
	}

	dto := e.ToDTO()

	if dto.ID != 1 || dto.Name != "Bench Press" || dto.Slug != "bench-press" {
		t.Error("basic fields mismatch")
	}
	if dto.Type != ExerciseTypeStrength || dto.TechnicalDifficulty != DifficultyIntermediate {
		t.Error("type/difficulty mismatch")
	}
	if *dto.AuthorName != "author" || *dto.AuthorUrl != "http://author.com" {
		t.Error("author mismatch")
	}
	if dto.Version != 3 || *dto.ParentExerciseID != 99 || dto.TemplateID != "test-exercise" {
		t.Error("version/parent/template mismatch")
	}
	if len(dto.Force) != 1 || dto.Force[0] != ForcePush {
		t.Errorf("Force = %v", dto.Force)
	}
	if len(dto.PrimaryMuscles) != 1 || dto.PrimaryMuscles[0] != MuscleChest {
		t.Errorf("PrimaryMuscles = %v", dto.PrimaryMuscles)
	}
	if len(dto.SecondaryMuscles) != 1 || dto.SecondaryMuscles[0] != MuscleTriceps {
		t.Errorf("SecondaryMuscles = %v", dto.SecondaryMuscles)
	}
	if len(dto.SuggestedMeasurementParadigms) != 1 || dto.SuggestedMeasurementParadigms[0] != MeasurementRepBased {
		t.Errorf("Paradigms = %v", dto.SuggestedMeasurementParadigms)
	}
	if len(dto.Instructions) != 1 || dto.Instructions[0] != "Step 1" {
		t.Errorf("Instructions = %v", dto.Instructions)
	}
	if len(dto.Images) != 1 || dto.Images[0] != "/img/bench.jpg" {
		t.Errorf("Images = %v", dto.Images)
	}
	if len(dto.AlternativeNames) != 1 || dto.AlternativeNames[0] != "Flat Bench" {
		t.Errorf("AlternativeNames = %v", dto.AlternativeNames)
	}
	if len(dto.EquipmentIDs) != 1 || dto.EquipmentIDs[0] != "barbell" {
		t.Errorf("EquipmentIDs = %v", dto.EquipmentIDs)
	}
}

func TestExerciseEntityToDTOEmpty(t *testing.T) {
	e := &ExerciseEntity{
		BaseModel: BaseModel{ID: 2},
		Name:      "Empty",
		Slug:      "empty",
		Type:      ExerciseTypeCardio,
		CreatedBy: "system",
	}
	dto := e.ToDTO()
	if dto.Force != nil || dto.PrimaryMuscles != nil || dto.SecondaryMuscles != nil {
		t.Error("expected nil slices for empty collections")
	}
	if dto.Instructions != nil || dto.Images != nil || dto.AlternativeNames != nil || dto.EquipmentIDs != nil {
		t.Error("expected nil slices for empty collections")
	}
}

func TestExerciseFromDTOFull(t *testing.T) {
	dto := Exercise{
		BaseModel:                     BaseModel{ID: 1},
		Name:                          "Squat",
		Slug:                          "squat",
		Type:                          ExerciseTypeStrength,
		Force:                         []Force{ForcePush},
		PrimaryMuscles:                []Muscle{MuscleQuads},
		SecondaryMuscles:              []Muscle{MuscleGlutes, MuscleHamstrings},
		TechnicalDifficulty:           DifficultyAdvanced,
		BodyWeightScaling:             1.0,
		SuggestedMeasurementParadigms: []MeasurementParadigm{MeasurementRepBased, MeasurementAMRAP},
		Description:                   "A leg exercise",
		Instructions:                  []string{"Step 1", "Step 2"},
		Images:                        []string{"/img/a.jpg", "/img/b.jpg"},
		AlternativeNames:              []string{"Back Squat"},
		CreatedBy:                     "system",
		TemplateID:                    "tid",
		EquipmentIDs:                  []string{"barbell", "rack"},
	}

	e := ExerciseFromDTO(dto)

	if e.Name != "Squat" || e.Slug != "squat" || e.Type != ExerciseTypeStrength {
		t.Error("basic fields mismatch")
	}
	if len(e.Forces) != 1 || e.Forces[0].Force != ForcePush {
		t.Errorf("Forces = %v", e.Forces)
	}
	if len(e.Muscles) != 3 {
		t.Fatalf("Muscles len = %d, want 3", len(e.Muscles))
	}
	primaryCount := 0
	for _, m := range e.Muscles {
		if m.IsPrimary {
			primaryCount++
		}
	}
	if primaryCount != 1 {
		t.Errorf("primary muscle count = %d, want 1", primaryCount)
	}
	if len(e.Paradigms) != 2 {
		t.Errorf("Paradigms len = %d", len(e.Paradigms))
	}
	if len(e.Instructions) != 2 || e.Instructions[0].Position != 0 || e.Instructions[1].Position != 1 {
		t.Errorf("Instructions = %v", e.Instructions)
	}
	if len(e.Images) != 2 || e.Images[0].Position != 0 || e.Images[1].Position != 1 {
		t.Errorf("Images = %v", e.Images)
	}
	if len(e.AlternativeNames) != 1 || e.AlternativeNames[0].Name != "Back Squat" {
		t.Errorf("AlternativeNames = %v", e.AlternativeNames)
	}
	if len(e.Equipment) != 2 || e.Equipment[0].EquipmentTemplateID != "barbell" {
		t.Errorf("Equipment = %v", e.Equipment)
	}
}

func TestExerciseFromDTOEmpty(t *testing.T) {
	dto := Exercise{
		Name:      "Empty",
		Slug:      "empty",
		Type:      ExerciseTypeCardio,
		CreatedBy: "system",
	}
	e := ExerciseFromDTO(dto)
	if e.Forces != nil || e.Muscles != nil || e.Paradigms != nil {
		t.Error("expected nil slices")
	}
	if e.Instructions != nil || e.Images != nil || e.AlternativeNames != nil || e.Equipment != nil {
		t.Error("expected nil slices")
	}
}
