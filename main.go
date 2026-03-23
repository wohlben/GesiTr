package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	equipmentHandlers "gesitr/internal/equipment/handlers"
	equipmentModels "gesitr/internal/equipment/models"
	fulfillmentHandlers "gesitr/internal/equipmentfulfillment/handlers"
	fulfillmentModels "gesitr/internal/equipmentfulfillment/models"
	equipmentRelHandlers "gesitr/internal/equipmentrelationship/handlers"
	equipmentRelModels "gesitr/internal/equipmentrelationship/models"
	exerciseHandlers "gesitr/internal/exercise/handlers"
	exerciseModels "gesitr/internal/exercise/models"
	exerciseGroupHandlers "gesitr/internal/exercisegroup/handlers"
	exerciseGroupModels "gesitr/internal/exercisegroup/models"
	exerciseRelHandlers "gesitr/internal/exerciserelationship/handlers"
	exerciseRelModels "gesitr/internal/exerciserelationship/models"
	"gesitr/internal/profile"
	profileHandlers "gesitr/internal/profile/handlers"
	profileModels "gesitr/internal/profile/models"
	exerciseLogHandlers "gesitr/internal/user/exerciselog/handlers"
	exerciseLogModels "gesitr/internal/user/exerciselog/models"
	workoutHandlers "gesitr/internal/user/workout/handlers"
	workoutModels "gesitr/internal/user/workout/models"
	workoutloghandlers "gesitr/internal/user/workoutlog/handlers"
	workoutlogmodels "gesitr/internal/user/workoutlog/models"

	"github.com/gin-gonic/gin"
)

//go:embed web/dist/browser/*
var staticFiles embed.FS

func autoMigrate() {
	database.DB.AutoMigrate(
		&profileModels.UserProfileEntity{},
		&exerciseModels.ExerciseEntity{},
		&exerciseModels.ExerciseForce{},
		&exerciseModels.ExerciseMuscle{},
		&exerciseModels.ExerciseMeasurementParadigm{},
		&exerciseModels.ExerciseInstruction{},
		&exerciseModels.ExerciseImage{},
		&exerciseModels.ExerciseAlternativeName{},
		&equipmentModels.EquipmentEntity{},
		&exerciseModels.ExerciseEquipment{},
		&fulfillmentModels.FulfillmentEntity{},
		&exerciseRelModels.ExerciseRelationshipEntity{},
		&exerciseGroupModels.ExerciseGroupEntity{},
		&exerciseGroupModels.ExerciseGroupMemberEntity{},
		&exerciseModels.ExerciseHistoryEntity{},
		&equipmentModels.EquipmentHistoryEntity{},
		&exerciseModels.ExerciseSchemeEntity{},
		&equipmentRelModels.EquipmentRelationshipEntity{},
		&workoutModels.WorkoutEntity{},
		&workoutModels.WorkoutSectionEntity{},
		&workoutModels.WorkoutSectionExerciseEntity{},
		&workoutlogmodels.WorkoutLogEntity{},
		&workoutlogmodels.WorkoutLogSectionEntity{},
		&workoutlogmodels.WorkoutLogExerciseEntity{},
		&workoutlogmodels.WorkoutLogExerciseSetEntity{},
		&exerciseLogModels.ExerciseLogEntity{},
	)
}

func setupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.Use(auth.UserID())
	api.Use(profile.EnsureProfile())

	api.GET("/profiles/:id", profileHandlers.GetProfile)

	exercises := api.Group("/exercises")
	{
		exercises.GET("", exerciseHandlers.ListExercises)
		exercises.POST("", exerciseHandlers.CreateExercise)
		exercises.GET("/:id", exerciseHandlers.GetExercise)
		exercises.PUT("/:id", exerciseHandlers.UpdateExercise)
		exercises.DELETE("/:id", exerciseHandlers.DeleteExercise)
		exercises.GET("/:id/versions", exerciseHandlers.ListExerciseVersions)
		exercises.GET("/templates/:templateId/versions/:version", exerciseHandlers.GetExerciseVersion)
	}

	equipment := api.Group("/equipment")
	{
		equipment.GET("", equipmentHandlers.ListEquipment)
		equipment.POST("", equipmentHandlers.CreateEquipment)
		equipment.GET("/:id", equipmentHandlers.GetEquipment)
		equipment.PUT("/:id", equipmentHandlers.UpdateEquipment)
		equipment.DELETE("/:id", equipmentHandlers.DeleteEquipment)
		equipment.GET("/:id/versions", equipmentHandlers.ListEquipmentVersions)
		equipment.GET("/templates/:templateId/versions/:version", equipmentHandlers.GetEquipmentVersion)
	}

	exerciseSchemes := api.Group("/exercise-schemes")
	{
		exerciseSchemes.GET("", exerciseHandlers.ListExerciseSchemes)
		exerciseSchemes.POST("", exerciseHandlers.CreateExerciseScheme)
		exerciseSchemes.GET("/:id", exerciseHandlers.GetExerciseScheme)
		exerciseSchemes.PUT("/:id", exerciseHandlers.UpdateExerciseScheme)
		exerciseSchemes.DELETE("/:id", exerciseHandlers.DeleteExerciseScheme)
	}

	fulfillments := api.Group("/fulfillments")
	{
		fulfillments.GET("", fulfillmentHandlers.ListFulfillments)
		fulfillments.POST("", fulfillmentHandlers.CreateFulfillment)
		fulfillments.DELETE("/:id", fulfillmentHandlers.DeleteFulfillment)
	}

	exerciseRelationships := api.Group("/exercise-relationships")
	{
		exerciseRelationships.GET("", exerciseRelHandlers.ListExerciseRelationships)
		exerciseRelationships.POST("", exerciseRelHandlers.CreateExerciseRelationship)
		exerciseRelationships.DELETE("/:id", exerciseRelHandlers.DeleteExerciseRelationship)
	}

	exerciseGroups := api.Group("/exercise-groups")
	{
		exerciseGroups.GET("", exerciseGroupHandlers.ListExerciseGroups)
		exerciseGroups.POST("", exerciseGroupHandlers.CreateExerciseGroup)
		exerciseGroups.GET("/:id", exerciseGroupHandlers.GetExerciseGroup)
		exerciseGroups.PUT("/:id", exerciseGroupHandlers.UpdateExerciseGroup)
		exerciseGroups.DELETE("/:id", exerciseGroupHandlers.DeleteExerciseGroup)
	}

	exerciseGroupMembers := api.Group("/exercise-group-members")
	{
		exerciseGroupMembers.GET("", exerciseGroupHandlers.ListExerciseGroupMembers)
		exerciseGroupMembers.POST("", exerciseGroupHandlers.CreateExerciseGroupMember)
		exerciseGroupMembers.DELETE("/:id", exerciseGroupHandlers.DeleteExerciseGroupMember)
	}

	equipmentRelationships := api.Group("/equipment-relationships")
	{
		equipmentRelationships.GET("", equipmentRelHandlers.ListEquipmentRelationships)
		equipmentRelationships.POST("", equipmentRelHandlers.CreateEquipmentRelationship)
		equipmentRelationships.DELETE("/:id", equipmentRelHandlers.DeleteEquipmentRelationship)
	}

	user := api.Group("/user")

	user.GET("/profile", profileHandlers.GetMyProfile)
	user.PUT("/profile", profileHandlers.UpdateMyProfile)

	workouts := user.Group("/workouts")
	{
		workouts.GET("", workoutHandlers.ListWorkouts)
		workouts.POST("", workoutHandlers.CreateWorkout)
		workouts.GET("/:id", workoutHandlers.GetWorkout)
		workouts.PUT("/:id", workoutHandlers.UpdateWorkout)
		workouts.DELETE("/:id", workoutHandlers.DeleteWorkout)
	}

	workoutSections := user.Group("/workout-sections")
	{
		workoutSections.GET("", workoutHandlers.ListWorkoutSections)
		workoutSections.POST("", workoutHandlers.CreateWorkoutSection)
		workoutSections.GET("/:id", workoutHandlers.GetWorkoutSection)
		workoutSections.DELETE("/:id", workoutHandlers.DeleteWorkoutSection)
	}

	workoutSectionExercises := user.Group("/workout-section-exercises")
	{
		workoutSectionExercises.GET("", workoutHandlers.ListWorkoutSectionExercises)
		workoutSectionExercises.POST("", workoutHandlers.CreateWorkoutSectionExercise)
		workoutSectionExercises.DELETE("/:id", workoutHandlers.DeleteWorkoutSectionExercise)
	}

	workoutLogs := user.Group("/workout-logs")
	{
		workoutLogs.GET("", workoutloghandlers.ListWorkoutLogs)
		workoutLogs.POST("", workoutloghandlers.CreateWorkoutLog)
		workoutLogs.GET("/:id", workoutloghandlers.GetWorkoutLog)
		workoutLogs.PATCH("/:id", workoutloghandlers.UpdateWorkoutLog)
		workoutLogs.DELETE("/:id", workoutloghandlers.DeleteWorkoutLog)
		workoutLogs.POST("/adhoc", workoutloghandlers.StartAdhocWorkoutLog)
		workoutLogs.POST("/:id/start", workoutloghandlers.StartWorkoutLog)
		workoutLogs.POST("/:id/abandon", workoutloghandlers.AbandonWorkoutLog)
		workoutLogs.POST("/:id/finish", workoutloghandlers.FinishWorkoutLog)
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

	exerciseLogs := user.Group("/exercise-logs")
	{
		exerciseLogs.GET("", exerciseLogHandlers.ListExerciseLogs)
		exerciseLogs.POST("", exerciseLogHandlers.CreateExerciseLog)
		exerciseLogs.GET("/:id", exerciseLogHandlers.GetExerciseLog)
		exerciseLogs.PATCH("/:id", exerciseLogHandlers.UpdateExerciseLog)
		exerciseLogs.DELETE("/:id", exerciseLogHandlers.DeleteExerciseLog)
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
