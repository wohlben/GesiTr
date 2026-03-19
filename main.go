package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"gesitr/internal/auth"
	compEquipment "gesitr/internal/compendium/equipment"
	compFulfillment "gesitr/internal/compendium/equipmentfulfillment"
	compExercise "gesitr/internal/compendium/exercise"
	compGroup "gesitr/internal/compendium/exercisegroup"
	compRelationship "gesitr/internal/compendium/exerciserelationship"
	"gesitr/internal/database"
	userExercise "gesitr/internal/user/exercise"
	"gesitr/internal/user/record"
	"gesitr/internal/user/workout"
	workoutloghandlers "gesitr/internal/user/workoutlog/handlers"
	workoutlogmodels "gesitr/internal/user/workoutlog/models"

	"github.com/gin-gonic/gin"
)

//go:embed web/dist/browser/*
var staticFiles embed.FS

func autoMigrate() {
	database.DB.AutoMigrate(
		&compExercise.ExerciseEntity{},
		&compExercise.ExerciseForce{},
		&compExercise.ExerciseMuscle{},
		&compExercise.ExerciseMeasurementParadigm{},
		&compExercise.ExerciseInstruction{},
		&compExercise.ExerciseImage{},
		&compExercise.ExerciseAlternativeName{},
		&compEquipment.EquipmentEntity{},
		&compExercise.ExerciseEquipment{},
		&compFulfillment.FulfillmentEntity{},
		&compRelationship.ExerciseRelationshipEntity{},
		&compGroup.ExerciseGroupEntity{},
		&compGroup.ExerciseGroupMemberEntity{},
		&compExercise.ExerciseHistoryEntity{},
		&compEquipment.EquipmentHistoryEntity{},
		&userExercise.UserExerciseEntity{},
		&userExercise.UserEquipmentEntity{},
		&userExercise.UserExerciseSchemeEntity{},
		&workout.WorkoutEntity{},
		&workout.WorkoutSectionEntity{},
		&workout.WorkoutSectionExerciseEntity{},
		&workoutlogmodels.WorkoutLogEntity{},
		&workoutlogmodels.WorkoutLogSectionEntity{},
		&workoutlogmodels.WorkoutLogExerciseEntity{},
		&workoutlogmodels.WorkoutLogExerciseSetEntity{},
		&record.UserRecordEntity{},
	)
}

func setupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.Use(auth.UserID())

	exercises := api.Group("/exercises")
	{
		exercises.GET("", compExercise.ListExercises)
		exercises.POST("", compExercise.CreateExercise)
		exercises.GET("/:id", compExercise.GetExercise)
		exercises.PUT("/:id", compExercise.UpdateExercise)
		exercises.DELETE("/:id", compExercise.DeleteExercise)
		exercises.GET("/:id/versions", compExercise.ListExerciseVersions)
		exercises.GET("/templates/:templateId/versions/:version", compExercise.GetExerciseVersion)
	}

	equipment := api.Group("/equipment")
	{
		equipment.GET("", compEquipment.ListEquipment)
		equipment.POST("", compEquipment.CreateEquipment)
		equipment.GET("/:id", compEquipment.GetEquipment)
		equipment.PUT("/:id", compEquipment.UpdateEquipment)
		equipment.DELETE("/:id", compEquipment.DeleteEquipment)
		equipment.GET("/:id/versions", compEquipment.ListEquipmentVersions)
		equipment.GET("/templates/:templateId/versions/:version", compEquipment.GetEquipmentVersion)
	}

	fulfillments := api.Group("/fulfillments")
	{
		fulfillments.GET("", compFulfillment.ListFulfillments)
		fulfillments.POST("", compFulfillment.CreateFulfillment)
		fulfillments.DELETE("/:id", compFulfillment.DeleteFulfillment)
	}

	exerciseRelationships := api.Group("/exercise-relationships")
	{
		exerciseRelationships.GET("", compRelationship.ListExerciseRelationships)
		exerciseRelationships.POST("", compRelationship.CreateExerciseRelationship)
		exerciseRelationships.DELETE("/:id", compRelationship.DeleteExerciseRelationship)
	}

	exerciseGroups := api.Group("/exercise-groups")
	{
		exerciseGroups.GET("", compGroup.ListExerciseGroups)
		exerciseGroups.POST("", compGroup.CreateExerciseGroup)
		exerciseGroups.GET("/:id", compGroup.GetExerciseGroup)
		exerciseGroups.PUT("/:id", compGroup.UpdateExerciseGroup)
		exerciseGroups.DELETE("/:id", compGroup.DeleteExerciseGroup)
	}

	exerciseGroupMembers := api.Group("/exercise-group-members")
	{
		exerciseGroupMembers.GET("", compGroup.ListExerciseGroupMembers)
		exerciseGroupMembers.POST("", compGroup.CreateExerciseGroupMember)
		exerciseGroupMembers.DELETE("/:id", compGroup.DeleteExerciseGroupMember)
	}

	user := api.Group("/user")

	userExercises := user.Group("/exercises")
	{
		userExercises.GET("", userExercise.ListUserExercises)
		userExercises.POST("", userExercise.CreateUserExercise)
		userExercises.GET("/:id", userExercise.GetUserExercise)
		userExercises.DELETE("/:id", userExercise.DeleteUserExercise)
	}

	userEquipment := user.Group("/equipment")
	{
		userEquipment.GET("", userExercise.ListUserEquipment)
		userEquipment.POST("", userExercise.CreateUserEquipment)
		userEquipment.GET("/:id", userExercise.GetUserEquipment)
		userEquipment.DELETE("/:id", userExercise.DeleteUserEquipment)
	}

	userExerciseSchemes := user.Group("/exercise-schemes")
	{
		userExerciseSchemes.GET("", userExercise.ListUserExerciseSchemes)
		userExerciseSchemes.POST("", userExercise.CreateUserExerciseScheme)
		userExerciseSchemes.GET("/:id", userExercise.GetUserExerciseScheme)
		userExerciseSchemes.PUT("/:id", userExercise.UpdateUserExerciseScheme)
		userExerciseSchemes.DELETE("/:id", userExercise.DeleteUserExerciseScheme)
	}

	workouts := user.Group("/workouts")
	{
		workouts.GET("", workout.ListWorkouts)
		workouts.POST("", workout.CreateWorkout)
		workouts.GET("/:id", workout.GetWorkout)
		workouts.PUT("/:id", workout.UpdateWorkout)
		workouts.DELETE("/:id", workout.DeleteWorkout)
	}

	workoutSections := user.Group("/workout-sections")
	{
		workoutSections.GET("", workout.ListWorkoutSections)
		workoutSections.POST("", workout.CreateWorkoutSection)
		workoutSections.GET("/:id", workout.GetWorkoutSection)
		workoutSections.DELETE("/:id", workout.DeleteWorkoutSection)
	}

	workoutSectionExercises := user.Group("/workout-section-exercises")
	{
		workoutSectionExercises.GET("", workout.ListWorkoutSectionExercises)
		workoutSectionExercises.POST("", workout.CreateWorkoutSectionExercise)
		workoutSectionExercises.DELETE("/:id", workout.DeleteWorkoutSectionExercise)
	}

	workoutLogs := user.Group("/workout-logs")
	{
		workoutLogs.GET("", workoutloghandlers.ListWorkoutLogs)
		workoutLogs.POST("", workoutloghandlers.CreateWorkoutLog)
		workoutLogs.GET("/:id", workoutloghandlers.GetWorkoutLog)
		workoutLogs.PATCH("/:id", workoutloghandlers.UpdateWorkoutLog)
		workoutLogs.DELETE("/:id", workoutloghandlers.DeleteWorkoutLog)
		workoutLogs.POST("/:id/start", workoutloghandlers.StartWorkoutLog)
		workoutLogs.POST("/:id/abandon", workoutloghandlers.AbandonWorkoutLog)
	}

	workoutLogSections := user.Group("/workout-log-sections")
	{
		workoutLogSections.GET("", workoutloghandlers.ListWorkoutLogSections)
		workoutLogSections.POST("", workoutloghandlers.CreateWorkoutLogSection)
		workoutLogSections.GET("/:id", workoutloghandlers.GetWorkoutLogSection)
		workoutLogSections.PATCH("/:id", workoutloghandlers.UpdateWorkoutLogSection)
		workoutLogSections.DELETE("/:id", workoutloghandlers.DeleteWorkoutLogSection)
	}

	workoutLogExercises := user.Group("/workout-log-exercises")
	{
		workoutLogExercises.GET("", workoutloghandlers.ListWorkoutLogExercises)
		workoutLogExercises.POST("", workoutloghandlers.CreateWorkoutLogExercise)
		workoutLogExercises.PATCH("/:id", workoutloghandlers.UpdateWorkoutLogExercise)
		workoutLogExercises.DELETE("/:id", workoutloghandlers.DeleteWorkoutLogExercise)
	}

	workoutLogExerciseSets := user.Group("/workout-log-exercise-sets")
	{
		workoutLogExerciseSets.GET("", workoutloghandlers.ListWorkoutLogExerciseSets)
		workoutLogExerciseSets.POST("", workoutloghandlers.CreateWorkoutLogExerciseSet)
		workoutLogExerciseSets.PATCH("/:id", workoutloghandlers.UpdateWorkoutLogExerciseSet)
		workoutLogExerciseSets.DELETE("/:id", workoutloghandlers.DeleteWorkoutLogExerciseSet)
	}

	records := user.Group("/records")
	{
		records.GET("", record.ListUserRecords)
		records.GET("/:id", record.GetUserRecord)
	}
}

func setupSPA(r *gin.Engine) {
	distFS, err := fs.Sub(staticFiles, "web/dist/browser")
	if err != nil {
		log.Fatal("Failed to load embedded files:", err)
	}
	indexHTML, err := fs.ReadFile(distFS, "index.html")
	if err != nil {
		log.Fatal("Failed to read index.html:", err)
	}
	r.NoRoute(func(c *gin.Context) {
		f, err := http.FS(distFS).Open(c.Request.URL.Path)
		if err == nil {
			f.Close()
			c.FileFromFS(c.Request.URL.Path, http.FS(distFS))
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
	})
}

func buildApp() *gin.Engine {
	database.Init()
	autoMigrate()

	r := gin.Default()
	setupRoutes(r)

	if os.Getenv("DEV") == "true" {
		r.POST("/api/ci/reset-db", func(c *gin.Context) {
			database.DB.Exec("PRAGMA foreign_keys = OFF")
			var tables []struct{ Name string }
			database.DB.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'").Scan(&tables)
			for _, t := range tables {
				database.DB.Exec("DELETE FROM " + t.Name)
			}
			database.DB.Exec("DELETE FROM sqlite_sequence")
			database.DB.Exec("PRAGMA foreign_keys = ON")
			c.JSON(http.StatusOK, gin.H{"status": "reset"})
		})
	} else {
		setupSPA(r)
	}
	return r
}

func main() {
	r := buildApp()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
