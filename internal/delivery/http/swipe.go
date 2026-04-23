package http

import (
	"net/http"

	"github.com/Fankemp/GameMatch/internal/service"
	"github.com/gin-gonic/gin"
)

type SwipeHandler struct {
	swipeService service.SwipeService
}

func NewSwipeHandler(swipeService service.SwipeService) *SwipeHandler {
	return &SwipeHandler{swipeService: swipeService}
}

func (h *SwipeHandler) Swipe(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var input service.SwipeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if input.TargetCardID == 0 || (input.Action != "like" && input.Action != "dislike") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "target_card_id and valid action (like/dislike) are required"})
		return
	}

	result, err := h.swipeService.Swipe(c.Request.Context(), userID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to swipe"})
		return
	}

	c.JSON(http.StatusOK, result)
}
