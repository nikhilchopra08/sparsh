package api

import (
	"log"
	"net/http"

	"licensebox/internal/auth"
	"licensebox/internal/db"
	"licensebox/internal/models"

	"github.com/gin-gonic/gin"
)

type Router struct {
	DB *db.Database
}

func NewRouter(database *db.Database) *Router {
	return &Router{DB: database}
}

func (r *Router) SetupRoutes() *gin.Engine {
	app := gin.Default()

	app.POST("/activate", r.handleActivate)

	// Internal helper for demo purposes to create a license
	// In production, this would be behind auth and triggered by payment gateway
	app.POST("/internal/create-license", r.handleCreateLicense)

	return app
}

func (r *Router) handleActivate(c *gin.Context) {
	var req models.ActivateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 1. Validate License
	license, err := r.DB.GetLicenseByKey(req.LicenseKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "License not found"})
		return
	}

	if license.Status != "active" {
		c.JSON(http.StatusForbidden, gin.H{"error": "License is no longer active"})
		return
	}

	// 2. Check if already activated on this device
	isActivated, err := r.DB.CheckActivation(license.ID, req.DeviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// 3. Register activation if new
	if !isActivated {
		if err := r.DB.CreateActivation(license.ID, req.DeviceID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register activation"})
			return
		}
	}

	// 4. Generate signed token
	token, err := auth.GenerateToken(license.ID, req.DeviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.ActivateResponse{
		ActivationToken: token,
	})
}

func (r *Router) handleCreateLicense(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required"`
		Key   string `json:"key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	license, err := r.DB.CreateLicense(req.Email, req.Key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Email sent to %s with key %s", req.Email, req.Key)
	c.JSON(http.StatusCreated, license)
}
