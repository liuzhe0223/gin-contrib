// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package logger

import (
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	loggerFeild = "logger"
)

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

func ErrorLogger() gin.HandlerFunc {
	return gin.ErrorLoggerT(gin.ErrorTypeAny)
}

func ErrorLoggerT(typ gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		// avoid writting if we already wrote into the response body
		if !c.Writer.Written() {
			errors := c.Errors.ByType(typ)
			if len(errors) > 0 {
				c.JSON(-1, errors)
			}
		}
	}
}

// Instance a Logger middleware with the specified writter buffer.
// Example: os.Stdout, a file opened in write mode, a socket...
func LoggerWithWriter(out io.Writer) gin.HandlerFunc {
	logger := log.New(out, "", log.Lshortfile)

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		// Set logger to context
		c.Set(loggerFeild, logger)

		path := c.Request.URL.Path
		reqId := c.Request.Header.Get(ReqIdHeaderFeild)
		if reqId == "" {
			reqId = NewReqId()
			c.Request.Header.Set(ReqIdHeaderFeild, reqId)
		}
		c.Writer.Header().Set(ReqIdHeaderFeild, reqId)

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		statusColor := colorForStatus(statusCode)
		methodColor := colorForMethod(method)
		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

		logger.Printf("[GIN] [%s] %v |%s %3d %s| %13v | %s |%s  %s %-7s %s\n%s",
			reqId,
			end.Format("2006/01/02 - 15:04:05"),
			statusColor, statusCode, reset,
			latency,
			clientIP,
			methodColor, reset, method,
			path,
			comment,
		)
	}
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}

func Default(c *gin.Context) *log.Logger {
	return c.MustGet(loggerFeild).(*log.Logger)
}
