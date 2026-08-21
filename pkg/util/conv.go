package util

import "student/internal/handler"

func ToRenderers[T handler.Renderer](s []T) []handler.Renderer {
	r := make([]handler.Renderer, len(s))
	for i := range s {
		r[i] = s[i]
	}
	return r
}
