package tasks

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"strconv"
	"strings"
)

// MethodNotAllowed is simple util function which writes MethodNotAllowed
// error response to given ResponseWriter.
func methodNotAllowed(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte(`{"error":"method not allowed"}`))
}

// Options is simple util function which writes CORS response to given
// responseWriter.
func options(w http.ResponseWriter, r *http.Request) {
	origin := "*"
	if ro := r.Header.Get("Origin"); ro != "" {
		origin = ro
	}

	w.Header().Add("Access-Control-Allow-Origin", origin)
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
}

// ParseTaskIDPath parses request URL and returns slice of TaskIDs or error
// if url is not valid.
func parseTaskIDPath(r *http.Request) ([]TaskID, error) {
	taskIDs := []TaskID{}

	parts := []string{}
	for _, part := range strings.Split(r.URL.Path, "/") {
		if len(part) > 0 {
			parts = append(parts, part)
		}
	}

	// it must have at least 1 part "/tasks"
	if len(parts) < 1 {
		return nil, ErrHandlerURLNotValid
	}

	// Tasks endpoint - skip first 2 values for part "/tasks/"
	for _, value := range parts[1:] {
		if len(value) > 0 {
			val, err := strconv.Atoi(value)
			if err != nil {
				log.Printf("(DEBUG) http: parsing TaskID path on value %q failed: %s\n", value, err)
				return nil, ErrHandlerURLNotValid
			}

			taskIDs = append(taskIDs, TaskID(val))
		}
	}

	return taskIDs, nil
}

var (
	// ErrBadMediaType is returned when request contains not supported
	// Content-Type.
	ErrBadMediaType error = errors.New("Bad media type")
)

// ParseBody parses a request body into an interface. Currently it supports
// only application/json Content-Type.
func parseBody(r *http.Request, v interface{}) error {
	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		log.Printf("(DEBUG) http: parsing request Content-Type %q failed: %s\n", contentType, err)
		return ErrBadMediaType
	}

	switch mediaType {
	case "application/json":
		dec := json.NewDecoder(r.Body)
		for {
			if err := dec.Decode(v); err == io.EOF {
				break
			} else if err != nil {
				log.Printf("(DEBUG) http: parsing request failed: %s\n", err)
				return err
			}
		}
		return nil

	default:
		log.Printf("(DEBUG) http: unsupported Content-Type %q\n", mediaType)
		return ErrBadMediaType
	}

	return nil
}

// JSONError is wrapper struct for errors thrown in business logic and returned
// in request response as JSON.
type JSONError struct {
	Error string `json:"error"`
}

// ErrorAsJSON is simple util funtion which returns given error as JSON
// payload with given status code.
func ErrorAsJSON(w http.ResponseWriter, statusCode int, err error) {
	ResponseAsJSON(w, statusCode, JSONError{err.Error()})
}

// ResponseOK is simple util function which returns given struct as JSON
// payload with status code OK (200).
func ResponseOK(w http.ResponseWriter, v interface{}) {
	ResponseAsJSON(w, http.StatusOK, v)
}

// ResponseCreated is simple util function which returns given struct as JSON
// paylaod with status code Created (201) and set's location header.
func ResponseCreated(w http.ResponseWriter, url string, v interface{}) {
	w.Header().Set("Location", url)
	ResponseAsJSON(w, http.StatusCreated, v)
}

// ResponseAsJSON encodes an interface into JSON.
func ResponseAsJSON(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	b, err := json.Marshal(v)
	if err != nil {
		log.Printf("(WARN) http: marshaling response failed: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
		return
	}

	w.WriteHeader(statusCode)
	if _, err := w.Write(b); err != nil {
		log.Printf("(WARN) http: writting JSON payload failed: %s", err)
	}
}
