package proxy

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

var CheckClient = &http.Client{Timeout: 30 * time.Second}

func NewProxy(upstream string) (http.Handler, error) {
	target, err := url.Parse(upstream)
	if err != nil {
		return nil, err
	}

	// Validate the scheme and hostname of the target URL
	if target.Scheme != "http" && target.Scheme != "https" {
		return nil, fmt.Errorf("unsupported scheme %q", target.Scheme)
	}
	if target.Hostname() == "" {
		return nil, errors.New("missing upstream hostname")
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			outURL := *target
			query := outURL.Query()
			for key, values := range r.In.URL.Query() {
				for _, value := range values {
					query.Add(key, value)
				}
			}
			outURL.RawQuery = query.Encode()
			r.Out.URL = &outURL
		},
	}

	return proxy, nil
}

func CheckUpstream(upstream string) error {
	request, err := http.NewRequest(http.MethodHead, upstream, nil)
	if err != nil {
		return err
	}

	response, err := CheckClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}
