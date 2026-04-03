package handlers

import (
	"fmt"
)

// Accepting an invitation without creating exercise schemes returns 400.
// Flow: alice creates a workout with 1 exercise, creates a group, invites bob.
// Bob tries to accept but hasn't set up his own schemes yet → 400.
func ExampleAcceptWorkoutGroupInvitation_denied() {
	setupExampleDB()
	r := newRouter()

	// Bob interacts with the API so the test router recognises him
	doRawAs(r, "GET", "/api/workouts", "", "bob")

	// Alice creates an exercise and scheme
	doRaw(r, "POST", "/api/exercises", `{
		"names": ["Bench Press"], "type": "STRENGTH",
		"technicalDifficulty": "beginner", "bodyWeightScaling": 0.5,
		"description": "Barbell bench press", "public": true
	}`)
	doRaw(r, "POST", "/api/exercise-schemes", `{
		"exerciseId": 1, "measurementType": "REP_BASED",
		"sets": 3, "reps": 10, "weight": 60.0
	}`)

	// Alice creates a workout with one section and one exercise item
	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/workout-sections", `{
		"workoutId": 1, "type": "main", "position": 0
	}`)
	doRaw(r, "POST", "/api/workout-section-items", `{
		"workoutSectionId": 1, "type": "exercise",
		"exerciseSchemeId": 1, "position": 0
	}`)

	// Alice creates a workout group and invites bob
	doRaw(r, "POST", "/api/user/workout-groups", `{
		"name": "Gym Buddies", "workoutId": 1
	}`)
	doRaw(r, "POST", "/api/user/workout-group-memberships", `{
		"groupId": 1, "userId": "bob", "role": "member"
	}`)

	// Bob tries to accept without having created his exercise schemes
	w := doRawAs(r, "POST", "/api/workouts/1/group/accept", "", "bob")
	fmt.Println(w.Code)
	// Output: 400
}

// Accepting an invitation after creating exercise schemes succeeds.
// Flow: same setup as above, but bob creates his own scheme for the
// workout item before accepting → 200.
func ExampleAcceptWorkoutGroupInvitation_success() {
	setupExampleDB()
	r := newRouter()

	// Bob interacts with the API to establish his profile
	doRawAs(r, "GET", "/api/workouts", "", "bob")

	// Alice creates an exercise and scheme
	doRaw(r, "POST", "/api/exercises", `{
		"names": ["Bench Press"], "type": "STRENGTH",
		"technicalDifficulty": "beginner", "bodyWeightScaling": 0.5,
		"description": "Barbell bench press", "public": true
	}`)
	doRaw(r, "POST", "/api/exercise-schemes", `{
		"exerciseId": 1, "measurementType": "REP_BASED",
		"sets": 3, "reps": 10, "weight": 60.0
	}`)

	// Alice creates a workout with one section and one exercise item
	doRaw(r, "POST", "/api/workouts", `{"name": "Push Day"}`)
	doRaw(r, "POST", "/api/workout-sections", `{
		"workoutId": 1, "type": "main", "position": 0
	}`)
	doRaw(r, "POST", "/api/workout-section-items", `{
		"workoutSectionId": 1, "type": "exercise",
		"exerciseSchemeId": 1, "position": 0
	}`)

	// Alice creates a workout group and invites bob
	doRaw(r, "POST", "/api/user/workout-groups", `{
		"name": "Gym Buddies", "workoutId": 1
	}`)
	doRaw(r, "POST", "/api/user/workout-group-memberships", `{
		"groupId": 1, "userId": "bob", "role": "member"
	}`)

	// Bob creates his own exercise scheme linked to the workout section item (id=1)
	doRawAs(r, "POST", "/api/exercise-schemes", `{
		"exerciseId": 1, "measurementType": "REP_BASED",
		"sets": 5, "reps": 5, "weight": 80.0,
		"workoutSectionItemId": 1
	}`, "bob")

	// Bob accepts the invitation — should succeed now
	w := doRawAs(r, "POST", "/api/workouts/1/group/accept", "", "bob")
	fmt.Println(w.Code)
	// Output: 200
}
