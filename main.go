package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"gesitr/internal/database"
	"gesitr/internal/handlers"
	"gesitr/internal/models"

	"github.com/gin-gonic/gin"
)

//go:embed web/dist/browser/*
var staticFiles embed.FS

func main() {
	database.Init()
	database.DB.AutoMigrate(&models.Item{})

	r := gin.Default()

	// API routes
	api := r.Group("/api")
	{
		api.GET("/items", handlers.ListItems)
		api.POST("/items", handlers.CreateItem)
		api.GET("/items/:id", handlers.GetItem)
		api.PUT("/items/:id", handlers.UpdateItem)
		api.DELETE("/items/:id", handlers.DeleteItem)
	}

	// Serve SPA in production (skip when developing with ng serve)
	if os.Getenv("DEV") != "true" {
		distFS, err := fs.Sub(staticFiles, "web/dist/browser")
		if err != nil {
			log.Fatal("Failed to load embedded files:", err)
		}
		r.NoRoute(func(c *gin.Context) {
			// Try to serve the requested file
			f, err := http.FS(distFS).Open(c.Request.URL.Path)
			if err == nil {
				f.Close()
				c.FileFromFS(c.Request.URL.Path, http.FS(distFS))
				return
			}
			// Fall back to index.html for SPA client-side routing
			c.FileFromFS("index.html", http.FS(distFS))
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
