package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/AksAman/gophercises/quietHN/settings"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("RECOVER", r)
					msg := "Something went wrong"
					w.Header().Set("Content-Type", "text/html")

					if settings.Settings.Debug {
						stackTrace := debug.Stack()
						msg = fmt.Sprintf("\n<h1 style='color:red'>DEBUG MODE ON</h1> <h3>%s</h3> <h3>Error Message: %s</h3> <h4>Stacktrace</h4><hr> <pre>%s</pre>\n", msg, r, stackTrace)

						http.Error(w, msg, http.StatusInternalServerError)
					} else {
						msg = fmt.Sprintf("\n<h1>%s</h1>\n", msg)
						http.Error(w, msg, http.StatusInternalServerError)
					}
				}
			}()
			rrw := &recoveryResponseWriter{ResponseWriter: w}
			next.ServeHTTP(rrw, r)
			rrw.flushToOriginalWriter()
		},
	)
}

type recoveryResponseWriter struct {
	http.ResponseWriter
	status int
	writes [][]byte
}

func (rrw *recoveryResponseWriter) WriteHeader(statusCode int) {
	fmt.Println("--- rrw WriteHeader ---", statusCode)
	rrw.status = statusCode
}

// Write writes the data to the connection as part of an HTTP reply.
func (rrw *recoveryResponseWriter) Write(b []byte) (int, error) {
	rrw.writes = append(rrw.writes, b)

	fmt.Println("--- rrw Write ---")
	fmt.Printf("rrw.ResponseWriter.Header(): %v\n", rrw.ResponseWriter.Header())

	return len(b), nil
}

func (rrw *recoveryResponseWriter) flushToOriginalWriter() error {
	fmt.Println("--- flushToOriginalWriter ---")
	if rrw.status != 0 {
		rrw.ResponseWriter.WriteHeader(rrw.status)
	}

	for _, write := range rrw.writes {
		_, err := rrw.ResponseWriter.Write(write)
		if err != nil {
			return err
		}
	}

	return nil

}
