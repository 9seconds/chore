package network

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"time"
)

const (
	connectTimeout = 2 * time.Second
	httpTimeout    = 30 * time.Second

	// practically, I'm pretty much sure we won't need to have a deal with
	// JSONs bigger than 1mb
	maxJSONSize = 1 * 1024 * 1024
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

func CloseResponse(resp *http.Response) error {
	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		return err
	}

	return resp.Body.Close()
}

func NewRequest(ctx context.Context, url string) *http.Request {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "chore")

	return req
}

func SendRequest(client *http.Client, req *http.Request) (*http.Response, error) {
	log.Printf("request %s", req.URL.String())

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot complete a request: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		CloseResponse(resp) //nolint: errcheck

		return nil, fmt.Errorf("unexpected response status code %d", resp.StatusCode)
	}

	log.Printf("response %s: %s (length: %d)", req.URL.String(), resp.Status, resp.ContentLength)

	return resp, err
}

func DoJSONRequestWithClient(
	ctx context.Context,
	client *http.Client,
	url string,
	target interface{},
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	req := NewRequest(ctx, url)
	req.Header.Set("Accept", "application/json")

	resp, err := SendRequest(client, req)
	if err != nil {
		return err
	}

	defer CloseResponse(resp) //nolint: errcheck

	decoder := json.NewDecoder(io.LimitReader(resp.Body, maxJSONSize))

	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("cannot decode JSON: %w", err)
	}

	return nil
}

func DoJSONRequest(ctx context.Context, url string, target interface{}) error {
	return DoJSONRequestWithClient(ctx, HTTPClient, url, target)
}
