package env

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	connectTimeout = 2 * time.Second
	httpTimeout    = 10 * time.Second

	userAgent = "chore"
)

type ipInfoResponse struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
}

var (
	netDialer = net.Dialer{
		Timeout:       connectTimeout,
		FallbackDelay: -1,
		KeepAlive:     -1,
	}
	httpClientV4 = http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
				return netDialer.DialContext(ctx, "tcp4", address)
			},
		},
		Timeout: httpTimeout,
	}
	httpClientV6 = http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
				return netDialer.DialContext(ctx, "tcp6", address)
			},
		},
		Timeout: httpTimeout,
	}
)

func generateNetworkFromIPInfo(ctx context.Context, result chan<- string) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://ipinfo.io/json", nil)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := httpClientV4.Do(req)
	if err != nil {
		log.Printf("cannot fetch ipinfo data: %v", err)

		return
	}

	defer resp.Body.Close()

	data := ipInfoResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("cannot load JSON response: %v", err)

		return
	}

	if data.IP != "" {
		sendEnvValue(ctx, result, EnvNetworkIPv4, data.IP)
	}

	if data.Hostname != "" {
		sendEnvValue(ctx, result, EnvNetworkHostname, data.Hostname)
	}

	if data.City != "" {
		sendEnvValue(ctx, result, EnvNetworkHostname, data.City)
	}

	if data.City != "" {
		sendEnvValue(ctx, result, EnvNetworkCity, data.City)
	}

	if data.Region != "" {
		sendEnvValue(ctx, result, EnvNetworkRegion, data.Region)
	}

	if data.Country != "" {
		sendEnvValue(ctx, result, EnvNetworkCountry, data.Country)
	}

	if data.Org != "" {
		sendEnvValue(ctx, result, EnvNetworkOrganization, data.Org)
	}

	if data.Postal != "" {
		sendEnvValue(ctx, result, EnvNetworkPostal, data.Postal)
	}

	if data.Timezone != "" {
		sendEnvValue(ctx, result, EnvNetworkTimezone, data.Timezone)
	}

	if lat, lon, ok := strings.Cut(data.Loc, ","); ok {
		sendEnvValue(ctx, result, EnvNetworkLatitude, lat)
		sendEnvValue(ctx, result, EnvNetworkLongitude, lon)
	}
}

func generateNetworkIPv6(ctx context.Context, result chan<- string) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://ifconfig.co/", nil)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/plain")

	resp, err := httpClientV6.Do(req)
	if err != nil {
		log.Printf("cannot fetch ipv6 address: %v", err)

		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("cannot read response body: %v", err)

		return
	}

	strBody := string(body)
	ipAddr := net.ParseIP(strings.TrimSpace(strBody))

	switch {
	case ipAddr == nil:
		log.Printf("incorrect ip address %s", strBody)
	case ipAddr.To16() == nil:
		log.Printf("incorrect ipv6 address %s", strBody)
	default:
		sendEnvValue(ctx, result, EnvNetworkIPv6, ipAddr.To16().String())
	}
}
