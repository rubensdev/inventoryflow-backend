package jsonutil

import (
	"fmt"
	"net/http"
)

type Output interface {
	Print(v ...any)
}

type JSONResponse struct {
	output Output
}

func NewJSONResponse(output Output) *JSONResponse {
	return &JSONResponse{
		output: output,
	}
}

func (j JSONResponse) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	j.logError(r, err)

	message := "the server has encountered a problem and could not process your request"
	j.ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func (j JSONResponse) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	j.ErrorResponse(w, r, http.StatusNotFound, message)
}

func (j JSONResponse) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	j.ErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (j JSONResponse) EditConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	j.ErrorResponse(w, r, http.StatusConflict, message)
}

func (j JSONResponse) ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	var data H

	switch v := message.(type) {
	case error:
		data = H{"error": v.Error()}
	default:
		data = H{"error": v}
	}

	err := WriteJSON(w, status, data, nil)
	if err != nil {
		j.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (j JSONResponse) logError(r *http.Request, err error) {
	// TODO: include http method and URL.
	j.output.Print(err)
}
