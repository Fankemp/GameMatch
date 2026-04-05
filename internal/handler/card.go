package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Fankemp/GameMatch/internal/model"
	"github.com/Fankemp/GameMatch/internal/service"
	"github.com/gin-gonic/gin"
)

type CardHandler struct {
	cardService service.CardService
}

func NewCardHandler(cardService service.CardService) *CardHandler {
	return &CardHandler{cardService: cardService}
}

func (h *CardHandler) Create(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var input service.CreateCardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if input.GameID == "" || input.Rank == "" || input.Role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "game_id, rank and role are required"})
		return
	}

	card, err := h.cardService.Create(c.Request.Context(), userID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create card"})
		return
	}

	c.JSON(http.StatusCreated, card)
}

func (h *CardHandler) GetMyCards(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	cards, err := h.cardService.GetMyCards(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cards"})
		return
	}

	if cards == nil {
		cards = []*model.GameCard{}
	}
	c.JSON(http.StatusOK, cards)
}

func (h *CardHandler) Update(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	cardID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card id"})
		return
	}

	var input service.UpdateCardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	card, err := h.cardService.Update(c.Request.Context(), userID, cardID, input)
	if err != nil {
		if errors.Is(err, service.ErrCardNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "card not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update card"})
		return
	}

	c.JSON(http.StatusOK, card)
}

func (h *CardHandler) Delete(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	cardID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card id"})
		return
	}

	if err := h.cardService.Delete(c.Request.Context(), userID, cardID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete card"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "card deleted"})
}
