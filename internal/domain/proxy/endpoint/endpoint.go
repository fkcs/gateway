// 根据路由拆分，内部参数拆解，用于数据过滤
package endpoint

import (
	"github.com/fkcs/gateway/internal/infrastructure/logger"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type EndPointer interface {
	ForwardProxy(rw http.ResponseWriter, req *http.Request)
	ReverseProxy(rw http.ResponseWriter, req *http.Request, destNode string)
}

type EndPoint struct {
}

func NewEndPoint() *EndPoint {
	return &EndPoint{}
}

// 正向代理
func (x *EndPoint) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	logger.Logger().Infof("Received request %s %s %s", req.Method, req.Host, req.RemoteAddr)
	outReq := new(http.Request)
	*outReq = *req
	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		outReq.Header.Set("X-Forwarded-For", clientIP)
	}
	res, err := http.DefaultTransport.RoundTrip(outReq)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}
	for key, value := range res.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}
	rw.WriteHeader(res.StatusCode)
	io.Copy(rw, res.Body)
	res.Body.Close()
}

func (x *EndPoint) ForwardProxy(rw http.ResponseWriter, req *http.Request) {
	x.ServeHTTP(rw, req)
}

// 反向代理
func (x *EndPoint) hostReverseProxy(w http.ResponseWriter, req *http.Request, host string) {
	remote, _ := url.Parse(host)
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, req)
}

func (x *EndPoint) ReverseProxy(rw http.ResponseWriter, req *http.Request, destNode string) {
	x.hostReverseProxy(rw, req, destNode)
}
