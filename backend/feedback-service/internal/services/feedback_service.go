package services

import (
	"database/sql"
	"github.com/xmashedxtomatox/feedback-service/internal/models"
)

type FeedbackService struct {
	db *sql.DB
}

func NewFeedbackService(db *sql.DB) *FeedbackService {
	return &FeedbackService{db: db}
}

func (s *FeedbackService) CreateFeedback(userID int, message string) (models.Feedback, error) {
	var fb models.Feedback
	err := s.db.QueryRow(
		`INSERT INTO feedback (user_id, message) VALUES ($1, $2) RETURNING id, user_id, message, created_at`,
		userID, message,
	).Scan(&fb.ID, &fb.UserID, &fb.Message, &fb.CreatedAt)
	return fb, err
}

func (s *FeedbackService) GetFeedbacks(userID int) ([]models.Feedback, error) {
	rows, err := s.db.Query(`SELECT id, user_id, message, created_at FROM feedback WHERE user_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []models.Feedback
	for rows.Next() {
		var fb models.Feedback
		if err := rows.Scan(&fb.ID, &fb.UserID, &fb.Message, &fb.CreatedAt); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, fb)
	}
	return feedbacks, nil
}
