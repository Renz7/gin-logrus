package gin_logrus

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestLugrus_Handler(t *testing.T) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	bf := bytes.Buffer{}
	logger.SetOutput(&bf)
	f := func(ctx *gin.Context) logrus.Fields {
		return logrus.Fields{
			"uid": ctx.GetString("uid"),
		}
	}
	handleFunc := Logger(logger).WithFields(f).HandleFunc()

	app := gin.New()
	w := httptest.NewRecorder()
	c := gin.CreateTestContextOnly(w, app)
	c.Set("uid", "uid value")
	c.Request = httptest.NewRequest("GET", "/", nil)
	handleFunc(c)

	var e = map[string]string{}
	json.Unmarshal(bf.Bytes(), &e)
	_, ok := e["uid"]
	assert.True(t, ok)
}
