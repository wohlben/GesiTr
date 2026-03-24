package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	equipmenthandlers "gesitr/internal/equipment/handlers"
	equipmentmodels "gesitr/internal/equipment/models"
	"gesitr/internal/exercise/models"
	"gesitr/internal/shared"

	profilemodels "gesitr/internal/profile/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupExampleDB() {
	os.Setenv("AUTH_FALLBACK_USER", "testuser")
	db, err := gorm.Open(sqlite.Open("file::memory:?_foreign_keys=on"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(
		&profilemodels.UserProfileEntity{},
		&models.ExerciseEntity{},
		&models.ExerciseForce{},
		&models.ExerciseMuscle{},
		&models.ExerciseMeasurementParadigm{},
		&models.ExerciseInstruction{},
		&models.ExerciseImage{},
		&models.ExerciseAlternativeName{},
		&models.ExerciseEquipment{},
		&models.ExerciseHistoryEntity{},
		&models.ExerciseSchemeEntity{},
		&equipmentmodels.EquipmentEntity{},
		&equipmentmodels.EquipmentHistoryEntity{},
	)
	db.Create(&profilemodels.UserProfileEntity{ID: "testuser", Name: "testuser"})
	db.Create(&profilemodels.UserProfileEntity{ID: "other", Name: "other"})
	database.DB = db
}

// newExampleRouter registers both exercise and equipment routes, so examples
// can demonstrate cross-API flows.
func newExampleRouter() *gin.Engine {
	r := newRouter()
	api := r.Group("/api")
	api.Use(auth.UserID())

	equipment := api.Group("/equipment")
	equipment.POST("", equipmenthandlers.CreateEquipment)

	return r
}

func doRawAs(r *gin.Engine, method, path, body, userID string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", userID)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// Owner sees their own exercises (both public and private) in the list.
func ExampleListExercises_owner() {
	setupExampleDB()
	r := newRouter()

	// Create a private exercise.
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Secret Move",
		"templateId": "secret-move",
		"type": "STRENGTH",
		"technicalDifficulty": "advanced",
		"bodyWeightScaling": 0,
		"description": "A private exercise",
		"version": 0
	}`)

	// Create a public exercise.
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Push-up",
		"templateId": "push-up",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 1.0,
		"description": "Bodyweight push-up",
		"version": 0,
		"public": true
	}`)

	w := doJSON(r, "GET", "/api/exercises", nil)

	var page paginatedJSON
	json.Unmarshal(w.Body.Bytes(), &page)
	fmt.Println(w.Code)
	fmt.Println(page.Total)
	// Output:
	// 200
	// 2
}

// Non-owner sees only public exercises in the list.
func ExampleListExercises_nonOwnerPublic() {
	setupExampleDB()
	r := newRouter()

	// Create a public exercise.
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Push-up",
		"templateId": "push-up",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 1.0,
		"description": "Bodyweight push-up",
		"version": 0,
		"public": true
	}`)

	// Another user sees the public exercise.
	w := doRawAs(r, "GET", "/api/exercises", "", "other")

	var page paginatedJSON
	json.Unmarshal(w.Body.Bytes(), &page)
	var exercises []models.Exercise
	json.Unmarshal(page.Items, &exercises)
	fmt.Println(w.Code)
	fmt.Println(page.Total)
	fmt.Println(exercises[0].Name)
	// Output:
	// 200
	// 1
	// Push-up
}

// Non-owner does not see private exercises in the list.
func ExampleListExercises_nonOwnerPrivate() {
	setupExampleDB()
	r := newRouter()

	// Create a private exercise.
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Secret Move",
		"templateId": "secret-move",
		"type": "STRENGTH",
		"technicalDifficulty": "advanced",
		"bodyWeightScaling": 0,
		"description": "A private exercise",
		"version": 0
	}`)

	// Another user sees an empty list.
	w := doRawAs(r, "GET", "/api/exercises", "", "other")

	var page paginatedJSON
	json.Unmarshal(w.Body.Bytes(), &page)
	fmt.Println(w.Code)
	fmt.Println(page.Total)
	// Output:
	// 200
	// 0
}

