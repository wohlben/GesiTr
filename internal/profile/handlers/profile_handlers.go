package handlers

import (
	"net/http"

	"gesitr/internal/auth"
	"gesitr/internal/database"
	"gesitr/internal/profile/models"

	"github.com/gin-gonic/gin"
)

func GetMyProfile(c *gin.Context) {
	userID := auth.GetUserID(c)

	var entity models.UserProfileEntity
	if err := database.DB.First(&entity, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	c.JSON(http.StatusOK, entity.ToDTO())
}

func UpdateMyProfile(c *gin.Context) {
	userID := auth.GetUserID(c)

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var entity models.UserProfileEntity
	if err := database.DB.First(&entity, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	entity.Name = req.Name
	if err := database.DB.Save(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, entity.ToDTO())
}

func GetProfile(c *gin.Context) {
	id := c.Param("id")

	var entity models.UserProfileEntity
	if err := database.DB.First(&entity, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	c.JSON(http.StatusOK, entity.ToDTO())
}
