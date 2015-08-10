package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Proxy struct {
	upstreams map[string]*url.URL
}

func NewReverseProxy() *Proxy {
	return &Proxy{map[string]*url.URL{}}
}

func (p *Proxy) AddUpstream(u *upstream) {
	p.upstreams[u.Host] = u.Target
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url, exist := p.upstreams[r.Host]
	if !exist {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	director := func(req *http.Request) {
		req.Host = url.Host
		req.URL.Scheme = url.Scheme
		req.URL.Host = url.Host
		req.URL.Path = singleJoiningSlash(url.Path, req.URL.Path)
	}
	reverseProxy := &httputil.ReverseProxy{Director: director}
	reverseProxy.ServeHTTP(w, r)
}

type upstream struct {
	Host   string
	Target *url.URL
}

func NewUpstream(host string, target *url.URL) *upstream {
	return &upstream{host, target}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
