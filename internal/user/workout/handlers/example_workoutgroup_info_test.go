package handlers

import (
	"encoding/json"
	"fmt"

	"gesitr/internal/user/workout/models"
)

// When bob lists workouts, shared workouts include workoutGroup info
// with the group name and his membership role.
func ExampleListWorkouts_groupMemberSeesWorkoutGroupInfo() {
	setupExampleDB()
	r := newRouter()

	// Bob interacts with the API to establish his profile
	doRawAs(r, "GET", "/api/workouts", "", "bob")

	// Alice creates a workout
	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)

	// Alice creates a workout group and invites bob
	doRaw(r, "POST", "/api/user/workout-groups", `{
		"name": "Gym Buddies", "workoutId": 1
	}`)
	doRaw(r, "POST", "/api/user/workout-group-memberships", `{
		"groupId": 1, "userId": "bob", "role": "member"
	}`)

	// Bob lists workouts — should see Alice's shared workout with group info
	w := doRawAs(r, "GET", "/api/workouts", "", "bob")

	var page struct {
		Items []models.Workout `json:"items"`
	}
	json.Unmarshal(w.Body.Bytes(), &page)
	fmt.Println(w.Code)
	fmt.Println(len(page.Items))
	fmt.Println(page.Items[0].Name)
	fmt.Println(page.Items[0].WorkoutGroup.GroupName)
	fmt.Println(page.Items[0].WorkoutGroup.Membership)
	// Output:
	// 200
	// 1
	// Push Day
	// Gym Buddies
	// invited
}

// The owner's own workouts do not include workoutGroup info,
// even if a group exists for that workout.
func ExampleListWorkouts_ownerDoesNotSeeWorkoutGroupInfo() {
	setupExampleDB()
	r := newRouter()

	// Alice creates a workout and a group for it
	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/user/workout-groups", `{
		"name": "Gym Buddies", "workoutId": 1
	}`)

	// Alice lists her workouts — workoutGroup should be absent
	w := doJSON(r, "GET", "/api/workouts", nil)

	var page struct {
		Items []models.Workout `json:"items"`
	}
	json.Unmarshal(w.Body.Bytes(), &page)
	fmt.Println(w.Code)
	fmt.Println(len(page.Items))
	fmt.Println(page.Items[0].Name)
	fmt.Println(page.Items[0].WorkoutGroup == nil)
	// Output:
	// 200
	// 1
	// Push Day
	// true
}

// GetWorkout for a group member includes workoutGroup info.
func ExampleGetWorkout_groupMemberSeesWorkoutGroupInfo() {
	setupExampleDB()
	r := newRouter()

	// Bob interacts with the API to establish his profile
	doRawAs(r, "GET", "/api/workouts", "", "bob")

	// Alice creates a workout, group, and invites bob
	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/user/workout-groups", `{
		"name": "Gym Buddies", "workoutId": 1
	}`)
	doRaw(r, "POST", "/api/user/workout-group-memberships", `{
		"groupId": 1, "userId": "bob", "role": "member"
	}`)

	// Bob fetches the workout directly
	w := doRawAs(r, "GET", "/api/workouts/1", "", "bob")

	var workout models.Workout
	json.Unmarshal(w.Body.Bytes(), &workout)
	fmt.Println(w.Code)
	fmt.Println(workout.Name)
	fmt.Println(workout.WorkoutGroup.GroupName)
	fmt.Println(workout.WorkoutGroup.Membership)
	// Output:
	// 200
	// Push Day
	// Gym Buddies
	// invited
}

// UpdateWorkout response for an admin group member includes workoutGroup info.
func ExampleUpdateWorkout_groupAdminSeesWorkoutGroupInfo() {
	setupExampleDB()
	r := newRouter()

	// Bob interacts with the API to establish his profile
	doRawAs(r, "GET", "/api/workouts", "", "bob")

	// Alice creates a workout, group, and adds bob as admin
	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/user/workout-groups", `{
		"name": "Gym Buddies", "workoutId": 1
	}`)
	doRaw(r, "POST", "/api/user/workout-group-memberships", `{
		"groupId": 1, "userId": "bob", "role": "member"
	}`)
	// Promote bob to admin so he can modify
	doRaw(r, "PUT", "/api/user/workout-group-memberships/1", `{"role": "admin"}`)

	// Bob updates the workout — response should include workoutGroup
	w := doRawAs(r, "PUT", "/api/workouts/1", `{"name": "Upper Push"}`, "bob")

	var workout models.Workout
	json.Unmarshal(w.Body.Bytes(), &workout)
	fmt.Println(w.Code)
	fmt.Println(workout.Name)
	fmt.Println(workout.WorkoutGroup.GroupName)
	fmt.Println(workout.WorkoutGroup.Membership)
	// Output:
	// 200
	// Upper Push
	// Gym Buddies
	// admin
}
