package handlers

import (
	"encoding/json"
	"fmt"

	"gesitr/internal/shared"
)

// Owner gets full permissions: READ, MODIFY, DELETE.
func ExampleGetWorkoutPermissions_owner() {
	setupExampleDB()
	r := newRouter()

	doRaw(r, "POST", "/api/user/workouts", `{"name": "Push Day"}`)

	w := doRaw(r, "GET", "/api/user/workouts/1/permissions", "")

	var resp shared.PermissionsResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	fmt.Println(w.Code)
	fmt.Println(resp.Permissions)
	// Output:
	// 200
	// [READ MODIFY DELETE]
}

// A group member with "member" role gets READ only.
func ExampleGetWorkoutPermissions_groupMember() {
	setupExampleDB()
	r := newRouter()

	doRawAs(r, "GET", "/api/user/workouts", "", "bob")

	doRaw(r, "POST", "/api/user/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/user/workout-groups", `{
		"name": "Gym Buddies", "workoutId": 1
	}`)
	doRaw(r, "POST", "/api/user/workout-group-memberships", `{
		"groupId": 1, "userId": "bob", "role": "member"
	}`)

	w := doRawAs(r, "GET", "/api/user/workouts/1/permissions", "", "bob")

	var resp shared.PermissionsResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	fmt.Println(w.Code)
	fmt.Println(resp.Permissions)
	// Output:
	// 200
	// [READ]
}

// A group admin gets READ and MODIFY.
func ExampleGetWorkoutPermissions_groupAdmin() {
	setupExampleDB()
	r := newRouter()

	doRawAs(r, "GET", "/api/user/workouts", "", "bob")

	doRaw(r, "POST", "/api/user/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/user/workout-groups", `{
		"name": "Gym Buddies", "workoutId": 1
	}`)
	doRaw(r, "POST", "/api/user/workout-group-memberships", `{
		"groupId": 1, "userId": "bob", "role": "member"
	}`)
	// Promote bob to admin
	doRaw(r, "PUT", "/api/user/workout-group-memberships/1", `{"role": "admin"}`)

	w := doRawAs(r, "GET", "/api/user/workouts/1/permissions", "", "bob")

	var resp shared.PermissionsResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	fmt.Println(w.Code)
	fmt.Println(resp.Permissions)
	// Output:
	// 200
	// [READ MODIFY]
}

// A user with no access gets 404.
func ExampleGetWorkoutPermissions_noAccess() {
	setupExampleDB()
	r := newRouter()

	doRawAs(r, "GET", "/api/user/workouts", "", "bob")

	doRaw(r, "POST", "/api/user/workouts", `{"name": "Push Day"}`)

	w := doRawAs(r, "GET", "/api/user/workouts/1/permissions", "", "bob")

	fmt.Println(w.Code)
	// Output:
	// 404
}

// An invited member gets READ only.
func ExampleGetWorkoutPermissions_invitedMember() {
	setupExampleDB()
	r := newRouter()

	doRawAs(r, "GET", "/api/user/workouts", "", "bob")

	doRaw(r, "POST", "/api/user/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/user/workout-groups", `{
		"name": "Gym Buddies", "workoutId": 1
	}`)
	doRaw(r, "POST", "/api/user/workout-group-memberships", `{
		"groupId": 1, "userId": "bob", "role": "invited"
	}`)

	w := doRawAs(r, "GET", "/api/user/workouts/1/permissions", "", "bob")

	var resp shared.PermissionsResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	fmt.Println(w.Code)
	fmt.Println(resp.Permissions)
	// Output:
	// 200
	// [READ]
}
