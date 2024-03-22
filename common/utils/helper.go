package utils

import (
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
)

// CompareHashAndPassword
//
//	@Description: 密码比较
//	@param e
//	@param p
//	@return bool
//	@return error
func CompareHashAndPassword(e string, p string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(e), []byte(p))
	if err != nil {
		return false, err
	}
	return true, nil
}

// Encrypt
//
//	@Description: 加密
//	@param password
//	@return err
func Encrypt(password string) (err error, pwd string) {
	if password == "" {
		return
	}
	var hash []byte
	if hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err != nil {
		return
	} else {
		pwd = string(hash)
		return
	}
}

// CreateTracer
//
//	@Description:
//	@param serviceName
//	@param header
//	@return opentracing.Tracer
//	@return opentracing.SpanContext
//	@return io.Closer
//	@return error
func CreateTracer(serviceName string, header http.Header) (opentracing.Tracer, opentracing.SpanContext, io.Closer, error) {
	var cfg = jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
			// 按实际情况替换你的 ip
			CollectorEndpoint: "http://192.168.111.143:6831",
		},
	}
	jLogger := jaegerlog.StdLogger
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
	)
	spanContext, _ := tracer.Extract(opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(header))
	return tracer, spanContext, closer, err
}

func CreateRequestId() string {
	requestId := uuid.New().String()
	return requestId
}