// Owner can update their exercise. Returns the updated exercise.
func ExampleUpdateExercise_owner() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/exercises", `{
		"name": "Squat",
		"templateId": "squat",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0.5,
		"description": "Barbell squat",
		"version": 0
	}`)

	w := doRaw(r, "PUT", "/api/exercises/1", `{
		"name": "Back Squat",
		"templateId": "squat",
		"type": "STRENGTH",
		"technicalDifficulty": "intermediate",
		"bodyWeightScaling": 0.5,
		"description": "Barbell back squat",
		"force": ["PUSH"],
		"primaryMuscles": ["QUADS"]
	}`)

	var exercise models.Exercise
	json.Unmarshal(w.Body.Bytes(), &exercise)
	fmt.Println(w.Code)
	fmt.Println(exercise.Name)
	fmt.Println(exercise.Owner)
	// Output:
	// 200
	// Back Squat
	// testuser
}

// Each update bumps the version and creates a history snapshot. Version 0
// preserves the original state, version 1 the updated state.
func ExampleUpdateExercise_versioning() {
	setupExampleDB()
	r := newRouter()

	// Create an exercise (version 0).
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Squat",
		"templateId": "squat",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0.5,
		"description": "Barbell squat",
		"version": 0
	}`)

	// Update it — version bumps to 1.
	w := doRaw(r, "PUT", "/api/exercises/1", `{
		"name": "Back Squat",
		"templateId": "squat",
		"type": "STRENGTH",
		"technicalDifficulty": "intermediate",
		"bodyWeightScaling": 0.5,
		"description": "Barbell back squat",
		"force": ["PUSH"],
		"primaryMuscles": ["QUADS"]
	}`)

	var exercise models.Exercise
	json.Unmarshal(w.Body.Bytes(), &exercise)
	fmt.Println("current version:", exercise.Version)

	// Version 0 preserves the original name.
	wv0 := doJSON(r, "GET", "/api/exercises/templates/squat/versions/0", nil)
	var v0 shared.VersionEntry
	json.Unmarshal(wv0.Body.Bytes(), &v0)
	var snap0 models.Exercise
	json.Unmarshal(v0.Snapshot, &snap0)
	fmt.Println("v0 name:", snap0.Name)

	// Version 1 has the updated name.
	wv1 := doJSON(r, "GET", "/api/exercises/templates/squat/versions/1", nil)
	var v1 shared.VersionEntry
	json.Unmarshal(wv1.Body.Bytes(), &v1)
	var snap1 models.Exercise
	json.Unmarshal(v1.Snapshot, &snap1)
	fmt.Println("v1 name:", snap1.Name)

	// Output:
	// current version: 1
	// v0 name: Squat
	// v1 name: Back Squat
}

// Non-owner cannot update a public exercise — returns 403.
func ExampleUpdateExercise_nonOwnerPublic() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/exercises", `{
		"name": "Push-up",
		"templateId": "push-up",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 1.0,
		"description": "Bodyweight push-up",
		"version": 0,
		"public": true
	}`)

	w := doRawAs(r, "PUT", "/api/exercises/1", `{
		"name": "Modified Push-up",
		"templateId": "push-up",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 1.0,
		"description": "Tampered"
	}`, "other")
	fmt.Println(w.Code)
	// Output: 403
}

// Non-owner cannot update a private exercise — returns 403.
func ExampleUpdateExercise_nonOwnerPrivate() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/exercises", `{
		"name": "Secret Move",
		"templateId": "secret-move",
		"type": "STRENGTH",
		"technicalDifficulty": "advanced",
		"bodyWeightScaling": 0,
		"description": "A private exercise",
		"version": 0
	}`)

	w := doRawAs(r, "PUT", "/api/exercises/1", `{
		"name": "Stolen Move",
		"templateId": "secret-move",
		"type": "STRENGTH",
		"technicalDifficulty": "advanced",
		"bodyWeightScaling": 0,
		"description": "Tampered"
	}`, "other")
	fmt.Println(w.Code)
	// Output: 403
}

// Owner gets full permissions (READ, MODIFY, DELETE) on their own exercise.
func ExampleGetExercisePermissions_owner() {
	setupExampleDB()
	r := newRouter()

	// Create a private exercise owned by testuser.
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Bench Press",
		"templateId": "bench-press",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0.5,
		"description": "Barbell bench press",
		"version": 0
	}`)

	// Query permissions as the owner.
	w := doJSON(r, "GET", "/api/exercises/1/permissions", nil)

	var resp shared.PermissionsResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	fmt.Println(resp.Permissions)
	// Output: [READ MODIFY DELETE]
}

