package common

import (
	"github.com/fkcs/gateway/internal/infrastructure/logger"
	"github.com/valyala/fasthttp"
	"net/http"
)

type Method struct {
	url         string
	contentType string
}

func NewMethod() *Method {
	return &Method{
		url:         "127.0.0.1:8080",
		contentType: "application/x-www-form-urlencoded",
	}
}

func (m *Method) SetContentType(contentType string) *Method {
	m.contentType = contentType
	return m
}

func (m *Method) SetHost(host string) *Method {
	m.url = host
	return m
}

func (m *Method) Do(method string, path string, body []byte, args map[string]string, header map[string]string) (int, []byte) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	if len(m.contentType) > 0 {
		req.Header.SetContentType(m.contentType)
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	req.Header.SetMethod(method)
	uri := fasthttp.AcquireURI()
	defer fasthttp.ReleaseURI(uri)
	uri.SetScheme("http")
	uri.SetHost(m.url)
	uri.SetPath(path)
	if len(args) > 0 {
		a := fasthttp.AcquireArgs()
		defer fasthttp.ReleaseArgs(a)
		for key, value := range args {
			a.Add(key, value)
		}
		uri.SetQueryStringBytes(a.QueryString())
	}
	req.SetRequestURIBytes(uri.FullURI())
	if len(body) > 0 {
		req.SetBody(body)
	}
	logger.Logger().Debugf("[Request] %s", req.String())
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	if err := fasthttp.Do(req, resp); err != nil {
		logger.Logger().Errorf("err:%v\n[Response]\n%s", err.Error(), resp.Body())
		return http.StatusInternalServerError, nil
	}
	logger.Logger().Debugf("[Response] %v", resp.String())
	return resp.StatusCode(), resp.Body()
}
