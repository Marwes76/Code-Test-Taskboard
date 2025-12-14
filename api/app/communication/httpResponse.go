package communication

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error	string	`json:"error"`
}

// 200 OK
func ResponseOK(w http.ResponseWriter, payload interface{}) {
	ResponseJson(w, http.StatusOK, payload)
}

// 201 Created
func ResponseCreated(w http.ResponseWriter, payload interface{}) {
	ResponseJson(w, http.StatusCreated, payload)
}

// 204 No Content
func ResponseNoContent(w http.ResponseWriter) {
	ResponseJson(w, http.StatusNoContent, nil)
}

// 400 Bad Request
func ResponseBadRequest(w http.ResponseWriter, err error) {
	ResponseErrorJson(w, http.StatusBadRequest, err)
}

// 404 Not Found
func ResponseNotFound(w http.ResponseWriter) {
	ResponseErrorJson(w, http.StatusNotFound, nil)
}

// 500 Internal Server Error
func ResponseInternalServerError(w http.ResponseWriter, err error) {
	ResponseErrorJson(w, http.StatusInternalServerError, err)
}

func ResponseJson(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func ResponseErrorJson(w http.ResponseWriter, status int, err error) {
	if err != nil {
		ResponseJson(w, status, ErrorResponse{
			Error:	err.Error(),
		})
	} else {
		ResponseJson(w, status, nil)
	}
}