package exercise

import (
	"testing"

	"gesitr/internal/shared"
)

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
		BaseModel:           shared.BaseModel{ID: 1},
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
		BaseModel: shared.BaseModel{ID: 2},
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
		BaseModel:                     shared.BaseModel{ID: 1},
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
