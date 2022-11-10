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

type Dialer interface {
	DialContext(context.Context, string, string) (net.Conn, error)
}

type Resolver interface {
	LookupHost(context.Context, string) ([]string, error)
	LookupAddr(context.Context, string) ([]string, error)
	LookupMX(context.Context, string) ([]*net.MX, error)
}

var (
	NetDialer Dialer = &net.Dialer{
		Timeout: connectTimeout,
	}

	DNSResolver Resolver = &net.Resolver{
		Dial: NetDialer.DialContext,
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
			DialContext: NetDialer.DialContext,
		},
		Timeout: httpTimeout,
	}

	HTTPClientV4 = &http.Client{
		Jar: CookieJar,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: func(ctx context.Context, _, address string) (net.Conn, error) {
				return NetDialer.DialContext(ctx, "tcp4", address)
			},
		},
		Timeout: httpTimeout,
	}

	HTTPClientV6 = &http.Client{
		Jar: CookieJar,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: func(ctx context.Context, _, address string) (net.Conn, error) {
				return NetDialer.DialContext(ctx, "tcp6", address)
			},
		},
		Timeout: httpTimeout,
	}
)
