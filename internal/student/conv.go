package student

import (
	"student/api"
	"student/internal/sqlc"
)

func ToStudentResponse(s *sqlc.Student) api.StudentResponse {
	return api.StudentResponse{
		ID:        s.ID,
		Name:      s.Name,
		Gender:    s.Gender,
		Birthday:  s.Birthday,
		Score:     *s.Score,
		Classroom: s.Classroom,
	}
}

func ToStudentResponses(students []sqlc.Student) []api.StudentResponse {
	responses := make([]api.StudentResponse, 0, len(students))
	for _, s := range students {
		responses = append(responses, ToStudentResponse(&s))
	}
	return responses
}
