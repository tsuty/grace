package internal

import (
	"context"
	"github.com/google/uuid"
	"log"
	"net/http"
	"runtime"
	"time"
)

type handler struct {
	logger *log.Logger
	writer Writer
}

func Handler(logger *log.Logger, writer Writer) http.Handler {
	return &handler{logger, writer}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler = http.HandlerFunc(h.Receive)
	handler = h.logging(handler)
	handler = h.requestID(handler)
	handler = h.recover(handler)
	handler.ServeHTTP(w, r)
}

func (h *handler) Receive(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	if r.Method == http.MethodPost {
		status = http.StatusCreated
	}
	w.WriteHeader(status)

	if err := h.writer.Write(r); err != nil {
		h.logger.Printf("write error (%T) %q\n", err, err.Error())
	}
}

type contextKey int

var requestID contextKey = 1

func (h *handler) requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestID, uuid.NewString())))
	})
}

func (h *handler) logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		next.ServeHTTP(w, r)
		h.logger.Printf("%s %s %s %q %s\n",
			r.Context().Value(requestID),
			r.Method, r.URL.Path, r.UserAgent(), time.Since(t))
	})
}

func (h *handler) recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if ret := recover(); ret != nil {
				h.logger.Println("panic recovered")
				h.logger.Printf("%#v\n", ret)
				for d := 1; d < 100; d++ {
					if pc, file, line, ok := runtime.Caller(d); ok {
						fn := ""
						if f := runtime.FuncForPC(pc); f != nil {
							fn = f.Name()
						}
						h.logger.Printf("\t%02d %s:%d %s\n", d, file, line, fn)
					} else {
						break
					}
				}

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
