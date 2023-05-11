package gin_logrus

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var timeFormat = "2006-01-02 15.04.05 -07"

// FieldFunc additional request info, such as seqId, userId.
type FieldFunc func(ctx *gin.Context) logrus.Fields

func StringFields(key string) FieldFunc {
	return func(ctx *gin.Context) logrus.Fields {
		if key != "" {
			return logrus.Fields{
				key: ctx.GetString(key),
			}
		}
		return nil
	}
}

type Lugrus struct {
	logger        *logrus.Logger
	skipPath      []string
	fieldHandlers []FieldFunc
}

func (l *Lugrus) HandleFunc() gin.HandlerFunc {
	skip := make(map[string]interface{})
	for _, p := range l.skipPath {
		skip[p] = struct{}{}
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if _, ok := skip[path]; ok {
			c.Next()
			return
		}
		start := time.Now()
		c.Next()
		end := time.Now()
		latency := end.Sub(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		method := c.Request.Method
		referer := c.Request.Referer()
		fields := logrus.Fields{
			"statusCode": statusCode,
			"latency":    latency,
			"clientIP":   clientIP,
			"method":     method,
			"path":       path,
			"referer":    referer,
			"userAgent":  clientUserAgent,
		}

		entry := l.logger.WithFields(fields)
		for _, handler := range l.fieldHandlers {
			entry = entry.WithFields(handler(c))
		}

		msg := fmt.Sprintf("%d %s %s %4s %s",
			statusCode,
			time.Now().Format(timeFormat),
			latency.String(),
			method,
			path)

		if statusCode >= http.StatusInternalServerError {
			entry.Error(msg)
		} else if statusCode >= http.StatusBadRequest {
			entry.Warn(msg)
		} else {
			entry.Info(msg)
		}
	}
}

func Logger(logger *logrus.Logger) *Lugrus {
	return &Lugrus{
		logger:        logger,
		skipPath:      []string{},
		fieldHandlers: []FieldFunc{},
	}
}

func (l *Lugrus) WithNoLog(path ...string) *Lugrus {
	l.skipPath = append(l.skipPath, path...)
	return l
}

func (l *Lugrus) WithFields(fields ...FieldFunc) *Lugrus {
	l.fieldHandlers = append(l.fieldHandlers, fields...)
	return l
}
