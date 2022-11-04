package network

import (
	"context"
	"net"
	"net/http"
	"net/http/cookiejar"
	"time"
)

const (
	connectTimeout = 2 * time.Second
	httpTimeout    = 10 * time.Second
)

var (
	Dialer = &net.Dialer{
		Timeout: connectTimeout,
	}

	CookieJar = func() *cookiejar.Jar {
		jar, err := cookiejar.New(nil)
		if err != nil {
			panic(err)
		}

		return jar
	}()

	HTTPClient = &http.Client{
		Jar: CookieJar,
		Transport: &http.Transport{
			Proxy:       http.ProxyFromEnvironment,
			DialContext: Dialer.DialContext,
		},
		Timeout: httpTimeout,
	}

	HTTPClientV4 = &http.Client{
		Jar: CookieJar,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: func(ctx context.Context, _, address string) (net.Conn, error) {
				return Dialer.DialContext(ctx, "tcp4", address)
			},
		},
		Timeout: httpTimeout,
	}

	HTTPClientV6 = &http.Client{
		Jar: CookieJar,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: func(ctx context.Context, _, address string) (net.Conn, error) {
				return Dialer.DialContext(ctx, "tcp6", address)
			},
		},
		Timeout: httpTimeout,
	}
)
