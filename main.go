package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/twanies/jinx/proxy"
)

const defaultAddr = ":9999"

func main() {
	var (
		conf = flag.String("conf", "conf/jinx.conf", "config filename")
	)
	flag.Parse()

	cfg := proxy.LoadConfig(*conf)
	p := proxy.NewReverseProxy()

	for address, settings := range cfg.Servers {
		if isHTTPS(address) {
			cfg.TLS = true
		}
		if settings.Scheme == "" {
			settings.Scheme = "http"
		}
		target := &url.URL{
			Scheme: settings.Scheme,
			Host:   parseHost(settings.Proxy),
			Path:   settings.Path,
		}
		up := proxy.NewUpstream(parseHost(address), target)
		p.AddUpstream(up)
	}

	listen := defaultAddr
	if cfg.Listen > 0 {
		listen = fmt.Sprintf(":%d", cfg.Listen)
	}
	if cfg.TLS {
		log.Fatal(http.ListenAndServeTLS(listen, cfg.CertFile, cfg.KeyFile, p))
	}
	log.Fatal(http.ListenAndServe(listen, p))
}

func isHTTPS(address string) bool {
	return strings.Contains(address, "https")
}

func parseHost(address string) string {
	if strings.Contains(address, "http://") {
		return address[7:]
	} else if strings.Contains(address, "https://") {
		return address[8:]
	}
	return address
}
