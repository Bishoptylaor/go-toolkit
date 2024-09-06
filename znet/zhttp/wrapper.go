package zhttp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type HttpClientWrapper interface {
	// Call 普通调用方法，返回标准 http.Response
	Call(ctx context.Context, method string, url string, payload map[string]any, headers map[string]string) (*http.Response, []byte, error)
	// CallWrite 与 Call 同参数，会提前解析结果到 resStruct 结构体中
	CallWrite(ctx context.Context, method string, url string, payload map[string]any, headers map[string]string, resStruct any) (*http.Response, error)
	// CallOp 使用 opts 的方式注入相关内容调用
	CallOp(ctx context.Context, payload map[string]any, opts ...CfgFunc) (*http.Response, []byte, error)
	CallOpOk(ctx context.Context, payload map[string]any, opts ...CfgFunc) (*http.Response, []byte, error)
	// CallOpWrite 与 CallOp 同参数 会提前解析结果到 resStruct 结构体中
	CallOpWrite(ctx context.Context, payload map[string]any, resStruct any, opts ...CfgFunc) (*http.Response, error)
	// CallOpWriteOk 与 CallOpWrite 同参数 会提前解析结果到 resStruct 结构体中；会提前根据 statuscode != 200 跳出
	CallOpWriteOk(ctx context.Context, payload map[string]any, resStruct any, opts ...CfgFunc) (*http.Response, error)
}

type HcWrapper struct {
	client *http.Client
	cfg    *httpConfig
}

func (h *HcWrapper) Call(ctx context.Context, method string, url string, payload map[string]any, headers map[string]string) (*http.Response, []byte, error) {
	if _, ok := _MethodMap[method]; !ok {
		return nil, nil, errors.New("not supported http method")
	}
	return h.CallOp(ctx, payload, _MethodMap[method](url), Header(headers))
}

func (h *HcWrapper) CallWrite(ctx context.Context, method string, url string, payload map[string]any, headers map[string]string, resStruct any) (*http.Response, error) {
	if _, ok := _MethodMap[method]; !ok {
		return nil, errors.New("not supported http method")
	}
	return h.CallOpWrite(ctx, payload, resStruct, _MethodMap[method](url), Header(headers))
}

func (h *HcWrapper) CallOp(ctx context.Context, payload map[string]any, opts ...CfgFunc) (*http.Response, []byte, error) {
	for _, opt := range opts {
		opt(h.cfg)
	}
	if h.cfg.method == "" || h.cfg.url == "" {
		return nil, nil, fmt.Errorf("method is empty")
	}
	// 自主设置超时时间
	var cancel context.CancelFunc
	if h.cfg.timeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, h.cfg.timeout)
		defer cancel()
	}

	// 组装request
	req, err := h.cfg.MakeReq(ctx, payload)
	if err != nil {
		return nil, nil, err
	}

	for _, p := range h.cfg.Prefix {
		p(ctx, req)
	}

	// 组装 header
	req.Header = h.cfg.headers
	// 组装 cookie
	if h.cfg.cookies != nil {
		req.AddCookie(h.cfg.cookies)
	}

	// 发送请求
	res, err := h.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("request do err, param: %v, %w", req, err)
	}
	defer res.Body.Close()

	for _, s := range h.cfg.Suffix {
		s(ctx, res)
	}

	bs, err := io.ReadAll(io.LimitReader(res.Body, int64(h.cfg.bodySize<<20))) // default 10MB change the size you want
	if err != nil {
		return nil, nil, err
	}
	return res, bs, nil
}

func (h *HcWrapper) CallOpOk(ctx context.Context, payload map[string]any, opts ...CfgFunc) (*http.Response, []byte, error) {
	res, bs, err := h.CallOp(ctx, payload, opts...)
	if err != nil {
		return nil, nil, err
	}
	if res.StatusCode != http.StatusOK {
		return res, []byte{}, fmt.Errorf("StatusCode(%d) != 200", res.StatusCode)
	}
	return res, bs, nil
}

func (h *HcWrapper) CallOpWrite(ctx context.Context, payload map[string]any, resStruct any, opts ...CfgFunc) (*http.Response, error) {
	res, bs, err := h.CallOp(ctx, payload, opts...)
	if err != nil {
		return nil, err
	}
	return h.cfg.MakeRes(ctx, res, bs, resStruct)
}

func (h *HcWrapper) CallOpWriteOk(ctx context.Context, payload map[string]any, resStruct any, opts ...CfgFunc) (*http.Response, error) {
	res, bs, err := h.CallOp(ctx, payload, opts...)
	if err != nil {
		return nil, err
	}
	return h.cfg.MakeResOk(ctx, res, bs, resStruct)
}

func NewHttpClientWrapper() HttpClientWrapper {
	return GetDefaultClient()
}