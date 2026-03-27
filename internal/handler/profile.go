package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Fankemp/GameMatch/internal/service"
)

type ProfileHandler struct {
	profileService service.ProfileService
}

func NewProfileHandler(profileService service.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

func (h *ProfileHandler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var input service.CreateProfileInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	profile, err := h.profileService.CreateProfile(r.Context(), userID, input)
	if err != nil {
		if errors.Is(err, service.ErrProfileAlreadyExists) {
			writeError(w, http.StatusConflict, "profile already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create profile")
		return
	}

	writeJSON(w, http.StatusCreated, profile)
}