// Non-owner can read a public exercise but cannot modify or delete it.
func ExampleGetExercisePermissions_nonOwnerPublic() {
	setupExampleDB()
	r := newRouter()

	// Create a public exercise owned by testuser.
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Squat",
		"templateId": "squat",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0.5,
		"description": "Barbell squat",
		"version": 0,
		"public": true
	}`)

	// Query permissions as a different user.
	w := doRawAs(r, "GET", "/api/exercises/1/permissions", "", "other")

	var resp shared.PermissionsResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	fmt.Println(resp.Permissions)
	// Output: [READ]
}

// Non-owner has no permissions on a private exercise.
func ExampleGetExercisePermissions_nonOwnerPrivate() {
	setupExampleDB()
	r := newRouter()

	// Create a private exercise owned by testuser.
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Deadlift",
		"templateId": "deadlift",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0.5,
		"description": "Barbell deadlift",
		"version": 0
	}`)

	// Query permissions as a different user.
	w := doRawAs(r, "GET", "/api/exercises/1/permissions", "", "other")

	var resp shared.PermissionsResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	fmt.Println(resp.Permissions)
	// Output: []
}

// Creating an exercise without equipment. The owner is set from the
// authenticated user, not from the request body.
func ExampleCreateExercise_simple() {
	setupExampleDB()
	r := newRouter()

	w := doRaw(r, "POST", "/api/exercises", `{
		"name": "Push-up",
		"templateId": "push-up",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 1.0,
		"description": "Bodyweight push-up",
		"version": 0
	}`)

	var exercise models.Exercise
	json.Unmarshal(w.Body.Bytes(), &exercise)
	fmt.Println(w.Code)
	fmt.Println(exercise.Name)
	fmt.Println(exercise.Owner)
	fmt.Println(exercise.EquipmentIDs)
	// Output:
	// 201
	// Push-up
	// testuser
	// []
}

// Creating an exercise that requires equipment. The equipment must be created
// first via the equipment API (see [gesitr/internal/equipment/handlers]).
func ExampleCreateExercise_withEquipment() {
	setupExampleDB()
	r := newExampleRouter()

	// First, create the equipment via the equipment API.
	doRaw(r, "POST", "/api/equipment", `{
		"name": "barbell",
		"displayName": "Barbell",
		"description": "Standard barbell",
		"category": "free_weights",
		"templateId": "barbell"
	}`)

	// Create an exercise referencing the equipment by ID.
	w := doRaw(r, "POST", "/api/exercises", `{
		"name": "Bench Press",
		"templateId": "bench-press",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0.5,
		"description": "Barbell bench press",
		"version": 0,
		"equipmentIds": [1]
	}`)

	var exercise models.Exercise
	json.Unmarshal(w.Body.Bytes(), &exercise)
	fmt.Println(w.Code)
	fmt.Println(exercise.Name)
	fmt.Println(exercise.EquipmentIDs)
	// Output:
	// 201
	// Bench Press
	// [1]
}

