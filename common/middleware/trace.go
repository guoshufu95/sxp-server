package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"sxp-server/common/utils"
)

// UseOpenTracing
//
//	@Description: jaeger链路追踪
//	@return gin.HandlerFunc
func UseOpenTracing() gin.HandlerFunc {
	handler := func(c *gin.Context) {
		tracer, spanContext, closer, _ := utils.CreateTracer("sxp-server", c.Request.Header)
		defer closer.Close()
		startSpan := tracer.StartSpan(c.Request.URL.Path, ext.RPCServerOption(spanContext))
		defer startSpan.Finish()
		ext.HTTPUrl.Set(startSpan, c.Request.URL.Path)
		ext.HTTPMethod.Set(startSpan, c.Request.Method)
		ext.Component.Set(startSpan, "sxp-server")
		// 在 header 中加上当前进程的上下文信息
		c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), startSpan))
		// 传递给下一个中间件
		c.Next()
		// 继续设置 tag
		ext.HTTPStatusCode.Set(startSpan, uint16(c.Writer.Status()))
	}

	return handler
}
