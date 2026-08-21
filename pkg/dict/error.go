package dict

import (
	"encoding/json"
	"net/http"
)

var (
	ErrNotFound     = "未找到记录"
	ErrInternal     = "内部错误"
	ErrBadRequest   = "请求错误"
	ErrUnauthorized = "未授权"
	ErrForbidden    = "禁止访问"
)

func HandleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	result := map[string]string{"error": err.Error()}
	jsonResult, _ := json.Marshal(result)
	w.Write(jsonResult)
}

func HandleNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	result := map[string]string{"error": ErrNotFound}
	jsonResult, _ := json.Marshal(result)
	w.Write(jsonResult)
}

func HandleInternal(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	result := map[string]string{"error": ErrInternal}
	jsonResult, _ := json.Marshal(result)
	w.Write(jsonResult)
}

func HandleBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	result := map[string]string{"error": ErrBadRequest}
	jsonResult, _ := json.Marshal(result)
	w.Write(jsonResult)
}

func HandleUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	result := map[string]string{"error": ErrUnauthorized}
	jsonResult, _ := json.Marshal(result)
	w.Write(jsonResult)
}

func HandleForbidden(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	result := map[string]string{"error": ErrForbidden}
	jsonResult, _ := json.Marshal(result)
	w.Write(jsonResult)
}
