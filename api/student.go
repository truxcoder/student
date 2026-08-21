package api

import (
	"net/http"
	"student/internal/sqlc"
	"time"

	"github.com/google/uuid"
)

type StudentRequest struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	Birthday  time.Time `json:"birthday"`
	Score     *float64  `json:"score"`
	Classroom string    `json:"classroom"`
}

type StudentResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	Birthday  time.Time `json:"birthday"`
	Score     float64   `json:"score"`
	Classroom string    `json:"classroom"`
}

func (s *StudentRequest) Bind(r *http.Request) error {
	return nil
}

func (s *StudentRequest) ToCreateParams() *sqlc.CreateStudentParams {
	return &sqlc.CreateStudentParams{
		ID:        uuid.NewString(),
		Name:      s.Name,
		Gender:    s.Gender,
		Birthday:  s.Birthday,
		Score:     s.Score,
		Classroom: s.Classroom,
	}
}

func (s *StudentRequest) ToUpdateParams() *sqlc.UpdateStudentParams {
	return &sqlc.UpdateStudentParams{
		ID:        s.ID,
		Name:      s.Name,
		Gender:    s.Gender,
		Birthday:  s.Birthday,
		Score:     s.Score,
		Classroom: s.Classroom,
	}
}

func (s StudentResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// encoder := json.NewEncoder(w)
	// return encoder.Encode(s)
	//
	return nil
}
