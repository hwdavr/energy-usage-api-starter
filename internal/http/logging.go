package http

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const maxLogBody = 1 << 20 // 1MB

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
	buf    bytes.Buffer
}

func (l *loggingResponseWriter) WriteHeader(status int) {
	l.status = status
	l.ResponseWriter.WriteHeader(status)
}

func (l *loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := l.ResponseWriter.Write(b)
	if n > 0 {
		l.size += n
		// keep a copy up to the maxLogBody
		if l.buf.Len() < maxLogBody {
			remaining := maxLogBody - l.buf.Len()
			if remaining < n {
				l.buf.Write(b[:remaining])
			} else {
				l.buf.Write(b[:n])
			}
		}
	}
	return n, err
}

// RequestResponseLogger returns a middleware that logs request and response info.
func RequestResponseLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Read entire request body so we can restore it for handlers.
			var reqBodyCopy []byte
			if r.Body != nil {
				body, err := io.ReadAll(r.Body)
				if err == nil {
					reqBodyCopy = body
				}
				// restore original body for downstream
				r.Body = io.NopCloser(bytes.NewReader(reqBodyCopy))
			}

			// prepare sanitized headers (don't log Authorization value)
			safeHeaders := make(map[string][]string)
			for k, v := range r.Header {
				if strings.ToLower(k) == "authorization" {
					// keep header key but redact value
					safeHeaders[k] = []string{"REDACTED"}
				} else {
					safeHeaders[k] = v
				}
			}

			lrw := &loggingResponseWriter{
				ResponseWriter: w,
				status:         0,
				size:           0,
			}

			next.ServeHTTP(lrw, r)

			duration := time.Since(start)

			// prepare truncated request/response bodies for logs
			var reqBodyLog string
			if len(reqBodyCopy) > 0 {
				if len(reqBodyCopy) > maxLogBody {
					reqBodyLog = string(reqBodyCopy[:maxLogBody]) + "...(truncated)"
				} else {
					reqBodyLog = string(reqBodyCopy)
				}
			}

			var respBodyLog string
			if lrw.buf.Len() > 0 {
				if lrw.buf.Len() > maxLogBody {
					respBodyLog = lrw.buf.String()[:maxLogBody] + "...(truncated)"
				} else {
					respBodyLog = lrw.buf.String()
				}
			}

			// default status if handler didn't write a header
			status := lrw.status
			if status == 0 {
				status = http.StatusOK
			}

			log.Printf(">>> method=%s path=%s remote=%s status=%d size=%d duration=%s headers=%v \n--> reqBodyLen=%d reqBody=%q\n<-- respBodyLen=%d respBody=%q",
				r.Method,
				r.URL.RequestURI(),
				r.RemoteAddr,
				status,
				lrw.size,
				duration.String(),
				safeHeaders,
				len(reqBodyCopy),
				truncateForLog(reqBodyLog),
				lrw.buf.Len(),
				truncateForLog(respBodyLog),
			)
		})
	}
}

// truncateForLog ensures the string is printable and trims newlines.
func truncateForLog(s string) string {
	// remove newlines for single-line logging
	s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}
