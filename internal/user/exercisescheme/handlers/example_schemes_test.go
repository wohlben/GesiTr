package handlers

import (
	"encoding/json"
	"fmt"

	"gesitr/internal/user/exercisescheme/models"
)

// Creating a rep-based exercise scheme for bicep curls. The exercise must
// exist first (see [CreateExercise]). The scheme defines how to perform
// the exercise: 3 sets of 12 reps at 15kg with 90s rest.
func ExampleCreateExerciseScheme_repBased() {
	setupExampleDB()
	r := newExampleRouter()

	// Create the exercise first.
	doRaw(r, "POST", "/api/exercises", `{
		"names": ["Bicep Curl"],
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0,
		"description": "Dumbbell bicep curl"
	}`)

	// Create a rep-based scheme for this exercise.
	w := doRaw(r, "POST", "/api/user/exercise-schemes", `{
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
	r := newExampleRouter()

	// Create the exercise first.
	doRaw(r, "POST", "/api/exercises", `{
		"names": ["Ergometer"],
		"type": "CARDIO",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0,
		"description": "Rowing ergometer"
	}`)

	// Create a time-based scheme for this exercise.
	w := doRaw(r, "POST", "/api/user/exercise-schemes", `{
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
	r := newExampleRouter()

	// Create the exercise and a scheme.
	doRaw(r, "POST", "/api/exercises", `{
		"names": ["Squat"],
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0.5,
		"description": "Barbell squat"
	}`)
	doRaw(r, "POST", "/api/user/exercise-schemes", `{
		"exerciseId": 1,
		"measurementType": "REP_BASED",
		"sets": 5,
		"reps": 5,
		"weight": 100.0,
		"restBetweenSets": 180
	}`)

	w := doJSON(r, "GET", "/api/user/exercise-schemes/1", nil)

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
	r := newExampleRouter()

	// Create a public exercise and a scheme for it.
	doRaw(r, "POST", "/api/exercises", `{
		"names": ["Push-up"],
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 1.0,
		"description": "Bodyweight push-up",
		"public": true
	}`)
	doRaw(r, "POST", "/api/user/exercise-schemes", `{
		"exerciseId": 1,
		"measurementType": "REP_BASED",
		"sets": 3,
		"reps": 20
	}`)

	// Another user can read the scheme because the exercise is public.
	w := doRawAs(r, "GET", "/api/user/exercise-schemes/1", "", "other")

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
	r := newExampleRouter()

	// Create a private exercise and a scheme for it.
	doRaw(r, "POST", "/api/exercises", `{
		"names": ["Secret Move"],
		"type": "STRENGTH",
		"technicalDifficulty": "advanced",
		"bodyWeightScaling": 0,
		"description": "A private exercise"
	}`)
	doRaw(r, "POST", "/api/user/exercise-schemes", `{
		"exerciseId": 1,
		"measurementType": "REP_BASED",
		"sets": 3,
		"reps": 8
	}`)

	// Another user is denied because the exercise is private.
	w := doRawAs(r, "GET", "/api/user/exercise-schemes/1", "", "other")
	fmt.Println(w.Code)
	// Output: 403
}

// Owner can update their exercise scheme.
func ExampleUpdateExerciseScheme_owner() {
	setupExampleDB()
	r := newExampleRouter()

	doRaw(r, "POST", "/api/exercises", `{
		"names": ["Bicep Curl"],
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0,
		"description": "Dumbbell bicep curl"
	}`)
	doRaw(r, "POST", "/api/user/exercise-schemes", `{
		"exerciseId": 1,
		"measurementType": "REP_BASED",
		"sets": 3,
		"reps": 10,
		"weight": 12.0,
		"restBetweenSets": 60
	}`)

	// Update the scheme — increase weight and reduce reps.
	w := doRaw(r, "PUT", "/api/user/exercise-schemes/1", `{
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
	r := newExampleRouter()

	doRaw(r, "POST", "/api/exercises", `{
		"names": ["Push-up"],
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 1.0,
		"description": "Bodyweight push-up",
		"public": true
	}`)
	doRaw(r, "POST", "/api/user/exercise-schemes", `{
		"exerciseId": 1,
		"measurementType": "REP_BASED",
		"sets": 3,
		"reps": 20
	}`)

	w := doRawAs(r, "PUT", "/api/user/exercise-schemes/1", `{
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
	r := newExampleRouter()

	doRaw(r, "POST", "/api/exercises", `{
		"names": ["Secret Move"],
		"type": "STRENGTH",
		"technicalDifficulty": "advanced",
		"bodyWeightScaling": 0,
		"description": "A private exercise"
	}`)
	doRaw(r, "POST", "/api/user/exercise-schemes", `{
		"exerciseId": 1,
		"measurementType": "REP_BASED",
		"sets": 3,
		"reps": 8
	}`)

	w := doRawAs(r, "PUT", "/api/user/exercise-schemes/1", `{
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
	r := newExampleRouter()

	// Create a private exercise and a scheme.
	doRaw(r, "POST", "/api/exercises", `{
		"names": ["Squat"],
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 0.5,
		"description": "Barbell squat"
	}`)
	doRaw(r, "POST", "/api/user/exercise-schemes", `{
		"exerciseId": 1,
		"measurementType": "REP_BASED",
		"sets": 5,
		"reps": 5
	}`)

	w := doJSON(r, "GET", "/api/user/exercise-schemes", nil)

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

// Non-owner does not see schemes for public exercises unless they share
// an ownership group with the scheme owner.
func ExampleListExerciseSchemes_nonOwnerPublicExercise() {
	setupExampleDB()
	r := newExampleRouter()

	// Create a public exercise and a scheme.
	doRaw(r, "POST", "/api/exercises", `{
		"names": ["Push-up"],
		"type": "STRENGTH",
		"technicalDifficulty": "beginner",
		"bodyWeightScaling": 1.0,
		"description": "Bodyweight push-up",
		"public": true
	}`)
	doRaw(r, "POST", "/api/user/exercise-schemes", `{
		"exerciseId": 1,
		"measurementType": "REP_BASED",
		"sets": 3,
		"reps": 20
	}`)

	// Another user does not see the scheme — not in the same ownership group.
	w := doRawAs(r, "GET", "/api/user/exercise-schemes", "", "other")

	var schemes []models.ExerciseScheme
	json.Unmarshal(w.Body.Bytes(), &schemes)
	fmt.Println(w.Code)
	fmt.Println(len(schemes))
	// Output:
	// 200
	// 0
}

// Non-owner does not see schemes for private exercises.
func ExampleListExerciseSchemes_nonOwnerPrivateExercise() {
	setupExampleDB()
	r := newExampleRouter()

	// Create a private exercise and a scheme.
	doRaw(r, "POST", "/api/exercises", `{
		"names": ["Secret Move"],
		"type": "STRENGTH",
		"technicalDifficulty": "advanced",
		"bodyWeightScaling": 0,
		"description": "A private exercise"
	}`)
	doRaw(r, "POST", "/api/user/exercise-schemes", `{
		"exerciseId": 1,
		"measurementType": "REP_BASED",
		"sets": 3,
		"reps": 8
	}`)

	// Another user sees an empty list.
	w := doRawAs(r, "GET", "/api/user/exercise-schemes", "", "other")

	var schemes []models.ExerciseScheme
	json.Unmarshal(w.Body.Bytes(), &schemes)
	fmt.Println(w.Code)
	fmt.Println(len(schemes))
	// Output:
	// 200
	// 0
}
