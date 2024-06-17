package jsonutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const BadlyFormedJSONErrorMsg = "body contains badly-formed JSON"
const EmptyBodyJSONErrorMsg = "body must not be empty"

type H map[string]any

func ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	err := json.NewDecoder(r.Body).Decode(dst)
	if err == nil {
		return nil
	}

	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var invalidUnmarshalError *json.InvalidUnmarshalError

	switch {
	case errors.As(err, &syntaxError):
		return BadlyFormedJSONError(syntaxError.Offset)
	case errors.Is(err, io.ErrUnexpectedEOF):
		return errors.New(BadlyFormedJSONErrorMsg)
	case errors.As(err, &unmarshalTypeError):
		if unmarshalTypeError.Field != "" {
			return FieldTypeError(unmarshalTypeError.Field)
		}
		return IncorrectJSONTypeError(unmarshalTypeError.Offset)
	case errors.Is(err, io.EOF):
		return errors.New(EmptyBodyJSONErrorMsg)
	case errors.As(err, &invalidUnmarshalError):
		panic(err)
	default:
		return err
	}
}

func WriteJSON(w http.ResponseWriter, status int, data H, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	//  Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func BadlyFormedJSONError(offset int64) error {
	return fmt.Errorf("body contains badly-formed JSON (at character %d)", offset)
}

func FieldTypeError(field string) error {
	return fmt.Errorf("body contains incorrect JSON type for field %q", field)
}

func IncorrectJSONTypeError(offset int64) error {
	return fmt.Errorf("body contains incorrect JSON type (at character %d)", offset)
}
