package middleware

import (
	"bufio"
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"strings"
	"sxp-server/common/logger"
	"time"
)

// LoggerMiddleware
//
//	@Description: 日志中间件
//	@return gin.HandlerFunc
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.GetLogger()
		// 开始时间
		startTime := time.Now()
		// 处理请求
		bf := bytes.NewBuffer(nil)
		wt := bufio.NewWriter(bf)
		_, err := io.Copy(wt, c.Request.Body)
		if err != nil {
			log.Errorf("copy body error, %s", err.Error())
			err = nil
		}
		rb, _ := ioutil.ReadAll(bf)
		param := strings.ReplaceAll(strings.ReplaceAll(string(rb), "\r\n", ""), " ", "")
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rb))
		c.Next()
		endTime := time.Now()
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		res, _ := c.Get("response")
		logData := map[string]interface{}{
			"latencyTime":  latencyTime,
			"method":       reqMethod,
			"uri":          reqUri,
			"requestParam": param,
			"response":     res,
			"responseCode": statusCode,
		}
		log.Info(logData)
		c.Next()
	}
}
