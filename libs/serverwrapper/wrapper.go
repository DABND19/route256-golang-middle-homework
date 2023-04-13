package serverwrapper

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"route256/libs/logger"

	"go.uber.org/zap"
)

type Validator interface {
	Validate() error
}

type Wrapper[Req Validator, Res any] struct {
	handler func(ctx context.Context, req Req) (Res, error)
}

func New[Req Validator, Res any](handler func(ctx context.Context, req Req) (Res, error)) *Wrapper[Req, Res] {
	return &Wrapper[Req, Res]{
		handler: handler,
	}
}

type HTTPError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"code"`
}

func (e HTTPError) Error() string {
	return e.Message
}

type ErrorResponsePayload struct {
	Error HTTPError `json:"error"`
}

func (wrapper *Wrapper[Req, Res]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError := HTTPError{
			Message:    "Method " + r.Method + " is not allowed.",
			StatusCode: http.StatusMethodNotAllowed,
		}
		err := writeError(w, httpError)
		if err != nil {
			logger.Error("Failed to encode error payload.", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	var reqPayload Req
	err := json.NewDecoder(r.Body).Decode(&reqPayload)
	if err != nil {
		httpError := HTTPError{
			Message:    "Invalid request payload: " + err.Error(),
			StatusCode: http.StatusBadRequest,
		}
		err := writeError(w, httpError)
		if err != nil {
			logger.Error("Failed to encode error payload.", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	err = reqPayload.Validate()
	if err != nil {
		httpError := HTTPError{
			Message:    "Invalid request payload: " + err.Error(),
			StatusCode: http.StatusUnprocessableEntity,
		}
		err := writeError(w, httpError)
		if err != nil {
			logger.Error("Failed to encode error payload.", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	ctx := r.Context()
	resPayload, err := wrapper.handler(ctx, reqPayload)
	if err != nil {
		var httpError HTTPError
		if errors.As(err, &httpError) {
			err := writeError(w, httpError)
			if err != nil {
				logger.Error("Failed to encode error payload.", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		} else {
			logger.Error("Internal server error.", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	err = writeJSON(w, resPayload, http.StatusOK)
	if err != nil {
		logger.Error("Failed to encode response payload.", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func writeJSON[Payload any](w http.ResponseWriter, payload Payload, statusCode int) error {
	resBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(resBody)
	return nil
}

func writeError(w http.ResponseWriter, httpError HTTPError) error {
	resPayload := ErrorResponsePayload{
		Error: httpError,
	}
	return writeJSON(w, resPayload, httpError.StatusCode)
}
