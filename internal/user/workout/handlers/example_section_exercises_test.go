package handlers

import (
	"encoding/json"
	"fmt"

	"gesitr/internal/user/workout/models"

	"github.com/gin-gonic/gin"
)

// createExerciseSchemeForExample creates an exercise and a scheme, returning
// the full router. Used as setup for section exercise examples.
func createExerciseSchemeForExample(r *gin.Engine) {
	doRaw(r, "POST", "/api/exercises", `{
		"name": "Bench Press",
		"templateId": "bench-press",
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0.5,
		"description": "Barbell bench press"
	}`)
	doRaw(r, "POST", "/api/exercise-schemes", `{
		"exerciseId": 1,
		"measurementType": "REP_BASED",
		"sets": 3,
		"reps": 10,
		"weight": 60.0,
		"restBetweenSets": 90
	}`)
}

// Creating a section exercise requires a workout, a section, and an exercise
// scheme. The full hierarchy: Workout → Section → SectionExercise → ExerciseScheme.
func ExampleCreateWorkoutSectionExercise() {
	setupExampleDB()
	r := newRouter()

	// Create the exercise and scheme first.
	createExerciseSchemeForExample(r)

	// Create a workout and section.
	doRaw(r, "POST", "/api/user/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/user/workout-sections", `{
		"workoutId": 1, "type": "main", "position": 0
	}`)

	// Add the exercise scheme to the section.
	w := doRaw(r, "POST", "/api/user/workout-section-exercises", `{
		"workoutSectionId": 1,
		"exerciseSchemeId": 1,
		"position": 0
	}`)

	var se models.WorkoutSectionExercise
	json.Unmarshal(w.Body.Bytes(), &se)
	fmt.Println(w.Code)
	fmt.Println(se.WorkoutSectionID)
	fmt.Println(se.ExerciseSchemeID)
	// Output:
	// 201
	// 1
	// 1
}

// Non-owner cannot add exercises to another user's section.
func ExampleCreateWorkoutSectionExercise_nonOwnerDenied() {
	setupExampleDB()
	r := newRouter()

	createExerciseSchemeForExample(r)
	doRaw(r, "POST", "/api/user/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/user/workout-sections", `{
		"workoutId": 1, "type": "main", "position": 0
	}`)

	w := doRawAs(r, "POST", "/api/user/workout-section-exercises", `{
		"workoutSectionId": 1,
		"exerciseSchemeId": 1,
		"position": 0
	}`, "bob")
	fmt.Println(w.Code)
	// Output: 403
}

// ListWorkoutSectionExercises returns exercises in sections owned by the
// current user. Filter by workoutSectionId.
func ExampleListWorkoutSectionExercises_owner() {
	setupExampleDB()
	r := newRouter()

	createExerciseSchemeForExample(r)
	doRaw(r, "POST", "/api/user/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/user/workout-sections", `{
		"workoutId": 1, "type": "main", "position": 0
	}`)
	doRaw(r, "POST", "/api/user/workout-section-exercises", `{
		"workoutSectionId": 1, "exerciseSchemeId": 1, "position": 0
	}`)

	w := doRaw(r, "GET", "/api/user/workout-section-exercises?workoutSectionId=1", "")

	var exercises []models.WorkoutSectionExercise
	json.Unmarshal(w.Body.Bytes(), &exercises)
	fmt.Println(w.Code)
	fmt.Println(len(exercises))
	fmt.Println(exercises[0].ExerciseSchemeID)
	// Output:
	// 200
	// 1
	// 1
}

// Non-owner sees an empty list for another user's section exercises.
func ExampleListWorkoutSectionExercises_nonOwner() {
	setupExampleDB()
	r := newRouter()

	createExerciseSchemeForExample(r)
	doRaw(r, "POST", "/api/user/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/user/workout-sections", `{
		"workoutId": 1, "type": "main", "position": 0
	}`)
	doRaw(r, "POST", "/api/user/workout-section-exercises", `{
		"workoutSectionId": 1, "exerciseSchemeId": 1, "position": 0
	}`)

	w := doRawAs(r, "GET", "/api/user/workout-section-exercises?workoutSectionId=1", "", "bob")

	var exercises []models.WorkoutSectionExercise
	json.Unmarshal(w.Body.Bytes(), &exercises)
	fmt.Println(w.Code)
	fmt.Println(len(exercises))
	// Output:
	// 200
	// 0
}

// Owner can delete an exercise from their section.
func ExampleDeleteWorkoutSectionExercise_owner() {
	setupExampleDB()
	r := newRouter()

	createExerciseSchemeForExample(r)
	doRaw(r, "POST", "/api/user/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/user/workout-sections", `{
		"workoutId": 1, "type": "main", "position": 0
	}`)
	doRaw(r, "POST", "/api/user/workout-section-exercises", `{
		"workoutSectionId": 1, "exerciseSchemeId": 1, "position": 0
	}`)

	w := doJSON(r, "DELETE", "/api/user/workout-section-exercises/1", nil)
	fmt.Println(w.Code)

	// Section exercise is gone.
	wl := doRaw(r, "GET", "/api/user/workout-section-exercises?workoutSectionId=1", "")
	var exercises []models.WorkoutSectionExercise
	json.Unmarshal(wl.Body.Bytes(), &exercises)
	fmt.Println(len(exercises))
	// Output:
	// 204
	// 0
}

// Non-owner cannot delete exercises from another user's section.
func ExampleDeleteWorkoutSectionExercise_nonOwnerDenied() {
	setupExampleDB()
	r := newRouter()

	createExerciseSchemeForExample(r)
	doRaw(r, "POST", "/api/user/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/user/workout-sections", `{
		"workoutId": 1, "type": "main", "position": 0
	}`)
	doRaw(r, "POST", "/api/user/workout-section-exercises", `{
		"workoutSectionId": 1, "exerciseSchemeId": 1, "position": 0
	}`)

	w := doRawAs(r, "DELETE", "/api/user/workout-section-exercises/1", "", "bob")
	fmt.Println(w.Code)
	// Output: 403
}

// Full hierarchy: create a workout with sections and exercises, then retrieve
// the complete tree via GetWorkout.
func ExampleGetWorkout_fullHierarchy() {
	setupExampleDB()
	r := newRouter()

	// Create exercise and scheme.
	createExerciseSchemeForExample(r)

	// Build the workout hierarchy.
	doRaw(r, "POST", "/api/user/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/user/workout-sections", `{
		"workoutId": 1, "type": "main", "label": "Compound", "position": 0
	}`)
	doRaw(r, "POST", "/api/user/workout-section-exercises", `{
		"workoutSectionId": 1, "exerciseSchemeId": 1, "position": 0
	}`)

	// GetWorkout returns the full tree.
	w := doJSON(r, "GET", "/api/user/workouts/1", nil)

	var workout models.Workout
	json.Unmarshal(w.Body.Bytes(), &workout)
	fmt.Println(w.Code)
	fmt.Println(workout.Name)
	fmt.Println(len(workout.Sections), "section(s)")
	fmt.Println(*workout.Sections[0].Label)
	fmt.Println(len(workout.Sections[0].Exercises), "exercise(s)")
	fmt.Println(workout.Sections[0].Exercises[0].ExerciseSchemeID)
	// Output:
	// 200
	// Push Day
	// 1 section(s)
	// Compound
	// 1 exercise(s)
	// 1
}
