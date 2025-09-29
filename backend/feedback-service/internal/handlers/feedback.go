package handlers

import (
	"encoding/json"
	"github.com/xmashedxtomatox/feedback-service/internal/middleware"
	"github.com/xmashedxtomatox/feedback-service/internal/services"
	"net/http"
)

type FeedbackHandler struct {
	service *services.FeedbackService
}

func NewFeedbackHandler(service *services.FeedbackService) *FeedbackHandler {
	return &FeedbackHandler{service: service}
}

type FeedbackResponse struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

func (h *FeedbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		h.createFeedback(w, r)
	case "GET":
		h.getFeedbacks(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *FeedbackHandler) createFeedback(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID).(int)

	var req struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	fb, err := h.service.CreateFeedback(userID, req.Message)
	if err != nil {
		http.Error(w, "could not insert feedback", http.StatusInternalServerError)
		return
	}

	resp := FeedbackResponse{
		ID:        fb.ID,
		UserID:    fb.UserID,
		Message:   fb.Message,
		CreatedAt: fb.CreatedAt.String(),
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *FeedbackHandler) getFeedbacks(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID).(int)

	feedbacks, err := h.service.GetFeedbacks(userID)
	if err != nil {
		http.Error(w, "could not fetch feedbacks", http.StatusInternalServerError)
		return
	}

	var responses []FeedbackResponse
	for _, fb := range feedbacks {
		responses = append(responses, FeedbackResponse{
			ID:        fb.ID,
			UserID:    fb.UserID,
			Message:   fb.Message,
			CreatedAt: fb.CreatedAt.String(),
		})
	}

	json.NewEncoder(w).Encode(responses)
}
