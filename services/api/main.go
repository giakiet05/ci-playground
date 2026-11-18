// Command api is a tiny HTTP service used to exercise the shared CI template.
//
// Routes:
//
//	GET /healthz       -> 200 "ok"
//	GET /hello/{name}  -> 200 JSON greeting, 400 when the name is blank
package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/giakiet05/ci-playground/libs/shared"
)

const shutdownTimeout = 5 * time.Second

// greetingResponse is the JSON body of GET /hello/{name}.
type greetingResponse struct {
	Message string `json:"message"`
}

// errorResponse is the JSON body of every error reply.
type errorResponse struct {
	Error string `json:"error"`
}

// newRouter registers the service's routes.
//
// Parameters:
//   - logger: structured logger used by the handlers.
//
// Returns an http.Handler serving every route.
func newRouter(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if _, err := w.Write([]byte("ok")); err != nil {
			logger.Error("failed to write health response", "error", err)
		}
	})
	mux.HandleFunc("GET /hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		message, err := shared.Greet(r.PathValue("name"))
		if err != nil {
			writeJSON(w, logger, http.StatusBadRequest, errorResponse{Error: err.Error()})
			return
		}
		writeJSON(w, logger, http.StatusOK, greetingResponse{Message: message})
	})
	return mux
}

// writeJSON encodes body as JSON with the given status code and logs encoding failures.
//
// Parameters:
//   - w: the response writer.
//   - logger: logger for encoding failures (the status line is already sent by then).
//   - status: HTTP status code.
//   - body: value encoded as the response body.
func writeJSON(w http.ResponseWriter, logger *slog.Logger, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.Error("failed to encode response", "status", status, "error", err)
	}
}

// run starts the HTTP server and blocks until ctx is cancelled, then shuts down gracefully.
//
// Parameters:
//   - ctx: cancelled on SIGINT/SIGTERM.
//   - logger: structured logger.
//   - addr: listen address, e.g. ":8080".
//
// Returns an error when the server fails to start or to shut down cleanly.
func run(ctx context.Context, logger *slog.Logger, addr string) error {
	server := &http.Server{
		Addr:              addr,
		Handler:           newRouter(logger),
		ReadHeaderTimeout: 5 * time.Second,
	}
	serverErr := make(chan error, 1)
	go func() {
		logger.Info("api listening", "addr", addr)
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	case <-ctx.Done():
	}

	logger.Info("api shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return err
	}
	return nil
}

// main wires the logger and signal handling, then runs the server.
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	addr := os.Getenv("API_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, logger, addr); err != nil {
		logger.Error("api stopped with error", "error", err)
		os.Exit(1)
	}
}
