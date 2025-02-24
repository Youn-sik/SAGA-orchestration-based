package utils

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"saga/logger"
	"time"
)

var PrintInfof = logger.Infof
var PrintErrorf = logger.Errorf
var PrintDebugf = logger.Debugf
var FatalIfError = logger.FatalIfError

func WrapHandler(fn func(*gin.Context) interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		began := time.Now()
		ret := fn(c)
		status, res := Result2HttpJSON(ret)

		b, _ := json.Marshal(res)
		if status == http.StatusOK || status == http.StatusTooEarly {
			PrintInfof("%2dms %d %s %s %s", time.Since(began).Milliseconds(), status, c.Request.Method, c.Request.RequestURI, string(b))
		} else {
			PrintInfof("%2dms %d %s %s %s", time.Since(began).Milliseconds(), status, c.Request.Method, c.Request.RequestURI, string(b))
		}
		c.JSON(status, res)
	}
}

func Result2HttpJSON(result interface{}) (code int, res interface{}) {
	err, _ := result.(error)
	if err == nil {
		code = http.StatusOK
		res = result
	} else {
		res = map[string]string{
			"error": err.Error(),
		}
		code = http.StatusInternalServerError
	}
	return
}

func GenerateTransactionID() int64 {
	var id int64

	// 1. 8바이트 난수 생성
	err := binary.Read(rand.Reader, binary.BigEndian, &id)
	if err != nil {
		panic(fmt.Sprintf("failed to generate random ID: %v", err)) // panic or handle error appropriately
	}

	// 2. 시간 정보 추가 (선택적) - 나노초 단위까지 고려
	nanoTime := time.Now().UnixNano()
	id ^= nanoTime // XOR 연산으로 섞어줌

	// 3. 음수 값 방지 (선택적)
	if id < 0 {
		id = -id
	}

	return id
}
