package unit

import (
	"testing"
	"time"

	"github.com/xmashedxtomatox/feedback-service/internal/services"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateFeedback(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	service := services.NewFeedbackService(db)

	// Expected DB interaction
	rows := sqlmock.NewRows([]string{"id", "user_id", "message", "created_at"}).
		AddRow(1, 123, "Great app!", time.Now())

	mock.ExpectQuery("INSERT INTO feedback").
		WithArgs(123, "Great app!").
		WillReturnRows(rows)

	fb, err := service.CreateFeedback(123, "Great app!")
	assert.NoError(t, err)
	assert.Equal(t, 123, fb.UserID)
	assert.Equal(t, "Great app!", fb.Message)
}

func TestGetFeedbacks(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	service := services.NewFeedbackService(db)

	createdAt := time.Now()
	rows := sqlmock.NewRows([]string{"id", "user_id", "message", "created_at"}).
		AddRow(1, 123, "Hello world", createdAt).
		AddRow(2, 123, "Second feedback", createdAt)

	mock.ExpectQuery("SELECT id, user_id, message, created_at FROM feedback").
		WithArgs(123).
		WillReturnRows(rows)

	feedbacks, err := service.GetFeedbacks(123)
	assert.NoError(t, err)
	assert.Len(t, feedbacks, 2)
	assert.Equal(t, "Hello world", feedbacks[0].Message)
	assert.Equal(t, "Second feedback", feedbacks[1].Message)
}