// Owner can retrieve their own exercise with full details.
func ExampleGetExercise_owner() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/exercises", `{
		"name": "Bench Press",
		"templateId": "bench-press",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0.5,
		"description": "Barbell bench press",
		"version": 0
	}`)

	w := doJSON(r, "GET", "/api/exercises/1", nil)

	var exercise models.Exercise
	json.Unmarshal(w.Body.Bytes(), &exercise)
	fmt.Println(w.Code)
	fmt.Println(exercise.Name)
	fmt.Println(exercise.Owner)
	// Output:
	// 200
	// Bench Press
	// testuser
}

// Non-owner can retrieve a public exercise.
func ExampleGetExercise_nonOwnerPublic() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/exercises", `{
		"name": "Squat",
		"templateId": "squat",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0.5,
		"description": "Barbell squat",
		"version": 0,
		"public": true
	}`)

	w := doRawAs(r, "GET", "/api/exercises/1", "", "other")

	var exercise models.Exercise
	json.Unmarshal(w.Body.Bytes(), &exercise)
	fmt.Println(w.Code)
	fmt.Println(exercise.Name)
	// Output:
	// 200
	// Squat
}

// Non-owner cannot access a private exercise — returns 403 Forbidden.
func ExampleGetExercise_nonOwnerPrivate() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/exercises", `{
		"name": "Secret Move",
		"templateId": "secret-move",
		"type": "STRENGTH",
		"technicalDifficulty": "advanced",
		"bodyWeightScaling": 0,
		"description": "A private exercise",
		"version": 0
	}`)

	w := doRawAs(r, "GET", "/api/exercises/1", "", "other")
	fmt.Println(w.Code)
	// Output: 403
}

// Creating a rep-based exercise scheme for bicep curls. The exercise must
// exist first (see [CreateExercise]). The scheme defines how to perform
// the exercise: 3 sets of 12 reps at 15kg with 90s rest.
func ExampleCreateExerciseScheme_repBased() {
	setupExampleDB()
	r := newRouter()

	// Create the exercise first.
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Bicep Curl",
		"templateId": "bicep-curl",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0,
		"description": "Dumbbell bicep curl",
		"version": 0
	}`)

	// Create a rep-based scheme for this exercise.
	w := doRaw(r, "POST", "/api/exercise-schemes", `{
		"exerciseId": 1,
		"measurementType": "REP_BASED",
		"sets": 3,
		"reps": 12,
		"weight": 15.0,
		"restBetweenSets": 90
	}`)

	var scheme models.ExerciseScheme
	json.Unmarshal(w.Body.Bytes(), &scheme)
	fmt.Println(w.Code)
	fmt.Println(scheme.MeasurementType)
	fmt.Println(*scheme.Sets, "sets,", *scheme.Reps, "reps,", *scheme.Weight, "kg")
	// Output:
	// 201
	// REP_BASED
	// 3 sets, 12 reps, 15 kg
}

// Creating a time-based exercise scheme for an ergometer session. The scheme
// defines a 30-minute cardio session with no sets or reps.
func ExampleCreateExerciseScheme_timeBased() {
	setupExampleDB()
	r := newRouter()

	// Create the exercise first.
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Ergometer",
		"templateId": "ergometer",
		"type": "CARDIO",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0,
		"description": "Rowing ergometer",
		"version": 0
	}`)

	// Create a time-based scheme for this exercise.
	w := doRaw(r, "POST", "/api/exercise-schemes", `{
		"exerciseId": 1,
		"measurementType": "TIME_BASED",
		"duration": 1800
	}`)

	var scheme models.ExerciseScheme
	json.Unmarshal(w.Body.Bytes(), &scheme)
	fmt.Println(w.Code)
	fmt.Println(scheme.MeasurementType)
	fmt.Println(*scheme.Duration, "seconds")
	fmt.Println(scheme.Sets == nil, "- no sets for cardio")
	// Output:
	// 201
	// TIME_BASED
	// 1800 seconds
	// true - no sets for cardio
}

// Owner can retrieve their own exercise scheme.
func ExampleGetExerciseScheme_owner() {
	setupExampleDB()
	r := newRouter()

	// Create the exercise and a scheme.
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Squat",
		"templateId": "squat",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0.5,
		"description": "Barbell squat",
		"version": 0
	}`)
	doRaw(r, "POST", "/api/exercise-schemes", `{
		"exerciseId": 1,
		"measurementType": "REP_BASED",
		"sets": 5,
		"reps": 5,
		"weight": 100.0,
		"restBetweenSets": 180
	}`)

	w := doJSON(r, "GET", "/api/exercise-schemes/1", nil)

	var scheme models.ExerciseScheme
	json.Unmarshal(w.Body.Bytes(), &scheme)
	fmt.Println(w.Code)
	fmt.Println(scheme.MeasurementType)
	// Output:
	// 200
	// REP_BASED
}

