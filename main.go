package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"gesitr/internal/database"
	"gesitr/internal/compendium/handlers"
	"gesitr/internal/compendium/models"
	userhandlers "gesitr/internal/user/handlers"
	usermodels "gesitr/internal/user/models"

	"github.com/gin-gonic/gin"
)

//go:embed web/dist/browser/*
var staticFiles embed.FS

func autoMigrate() {
	database.DB.AutoMigrate(
		&models.ExerciseEntity{},
		&models.ExerciseForce{},
		&models.ExerciseMuscle{},
		&models.ExerciseMeasurementParadigm{},
		&models.ExerciseInstruction{},
		&models.ExerciseImage{},
		&models.ExerciseAlternativeName{},
		&models.EquipmentEntity{},
		&models.ExerciseEquipment{},
		&models.FulfillmentEntity{},
		&models.ExerciseRelationshipEntity{},
		&models.ExerciseGroupEntity{},
		&models.ExerciseGroupMemberEntity{},
		&usermodels.UserExerciseEntity{},
		&usermodels.UserEquipmentEntity{},
		&usermodels.UserExerciseSchemeEntity{},
		&usermodels.WorkoutEntity{},
		&usermodels.WorkoutSectionEntity{},
		&usermodels.WorkoutSectionExerciseEntity{},
		&usermodels.WorkoutLogEntity{},
		&usermodels.WorkoutLogSectionEntity{},
		&usermodels.WorkoutLogExerciseEntity{},
		&usermodels.WorkoutLogExerciseSetEntity{},
		&usermodels.UserRecordEntity{},
	)
}

func setupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	exercises := api.Group("/exercises")
	{
		exercises.GET("", handlers.ListExercises)
		exercises.POST("", handlers.CreateExercise)
		exercises.GET("/:id", handlers.GetExercise)
		exercises.PUT("/:id", handlers.UpdateExercise)
		exercises.DELETE("/:id", handlers.DeleteExercise)
	}

	equipment := api.Group("/equipment")
	{
		equipment.GET("", handlers.ListEquipment)
		equipment.POST("", handlers.CreateEquipment)
		equipment.GET("/:id", handlers.GetEquipment)
		equipment.PUT("/:id", handlers.UpdateEquipment)
		equipment.DELETE("/:id", handlers.DeleteEquipment)
	}

	fulfillments := api.Group("/fulfillments")
	{
		fulfillments.GET("", handlers.ListFulfillments)
		fulfillments.POST("", handlers.CreateFulfillment)
		fulfillments.DELETE("/:id", handlers.DeleteFulfillment)
	}

	exerciseRelationships := api.Group("/exercise-relationships")
	{
		exerciseRelationships.GET("", handlers.ListExerciseRelationships)
		exerciseRelationships.POST("", handlers.CreateExerciseRelationship)
		exerciseRelationships.DELETE("/:id", handlers.DeleteExerciseRelationship)
	}

	exerciseGroups := api.Group("/exercise-groups")
	{
		exerciseGroups.GET("", handlers.ListExerciseGroups)
		exerciseGroups.POST("", handlers.CreateExerciseGroup)
		exerciseGroups.GET("/:id", handlers.GetExerciseGroup)
		exerciseGroups.DELETE("/:id", handlers.DeleteExerciseGroup)
	}

	exerciseGroupMembers := api.Group("/exercise-group-members")
	{
		exerciseGroupMembers.GET("", handlers.ListExerciseGroupMembers)
		exerciseGroupMembers.POST("", handlers.CreateExerciseGroupMember)
		exerciseGroupMembers.DELETE("/:id", handlers.DeleteExerciseGroupMember)
	}

	user := api.Group("/user")

	userExercises := user.Group("/exercises")
	{
		userExercises.GET("", userhandlers.ListUserExercises)
		userExercises.POST("", userhandlers.CreateUserExercise)
		userExercises.GET("/:id", userhandlers.GetUserExercise)
		userExercises.DELETE("/:id", userhandlers.DeleteUserExercise)
	}

	userEquipment := user.Group("/equipment")
	{
		userEquipment.GET("", userhandlers.ListUserEquipment)
		userEquipment.POST("", userhandlers.CreateUserEquipment)
		userEquipment.GET("/:id", userhandlers.GetUserEquipment)
		userEquipment.DELETE("/:id", userhandlers.DeleteUserEquipment)
	}

	userExerciseSchemes := user.Group("/exercise-schemes")
	{
		userExerciseSchemes.GET("", userhandlers.ListUserExerciseSchemes)
		userExerciseSchemes.POST("", userhandlers.CreateUserExerciseScheme)
		userExerciseSchemes.GET("/:id", userhandlers.GetUserExerciseScheme)
		userExerciseSchemes.PUT("/:id", userhandlers.UpdateUserExerciseScheme)
		userExerciseSchemes.DELETE("/:id", userhandlers.DeleteUserExerciseScheme)
	}

	workouts := user.Group("/workouts")
	{
		workouts.GET("", userhandlers.ListWorkouts)
		workouts.POST("", userhandlers.CreateWorkout)
		workouts.GET("/:id", userhandlers.GetWorkout)
		workouts.PUT("/:id", userhandlers.UpdateWorkout)
		workouts.DELETE("/:id", userhandlers.DeleteWorkout)
	}

	workoutSections := user.Group("/workout-sections")
	{
		workoutSections.GET("", userhandlers.ListWorkoutSections)
		workoutSections.POST("", userhandlers.CreateWorkoutSection)
		workoutSections.GET("/:id", userhandlers.GetWorkoutSection)
		workoutSections.DELETE("/:id", userhandlers.DeleteWorkoutSection)
	}

	workoutSectionExercises := user.Group("/workout-section-exercises")
	{
		workoutSectionExercises.GET("", userhandlers.ListWorkoutSectionExercises)
		workoutSectionExercises.POST("", userhandlers.CreateWorkoutSectionExercise)
		workoutSectionExercises.DELETE("/:id", userhandlers.DeleteWorkoutSectionExercise)
	}

	workoutLogs := user.Group("/workout-logs")
	{
		workoutLogs.GET("", userhandlers.ListWorkoutLogs)
		workoutLogs.POST("", userhandlers.CreateWorkoutLog)
		workoutLogs.GET("/:id", userhandlers.GetWorkoutLog)
		workoutLogs.PUT("/:id", userhandlers.UpdateWorkoutLog)
		workoutLogs.DELETE("/:id", userhandlers.DeleteWorkoutLog)
	}

	workoutLogSections := user.Group("/workout-log-sections")
	{
		workoutLogSections.GET("", userhandlers.ListWorkoutLogSections)
		workoutLogSections.POST("", userhandlers.CreateWorkoutLogSection)
		workoutLogSections.GET("/:id", userhandlers.GetWorkoutLogSection)
		workoutLogSections.DELETE("/:id", userhandlers.DeleteWorkoutLogSection)
	}

	workoutLogExercises := user.Group("/workout-log-exercises")
	{
		workoutLogExercises.GET("", userhandlers.ListWorkoutLogExercises)
		workoutLogExercises.POST("", userhandlers.CreateWorkoutLogExercise)
		workoutLogExercises.PUT("/:id", userhandlers.UpdateWorkoutLogExercise)
		workoutLogExercises.DELETE("/:id", userhandlers.DeleteWorkoutLogExercise)
	}

	workoutLogExerciseSets := user.Group("/workout-log-exercise-sets")
	{
		workoutLogExerciseSets.GET("", userhandlers.ListWorkoutLogExerciseSets)
		workoutLogExerciseSets.POST("", userhandlers.CreateWorkoutLogExerciseSet)
		workoutLogExerciseSets.PUT("/:id", userhandlers.UpdateWorkoutLogExerciseSet)
		workoutLogExerciseSets.DELETE("/:id", userhandlers.DeleteWorkoutLogExerciseSet)
	}

	records := user.Group("/records")
	{
		records.GET("", userhandlers.ListUserRecords)
		records.GET("/:id", userhandlers.GetUserRecord)
	}
}

func setupSPA(r *gin.Engine) {
	distFS, err := fs.Sub(staticFiles, "web/dist/browser")
	if err != nil {
		log.Fatal("Failed to load embedded files:", err)
	}
	r.NoRoute(func(c *gin.Context) {
		f, err := http.FS(distFS).Open(c.Request.URL.Path)
		if err == nil {
			f.Close()
			c.FileFromFS(c.Request.URL.Path, http.FS(distFS))
			return
		}
		c.FileFromFS("index.html", http.FS(distFS))
	})
}

func buildApp() *gin.Engine {
	database.Init()
	autoMigrate()

	r := gin.Default()
	setupRoutes(r)

	if os.Getenv("DEV") != "true" {
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
