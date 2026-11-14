package student

import (
	"student/internal/sqlc"
)

func ToStudentResponse(s *sqlc.Student) StudentResponse {
	return StudentResponse{
		ID:        s.ID,
		Name:      s.Name,
		Gender:    s.Gender,
		Birthday:  s.Birthday,
		Score:     s.Score,
		Classroom: s.Classroom,
	}
}

func ToStudentResponses(students []sqlc.Student) []StudentResponse {
	responses := make([]StudentResponse, 0, len(students))
	for _, s := range students {
		responses = append(responses, ToStudentResponse(&s))
	}
	return responses
}
