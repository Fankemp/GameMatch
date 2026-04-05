package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Fankemp/GameMatch/internal/service"
	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	profileService service.ProfileService
}

func NewProfileHandler(profileService service.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var input service.CreateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	profile, err := h.profileService.CreateProfile(c.Request.Context(), userID, input)
	if err != nil {
		if errors.Is(err, service.ErrProfileAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "profile already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create profile"})
		return
	}

	c.JSON(http.StatusCreated, profile)
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var input service.UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	profile, err := h.profileService.UpdateProfile(c.Request.Context(), userID, input)
	if err != nil {
		if errors.Is(err, service.ErrProfileNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *ProfileHandler) GetMyProfile(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	profile, err := h.profileService.GetMyProfile(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrProfileNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *ProfileHandler) GetProfileByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid profile id"})
		return
	}

	profile, err := h.profileService.GetProfileByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrProfileNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}
