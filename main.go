package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"gesitr/internal/auth"
	compEquipmentHandlers "gesitr/internal/compendium/equipment/handlers"
	compEquipmentModels "gesitr/internal/compendium/equipment/models"
	compFulfillmentHandlers "gesitr/internal/compendium/equipmentfulfillment/handlers"
	compFulfillmentModels "gesitr/internal/compendium/equipmentfulfillment/models"
	compExerciseHandlers "gesitr/internal/compendium/exercise/handlers"
	compExerciseModels "gesitr/internal/compendium/exercise/models"
	compGroupHandlers "gesitr/internal/compendium/exercisegroup/handlers"
	compGroupModels "gesitr/internal/compendium/exercisegroup/models"
	compRelationshipHandlers "gesitr/internal/compendium/exerciserelationship/handlers"
	compRelationshipModels "gesitr/internal/compendium/exerciserelationship/models"
	"gesitr/internal/database"
	"gesitr/internal/profile"
	profileHandlers "gesitr/internal/profile/handlers"
	profileModels "gesitr/internal/profile/models"
	userEquipmentHandlers "gesitr/internal/user/equipment/handlers"
	userEquipmentModels "gesitr/internal/user/equipment/models"
	userExerciseHandlers "gesitr/internal/user/exercise/handlers"
	userExerciseModels "gesitr/internal/user/exercise/models"
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
		&compExerciseModels.ExerciseEntity{},
		&compExerciseModels.ExerciseForce{},
		&compExerciseModels.ExerciseMuscle{},
		&compExerciseModels.ExerciseMeasurementParadigm{},
		&compExerciseModels.ExerciseInstruction{},
		&compExerciseModels.ExerciseImage{},
		&compExerciseModels.ExerciseAlternativeName{},
		&compEquipmentModels.EquipmentEntity{},
		&compExerciseModels.ExerciseEquipment{},
		&compFulfillmentModels.FulfillmentEntity{},
		&compRelationshipModels.ExerciseRelationshipEntity{},
		&compGroupModels.ExerciseGroupEntity{},
		&compGroupModels.ExerciseGroupMemberEntity{},
		&compExerciseModels.ExerciseHistoryEntity{},
		&compEquipmentModels.EquipmentHistoryEntity{},
		&userExerciseModels.UserExerciseEntity{},
		&userExerciseModels.UserExerciseSchemeEntity{},
		&userEquipmentModels.UserEquipmentEntity{},
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
		exercises.GET("", compExerciseHandlers.ListExercises)
		exercises.POST("", compExerciseHandlers.CreateExercise)
		exercises.GET("/:id", compExerciseHandlers.GetExercise)
		exercises.PUT("/:id", compExerciseHandlers.UpdateExercise)
		exercises.DELETE("/:id", compExerciseHandlers.DeleteExercise)
		exercises.GET("/:id/versions", compExerciseHandlers.ListExerciseVersions)
		exercises.GET("/templates/:templateId/versions/:version", compExerciseHandlers.GetExerciseVersion)
	}

	equipment := api.Group("/equipment")
	{
		equipment.GET("", compEquipmentHandlers.ListEquipment)
		equipment.POST("", compEquipmentHandlers.CreateEquipment)
		equipment.GET("/:id", compEquipmentHandlers.GetEquipment)
		equipment.PUT("/:id", compEquipmentHandlers.UpdateEquipment)
		equipment.DELETE("/:id", compEquipmentHandlers.DeleteEquipment)
		equipment.GET("/:id/versions", compEquipmentHandlers.ListEquipmentVersions)
		equipment.GET("/templates/:templateId/versions/:version", compEquipmentHandlers.GetEquipmentVersion)
	}

	fulfillments := api.Group("/fulfillments")
	{
		fulfillments.GET("", compFulfillmentHandlers.ListFulfillments)
		fulfillments.POST("", compFulfillmentHandlers.CreateFulfillment)
		fulfillments.DELETE("/:id", compFulfillmentHandlers.DeleteFulfillment)
	}

	exerciseRelationships := api.Group("/exercise-relationships")
	{
		exerciseRelationships.GET("", compRelationshipHandlers.ListExerciseRelationships)
		exerciseRelationships.POST("", compRelationshipHandlers.CreateExerciseRelationship)
		exerciseRelationships.DELETE("/:id", compRelationshipHandlers.DeleteExerciseRelationship)
	}

	exerciseGroups := api.Group("/exercise-groups")
	{
		exerciseGroups.GET("", compGroupHandlers.ListExerciseGroups)
		exerciseGroups.POST("", compGroupHandlers.CreateExerciseGroup)
		exerciseGroups.GET("/:id", compGroupHandlers.GetExerciseGroup)
		exerciseGroups.PUT("/:id", compGroupHandlers.UpdateExerciseGroup)
		exerciseGroups.DELETE("/:id", compGroupHandlers.DeleteExerciseGroup)
	}

	exerciseGroupMembers := api.Group("/exercise-group-members")
	{
		exerciseGroupMembers.GET("", compGroupHandlers.ListExerciseGroupMembers)
		exerciseGroupMembers.POST("", compGroupHandlers.CreateExerciseGroupMember)
		exerciseGroupMembers.DELETE("/:id", compGroupHandlers.DeleteExerciseGroupMember)
	}

	user := api.Group("/user")

	user.GET("/profile", profileHandlers.GetMyProfile)
	user.PUT("/profile", profileHandlers.UpdateMyProfile)

	userExercises := user.Group("/exercises")
	{
		userExercises.GET("", userExerciseHandlers.ListUserExercises)
		userExercises.POST("", userExerciseHandlers.CreateUserExercise)
		userExercises.GET("/:id", userExerciseHandlers.GetUserExercise)
		userExercises.DELETE("/:id", userExerciseHandlers.DeleteUserExercise)
	}

	userEquipment := user.Group("/equipment")
	{
		userEquipment.GET("", userEquipmentHandlers.ListUserEquipment)
		userEquipment.POST("", userEquipmentHandlers.CreateUserEquipment)
		userEquipment.GET("/:id", userEquipmentHandlers.GetUserEquipment)
		userEquipment.DELETE("/:id", userEquipmentHandlers.DeleteUserEquipment)
	}

	userExerciseSchemes := user.Group("/exercise-schemes")
	{
		userExerciseSchemes.GET("", userExerciseHandlers.ListUserExerciseSchemes)
		userExerciseSchemes.POST("", userExerciseHandlers.CreateUserExerciseScheme)
		userExerciseSchemes.GET("/:id", userExerciseHandlers.GetUserExerciseScheme)
		userExerciseSchemes.PUT("/:id", userExerciseHandlers.UpdateUserExerciseScheme)
		userExerciseSchemes.DELETE("/:id", userExerciseHandlers.DeleteUserExerciseScheme)
	}

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
