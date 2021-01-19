package config

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	ozzo "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

type jsendSuccess struct {
	Code interface{} `json:"code"`
	Data interface{} `json:"data"`
}

type JSendError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Err     error       `json:"-"`
}

// LoadEnv loads environment variables into env
func LoadEnv(env interface{}) error {
	err := godotenv.Load()
	if err != nil {
		perr, ok := err.(*os.PathError)
		if !ok || !errors.Is(perr.Unwrap(), os.ErrNotExist) {
			return err
		}
	}

	return envconfig.Process("", env)
}

// WithCancel replicates context.WithCancel but listens for Interrupt and SIGTERM signals
func WithCancel(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	go func() {
		defer cancel()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		<-quit
	}()

	return ctx, cancel
}

func ReadJSON(r *http.Request, v interface{}) {
	contentType := r.Header.Get("Content-Type")

	if !strings.Contains(contentType, "application/json") {
		panic(JSendError{
			Code:    http.StatusUnsupportedMediaType,
			Message: http.StatusText(http.StatusUnsupportedMediaType),
		},
		)
	}

	err := json.NewDecoder(r.Body).Decode(v)
	switch {
	case err == io.EOF:
		err := ozzo.Validate(v)
		if err != nil {
			panic(JSendError{
				Code:    http.StatusUnprocessableEntity,
				Message: "We could not validate your request.",
			},
			)
		}
	case err != nil:
		panic(JSendError{
			Code:    http.StatusBadRequest,
			Message: "We cannot parse your request body.",
		},
		)
	default:
		err := ozzo.Validate(v)
		if err != nil {
			panic(JSendError{
				Code:    http.StatusUnprocessableEntity,
				Message: "We could not validate your request.",
			},
			)
		}
	}
}

func Send(w http.ResponseWriter, code int, data []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(code)
	if _, err := w.Write(data); err != nil {
		panic(err.Error())
	}
}

func SendSuccess(w http.ResponseWriter, r *http.Request, v interface{}) {
	log := zerolog.Ctx(r.Context())
	raw := getJSON(log, jsendSuccess{http.StatusOK, v})

	log.Info().Msg("")

	Send(w, http.StatusOK, raw)
}

func SendError(w http.ResponseWriter, r *http.Request, err JSendError) {
	log := zerolog.Ctx(r.Context())
	raw := getJSON(log, err)

	log.Err(err).Msg("")

	Send(w, err.Code, raw)
}

func getJSON(log *zerolog.Logger, v interface{}) []byte {
	raw, _ := json.Marshal(v)

	if v != nil {
		log.UpdateContext(func(ctx zerolog.Context) zerolog.Context {
			buffer := new(bytes.Buffer)

			if err := json.Compact(buffer, raw); err != nil {
				panic(err.Error())
			}

			return ctx.RawJSON("response", buffer.Bytes())
		})
	}

	return raw
}

func (e JSendError) Error() string {
	if e.Err == nil {
		return e.Message
	}

	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}