// Non-owner can read a scheme if the linked exercise is public.
func ExampleGetExerciseScheme_nonOwnerPublicExercise() {
	setupExampleDB()
	r := newRouter()

	// Create a public exercise and a scheme for it.
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Push-up",
		"templateId": "push-up",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 1.0,
		"description": "Bodyweight push-up",
		"version": 0,
		"public": true
	}`)
	doRaw(r, "POST", "/api/exercise-schemes", `{
		"exerciseId": 1,
		"measurementType": "REP_BASED",
		"sets": 3,
		"reps": 20
	}`)

	// Another user can read the scheme because the exercise is public.
	w := doRawAs(r, "GET", "/api/exercise-schemes/1", "", "other")

	var scheme models.ExerciseScheme
	json.Unmarshal(w.Body.Bytes(), &scheme)
	fmt.Println(w.Code)
	fmt.Println(scheme.MeasurementType)
	// Output:
	// 200
	// REP_BASED
}

// Non-owner cannot read a scheme if the linked exercise is private.
func ExampleGetExerciseScheme_nonOwnerPrivateExercise() {
	setupExampleDB()
	r := newRouter()

	// Create a private exercise and a scheme for it.
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Secret Move",
		"templateId": "secret-move",
		"type": "STRENGTH",
		"technicalDifficulty": "advanced",
		"bodyWeightScaling": 0,
		"description": "A private exercise",
		"version": 0
	}`)
	doRaw(r, "POST", "/api/exercise-schemes", `{
		"exerciseId": 1,
		"measurementType": "REP_BASED",
		"sets": 3,
		"reps": 8
	}`)

	// Another user is denied because the exercise is private.
	w := doRawAs(r, "GET", "/api/exercise-schemes/1", "", "other")
	fmt.Println(w.Code)
	// Output: 403
}

// Owner can update their exercise scheme.
func ExampleUpdateExerciseScheme_owner() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/exercises", `{
		"name": "Bicep Curl",
		"templateId": "bicep-curl",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0,
		"description": "Dumbbell bicep curl",
		"version": 0
	}`)
	doRaw(r, "POST", "/api/exercise-schemes", `{
		"exerciseId": 1,
		"measurementType": "REP_BASED",
		"sets": 3,
		"reps": 10,
		"weight": 12.0,
		"restBetweenSets": 60
	}`)

	// Update the scheme — increase weight and reduce reps.
	w := doRaw(r, "PUT", "/api/exercise-schemes/1", `{
		"exerciseId": 1,
		"measurementType": "REP_BASED",
		"sets": 3,
		"reps": 8,
		"weight": 15.0,
		"restBetweenSets": 90
	}`)

	var scheme models.ExerciseScheme
	json.Unmarshal(w.Body.Bytes(), &scheme)
	fmt.Println(w.Code)
	fmt.Println(*scheme.Reps, "reps,", *scheme.Weight, "kg")
	// Output:
	// 200
	// 8 reps, 15 kg
}

// Non-owner cannot update a scheme even if the linked exercise is public.
func ExampleUpdateExerciseScheme_nonOwnerPublicExercise() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/exercises", `{
		"name": "Push-up",
		"templateId": "push-up",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 1.0,
		"description": "Bodyweight push-up",
		"version": 0,
		"public": true
	}`)
	doRaw(r, "POST", "/api/exercise-schemes", `{
		"exerciseId": 1,
		"measurementType": "REP_BASED",
		"sets": 3,
		"reps": 20
	}`)

	w := doRawAs(r, "PUT", "/api/exercise-schemes/1", `{
		"exerciseId": 1,
		"measurementType": "REP_BASED",
		"sets": 5,
		"reps": 50
	}`, "other")
	fmt.Println(w.Code)
	// Output: 403
}

// Non-owner cannot update a scheme for a private exercise.
func ExampleUpdateExerciseScheme_nonOwnerPrivateExercise() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/exercises", `{
		"name": "Secret Move",
		"templateId": "secret-move",
		"type": "STRENGTH",
		"technicalDifficulty": "advanced",
		"bodyWeightScaling": 0,
		"description": "A private exercise",
		"version": 0
	}`)
	doRaw(r, "POST", "/api/exercise-schemes", `{
		"exerciseId": 1,
		"measurementType": "REP_BASED",
		"sets": 3,
		"reps": 8
	}`)

	w := doRawAs(r, "PUT", "/api/exercise-schemes/1", `{
		"exerciseId": 1,
		"measurementType": "REP_BASED",
		"sets": 5,
		"reps": 50
	}`, "other")
	fmt.Println(w.Code)
	// Output: 403
}

// Owner sees their own schemes in the list.
func ExampleListExerciseSchemes_owner() {
	setupExampleDB()
	r := newRouter()

	// Create a private exercise and a scheme.
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Squat",
		"templateId": "squat",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0.5,
		"description": "Barbell squat",
		"version": 0
	}`)
	doRaw(r, "POST", "/api/exercise-schemes", `{
		"exerciseId": 1,
		"measurementType": "REP_BASED",
		"sets": 5,
		"reps": 5
	}`)

	w := doJSON(r, "GET", "/api/exercise-schemes", nil)

	var schemes []models.ExerciseScheme
	json.Unmarshal(w.Body.Bytes(), &schemes)
	fmt.Println(w.Code)
	fmt.Println(len(schemes))
	fmt.Println(schemes[0].MeasurementType)
	// Output:
	// 200
	// 1
	// REP_BASED
}

