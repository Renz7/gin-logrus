# gin-logrus

[logrus](https://github.com/sirupsen/logrus) middleware for [Gin](https://github.com/gin-gonic/gin).

```go
package main

import (
	"github.com/gin-gonic/gin"
	ginlogrus "github.com/renz7/gin-logrus"
	"github.com/sirupsen/logrus"
)

func main() {
	g := gin.New()
	logger := logrus.New()
	reqIdFunc := ginlogrus.StringFields("reqId")
	middleware := ginlogrus.Logger(logger).WithNoLog("'/").WithFields(reqIdFunc).HandleFunc()
	g.Use(middleware)

	g.GET("/ping", func(context *gin.Context) {
		context.String(200, "pong")
	})
	g.Run(":8080")
}
```
