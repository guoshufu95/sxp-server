package helper

import (
	"context"
	"errors"
	"google.golang.org/grpc/metadata"
	"time"
)

const defaultDeadline = 5

// BuildTokenCtx
//
//	@Description: ctx中注入token
//	@param token
//	@return ctx
func BuildTokenCtx(token string) context.Context {
	// 创建metadata和context.
	md := metadata.Pairs("token", token)
	ctx, _ := context.WithTimeout(context.Background(), defaultDeadline*time.Second)
	c := metadata.NewOutgoingContext(ctx, md)
	return c
}

// CheckTokenRes
//
//	@Description: grpc服务端返回的token校验
//	@param header
//	@return err
//	@return flag
func CheckTokenRes(header metadata.MD) (err error, flag bool) {
	// 从返回响应的header中读取数据.
	h, ok := header["check_token"]
	if !ok {
		err = errors.New("获取token校验失败")
		return
	}
	if h[0] != "1" { // string 1：正常 0：失败
		err = errors.New("grpc服务端token校验失败")
		return
	}
	flag = true
	return
}

// CheckSign
//
//	@Description: grpc服务端返回标志位校验
//	@param trailer
//	@return err
//	@return falg
func CheckSign(trailer metadata.MD) (err error, falg bool) {
	// 从返回响应的trailer中读取sign.
	tl, ok := trailer["sign"]
	if !ok {
		err = errors.New("获取返回标志失败")
		return
	}
	if tl[0] != "sxp-alan" { // 可根据自己的业务逻辑进行校验
		err = errors.New("sign校验失败")
		return
	}
	falg = true
	return
}
