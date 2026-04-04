package handlers

import (
	"encoding/json"
	"fmt"

	"gesitr/internal/compendium/workout/models"

	"github.com/gin-gonic/gin"
)

// createExerciseForExample creates an exercise, returning the full router.
// Used as setup for section item examples.
func createExerciseForExample(r *gin.Engine) {
	doRaw(r, "POST", "/api/exercises", `{
		"names": ["Bench Press"],
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0.5,
		"description": "Barbell bench press"
	}`)
}

// Creating a section item requires a workout, a section, and an exercise.
// The full hierarchy: Workout → Section → SectionItem → Exercise.
func ExampleCreateWorkoutSectionItem() {
	setupExampleDB()
	r := newRouter()

	// Create the exercise first.
	createExerciseForExample(r)

	// Create a workout and section.
	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/workout-sections", `{
		"workoutId": 1, "type": "main", "position": 0
	}`)

	// Add the exercise to the section.
	w := doRaw(r, "POST", "/api/workout-section-items", `{
		"workoutSectionId": 1,
		"type": "exercise",
		"exerciseId": 1,
		"position": 0
	}`)

	var item models.WorkoutSectionItem
	json.Unmarshal(w.Body.Bytes(), &item)
	fmt.Println(w.Code)
	fmt.Println(item.WorkoutSectionID)
	fmt.Println(item.Type)
	fmt.Println(*item.ExerciseID)
	// Output:
	// 201
	// 1
	// exercise
	// 1
}

// Non-owner cannot add items to another user's section.
func ExampleCreateWorkoutSectionItem_nonOwnerDenied() {
	setupExampleDB()
	r := newRouter()

	createExerciseForExample(r)
	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/workout-sections", `{
		"workoutId": 1, "type": "main", "position": 0
	}`)

	w := doRawAs(r, "POST", "/api/workout-section-items", `{
		"workoutSectionId": 1,
		"type": "exercise",
		"exerciseId": 1,
		"position": 0
	}`, "bob")
	fmt.Println(w.Code)
	// Output: 403
}

// ListWorkoutSectionItems returns items in sections owned by the
// current user. Filter by workoutSectionId.
func ExampleListWorkoutSectionItems_owner() {
	setupExampleDB()
	r := newRouter()

	createExerciseForExample(r)
	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/workout-sections", `{
		"workoutId": 1, "type": "main", "position": 0
	}`)
	doRaw(r, "POST", "/api/workout-section-items", `{
		"workoutSectionId": 1, "type": "exercise", "exerciseId": 1, "position": 0
	}`)

	w := doRaw(r, "GET", "/api/workout-section-items?workoutSectionId=1", "")

	var items []models.WorkoutSectionItem
	json.Unmarshal(w.Body.Bytes(), &items)
	fmt.Println(w.Code)
	fmt.Println(len(items))
	fmt.Println(*items[0].ExerciseID)
	// Output:
	// 200
	// 1
	// 1
}

// Non-owner sees an empty list for another user's section items.
func ExampleListWorkoutSectionItems_nonOwner() {
	setupExampleDB()
	r := newRouter()

	createExerciseForExample(r)
	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/workout-sections", `{
		"workoutId": 1, "type": "main", "position": 0
	}`)
	doRaw(r, "POST", "/api/workout-section-items", `{
		"workoutSectionId": 1, "type": "exercise", "exerciseId": 1, "position": 0
	}`)

	w := doRawAs(r, "GET", "/api/workout-section-items?workoutSectionId=1", "", "bob")

	var items []models.WorkoutSectionItem
	json.Unmarshal(w.Body.Bytes(), &items)
	fmt.Println(w.Code)
	fmt.Println(len(items))
	// Output:
	// 200
	// 0
}

// Owner can delete an item from their section.
func ExampleDeleteWorkoutSectionItem_owner() {
	setupExampleDB()
	r := newRouter()

	createExerciseForExample(r)
	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/workout-sections", `{
		"workoutId": 1, "type": "main", "position": 0
	}`)
	doRaw(r, "POST", "/api/workout-section-items", `{
		"workoutSectionId": 1, "type": "exercise", "exerciseId": 1, "position": 0
	}`)

	w := doJSON(r, "DELETE", "/api/workout-section-items/1", nil)
	fmt.Println(w.Code)

	// Section item is gone.
	wl := doRaw(r, "GET", "/api/workout-section-items?workoutSectionId=1", "")
	var items []models.WorkoutSectionItem
	json.Unmarshal(wl.Body.Bytes(), &items)
	fmt.Println(len(items))
	// Output:
	// 204
	// 0
}

// Non-owner cannot delete items from another user's section.
func ExampleDeleteWorkoutSectionItem_nonOwnerDenied() {
	setupExampleDB()
	r := newRouter()

	createExerciseForExample(r)
	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/workout-sections", `{
		"workoutId": 1, "type": "main", "position": 0
	}`)
	doRaw(r, "POST", "/api/workout-section-items", `{
		"workoutSectionId": 1, "type": "exercise", "exerciseId": 1, "position": 0
	}`)

	w := doRawAs(r, "DELETE", "/api/workout-section-items/1", "", "bob")
	fmt.Println(w.Code)
	// Output: 403
}

// Full hierarchy: create a workout with sections and items, then retrieve
// the complete tree via GetWorkout.
func ExampleGetWorkout_fullHierarchy() {
	setupExampleDB()
	r := newRouter()

	// Create exercise.
	createExerciseForExample(r)

	// Build the workout hierarchy.
	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/workout-sections", `{
		"workoutId": 1, "type": "main", "label": "Compound", "position": 0
	}`)
	doRaw(r, "POST", "/api/workout-section-items", `{
		"workoutSectionId": 1, "type": "exercise", "exerciseId": 1, "position": 0
	}`)

	// GetWorkout returns the full tree.
	w := doJSON(r, "GET", "/api/workouts/1", nil)

	var workout models.Workout
	json.Unmarshal(w.Body.Bytes(), &workout)
	fmt.Println(w.Code)
	fmt.Println(workout.Name)
	fmt.Println(len(workout.Sections), "section(s)")
	fmt.Println(*workout.Sections[0].Label)
	fmt.Println(len(workout.Sections[0].Items), "item(s)")
	fmt.Println(*workout.Sections[0].Items[0].ExerciseID)
	// Output:
	// 200
	// Push Day
	// 1 section(s)
	// Compound
	// 1 item(s)
	// 1
}
