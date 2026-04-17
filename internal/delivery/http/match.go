package http

import (
	"net/http"

	"github.com/Fankemp/GameMatch/internal/model"
	"github.com/Fankemp/GameMatch/internal/repository"
	"github.com/gin-gonic/gin"
)

type MatchHandler struct {
	matchRepo repository.MatchRepository
}

func NewMatchHandler(matchRepo repository.MatchRepository) *MatchHandler {
	return &MatchHandler{matchRepo: matchRepo}
}

func (h *MatchHandler) GetMatches(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	matches, err := h.matchRepo.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get matches"})
		return
	}

	if matches == nil {
		matches = []*model.MatchWithUser{}
	}
	c.JSON(http.StatusOK, matches)
}
