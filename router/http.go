package router

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"saga/utils"
	"saga/workflow"
)

const (
	API_PATH = "/api/rtu"
	API_PORT = 8998
)

func RunHTTP(app *gin.Engine) {
	log.Printf("Starting SAGA HTTP GIN Server at: %d", API_PORT)
	err := app.Run(fmt.Sprintf(":%d", API_PORT))
	utils.FatalIfError(err)
}

func SetUp() *gin.Engine {
	app := GetGinApp()

	AddRoute(app)

	return app
}

func GetGinApp() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	app.Use(gin.Recovery())

	return app
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 요청 본문을 읽습니다.
		bodyBytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			utils.PrintErrorf("REQUEST READ FAILED: %v", err)
		} else {
			// 본문을 복원하여 이후 핸들러에서도 읽을 수 있도록 합니다.
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			utils.PrintInfof("<- REQUEST: [%s] %s | Body: %s", c.Request.Method, c.Request.RequestURI, string(bodyBytes))
		}
		// 다음 핸들러로 제어를 넘깁니다.
		c.Next()
	}
}

func AddRoute(app *gin.Engine) {
	app.Use(LoggingMiddleware())
	app.POST(API_PATH+"/submit", utils.WrapHandler(func(c *gin.Context) interface{} {
		data, err := ioutil.ReadAll(c.Request.Body)
		utils.FatalIfError(err)
		return workflow.Execute(c.Request.URL.Query(), data)
	}))
}