// Non-owner sees schemes for public exercises in the list.
func ExampleListExerciseSchemes_nonOwnerPublicExercise() {
	setupExampleDB()
	r := newRouter()

	// Create a public exercise and a scheme.
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Push-up",
		"templateId": "push-up",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 1.0,
		"description": "Bodyweight push-up",
		"version": 0,
		"public": true
	}`)
	doRaw(r, "POST", "/api/exercise-schemes", `{
		"exerciseId": 1,
		"measurementType": "REP_BASED",
		"sets": 3,
		"reps": 20
	}`)

	// Another user sees schemes for public exercises.
	w := doRawAs(r, "GET", "/api/exercise-schemes", "", "other")

	var schemes []models.ExerciseScheme
	json.Unmarshal(w.Body.Bytes(), &schemes)
	fmt.Println(w.Code)
	fmt.Println(len(schemes))
	// Output:
	// 200
	// 1
}

// Non-owner does not see schemes for private exercises.
func ExampleListExerciseSchemes_nonOwnerPrivateExercise() {
	setupExampleDB()
	r := newRouter()

	// Create a private exercise and a scheme.
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Secret Move",
		"templateId": "secret-move",
		"type": "STRENGTH",
		"technicalDifficulty": "advanced",
		"bodyWeightScaling": 0,
		"description": "A private exercise",
		"version": 0
	}`)
	doRaw(r, "POST", "/api/exercise-schemes", `{
		"exerciseId": 1,
		"measurementType": "REP_BASED",
		"sets": 3,
		"reps": 8
	}`)

	// Another user sees an empty list.
	w := doRawAs(r, "GET", "/api/exercise-schemes", "", "other")

	var schemes []models.ExerciseScheme
	json.Unmarshal(w.Body.Bytes(), &schemes)
	fmt.Println(w.Code)
	fmt.Println(len(schemes))
	// Output:
	// 200
	// 0
}

