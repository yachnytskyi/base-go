package middleware

/*
 * Inspired by Golang's TimeoutHandler: https://golang.org/src/net/http/server.go?s=101514:101582#L3212
 * and gin-timeout: https://github.com/vearne/gin-timeout
 */

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/base-go/account/model/apperrors"
)

// Timeout wraps the request context with a timeout.
func Timeout(timeout time.Duration, errorTimeout *apperrors.Error) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set Gin's writer as our custom writer.
		timeWriter := &timeoutimeWriterriter{ResponseWriter: c.Writer, h: make(http.Header)}
		c.Writer = timeWriter

		// Wrap the request context with a timeout.
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Update gin request context.
		c.Request = c.Request.WithContext(ctx)

		finished := make(chan struct{})        // To indicate handler finished.
		panicChan := make(chan interface{}, 1) // Used to handle panics if we cannot recover.

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			c.Next() // Calls subsequent middleware(s) and handler.
			finished <- struct{}{}
		}()

		select {
		case <-panicChan:
			// If we cannot recover from panic,
			// send internal server error.
			e := apperrors.NewInternal()
			timeWriter.ResponseWriter.WriteHeader(e.Status())
			errorResponse, _ := json.Marshal(gin.H{
				"error": e,
			})
			timeWriter.ResponseWriter.Write(errorResponse)
		case <-finished:
			// If finished, set headers and write resp.
			timeWriter.mu.Lock()
			defer timeWriter.mu.Unlock()
			// Map Headers from timeWriter.Header() (written to by gin)
			// to timeWriter.ResponseWriter for response.
			dst := timeWriter.ResponseWriter.Header()
			for k, vv := range timeWriter.Header() {
				dst[k] = vv
			}
			timeWriter.ResponseWriter.WriteHeader(timeWriter.code)
			// timeWriter.writeBuffer will have been written to already when gin writes to timeWriter.Write().
			timeWriter.ResponseWriter.Write(timeWriter.writeBuffer.Bytes())
		case <-ctx.Done():
			// Timeout has occurred, send errorTimeout and write headers.
			timeWriter.mu.Lock()
			defer timeWriter.mu.Unlock()
			// ResponseWriter from gin.
			timeWriter.ResponseWriter.Header().Set("Content-Type", "application/json")
			timeWriter.ResponseWriter.WriteHeader(errorTimeout.Status())
			errorResponse, _ := json.Marshal(gin.H{
				"error": errorTimeout,
			})
			timeWriter.ResponseWriter.Write(errorResponse)
			c.Abort()
			timeWriter.SetTimedOut()
		}
	}
}

// Implements http.Writer, but tracks if Writer has timed out
// or has already written its header to prevent
// header and body overwrites.
// Also locks access to this writer to prevent race conditions
// holds the gin.ResponseWriter which we will manually call Write()
// on in the middleware function to send response.
type timeoutimeWriterriter struct {
	gin.ResponseWriter
	h           http.Header
	writeBuffer bytes.Buffer // The zero value for Buffer is an empty buffer ready to use.

	mu          sync.Mutex
	timedOut    bool
	wroteHeader bool
	code        int
}

// Writes the response, but first makes sure there
// has not already been a timeout
// In http.ResponseWriter interface.
func (timeWriter *timeoutimeWriterriter) Write(b []byte) (int, error) {
	timeWriter.mu.Lock()
	defer timeWriter.mu.Unlock()
	if timeWriter.timedOut {
		return 0, nil
	}

	return timeWriter.writeBuffer.Write(b)
}

// In http.ResponseWriter interface.
func (timeWriter *timeoutimeWriterriter) WriteHeader(code int) {
	checkWriteHeaderCode(code)
	timeWriter.mu.Lock()
	defer timeWriter.mu.Unlock()
	// We do not write the header if we have timed out or written the header
	if timeWriter.timedOut || timeWriter.wroteHeader {
		return
	}
	timeWriter.writeHeader(code)
}

// Set that the header has been written.
func (timeWriter *timeoutimeWriterriter) writeHeader(code int) {
	timeWriter.wroteHeader = true
	timeWriter.code = code
}

// Header "relays" the header, h, set in struct
// In http.ResponseWriter interface.
func (timeWriter *timeoutimeWriterriter) Header() http.Header {
	return timeWriter.h
}

// SetTimeOut sets timedOut field to true.
func (timeWriter *timeoutimeWriterriter) SetTimedOut() {
	timeWriter.timedOut = true
}

func checkWriteHeaderCode(code int) {
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid WriteHeader code %v", code))
	}
}