// Owner can retrieve version history. Creating an exercise produces version 0,
// updating it produces version 1. Each version contains a snapshot of the
// exercise at that point in time.
func ExampleGetExerciseVersion_owner() {
	setupExampleDB()
	r := newRouter()

	// Create an exercise (version 0).
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Squat",
		"templateId": "squat",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0.5,
		"description": "Barbell squat",
		"version": 0
	}`)

	// Verify version 0 exists.
	w0 := doJSON(r, "GET", "/api/exercises/templates/squat/versions/0", nil)
	var v0 shared.VersionEntry
	json.Unmarshal(w0.Body.Bytes(), &v0)
	fmt.Println("v0:", w0.Code, "version", v0.Version)

	// Parse the snapshot to check the exercise name at version 0.
	var snap0 models.Exercise
	json.Unmarshal(v0.Snapshot, &snap0)
	fmt.Println("v0 name:", snap0.Name)

	// Update the exercise (creates version 1).
	doRaw(r, "PUT", "/api/exercises/1", `{
		"name": "Back Squat",
		"templateId": "squat",
		"type": "STRENGTH",
		"technicalDifficulty": "intermediate",
		"bodyWeightScaling": 0.5,
		"description": "Barbell back squat",
		"force": ["PUSH"],
		"primaryMuscles": ["QUADS"]
	}`)

	// Verify version 1 exists with the updated name.
	w1 := doJSON(r, "GET", "/api/exercises/templates/squat/versions/1", nil)
	var v1 shared.VersionEntry
	json.Unmarshal(w1.Body.Bytes(), &v1)
	fmt.Println("v1:", w1.Code, "version", v1.Version)

	var snap1 models.Exercise
	json.Unmarshal(v1.Snapshot, &snap1)
	fmt.Println("v1 name:", snap1.Name)

	// Output:
	// v0: 200 version 0
	// v0 name: Squat
	// v1: 200 version 1
	// v1 name: Back Squat
}

// Non-owner can retrieve version history for a public exercise.
func ExampleGetExerciseVersion_nonOwnerPublic() {
	setupExampleDB()
	r := newRouter()

	// Create a public exercise (version 0).
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Push-up",
		"templateId": "push-up",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 1.0,
		"description": "Bodyweight push-up",
		"version": 0,
		"public": true
	}`)

	// Update it (version 1).
	doRaw(r, "PUT", "/api/exercises/1", `{
		"name": "Push-up",
		"templateId": "push-up",
		"type": "STRENGTH",
		"technicalDifficulty": "intermediate",
		"bodyWeightScaling": 1.0,
		"description": "Strict push-up",
		"public": true,
		"force": ["PUSH"],
		"primaryMuscles": ["CHEST"]
	}`)

	// Another user can access both versions.
	w0 := doRawAs(r, "GET", "/api/exercises/templates/push-up/versions/0", "", "other")
	var v0 shared.VersionEntry
	json.Unmarshal(w0.Body.Bytes(), &v0)
	fmt.Println("v0:", w0.Code, "version", v0.Version)

	w1 := doRawAs(r, "GET", "/api/exercises/templates/push-up/versions/1", "", "other")
	var v1 shared.VersionEntry
	json.Unmarshal(w1.Body.Bytes(), &v1)
	fmt.Println("v1:", w1.Code, "version", v1.Version)

	var snap1 models.Exercise
	json.Unmarshal(v1.Snapshot, &snap1)
	fmt.Println("v1 description:", snap1.Description)

	// Output:
	// v0: 200 version 0
	// v1: 200 version 1
	// v1 description: Strict push-up
}

// Non-owner cannot retrieve version history for a private exercise.
func ExampleGetExerciseVersion_nonOwnerPrivate() {
	setupExampleDB()
	r := newRouter()

	// Create a private exercise.
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Secret Move",
		"templateId": "secret-move",
		"type": "STRENGTH",
		"technicalDifficulty": "advanced",
		"bodyWeightScaling": 0,
		"description": "A private exercise",
		"version": 0
	}`)

	w := doRawAs(r, "GET", "/api/exercises/templates/secret-move/versions/0", "", "other")
	fmt.Println(w.Code)
	// Output: 403
}

// Version history survives exercise deletion. The exercise is hard-deleted
// but the history snapshots remain accessible via the templateId.
func ExampleGetExerciseVersion_afterDelete() {
	setupExampleDB()
	r := newRouter()

	// Create an exercise (version 0).
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Overhead Press",
		"templateId": "ohp",
		"type": "STRENGTH",
		"technicalDifficulty": "intermediate",
		"bodyWeightScaling": 0,
		"description": "Standing overhead press",
		"version": 0
	}`)

	// Verify it exists.
	w := doJSON(r, "GET", "/api/exercises/1", nil)
	fmt.Println("before delete:", w.Code)

	// Delete the exercise.
	wd := doJSON(r, "DELETE", "/api/exercises/1", nil)
	fmt.Println("delete:", wd.Code)

	// Exercise is gone.
	wg := doJSON(r, "GET", "/api/exercises/1", nil)
	fmt.Println("after delete:", wg.Code)

	// Version history is still accessible.
	wv := doJSON(r, "GET", "/api/exercises/templates/ohp/versions/0", nil)
	var v0 shared.VersionEntry
	json.Unmarshal(wv.Body.Bytes(), &v0)
	fmt.Println("version:", wv.Code, "version", v0.Version)

	var snap models.Exercise
	json.Unmarshal(v0.Snapshot, &snap)
	fmt.Println("snapshot name:", snap.Name)

	// Output:
	// before delete: 200
	// delete: 204
	// after delete: 404
	// version: 200 version 0
	// snapshot name: Overhead Press
}
