package env

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/9seconds/chore/internal/network"
)

const (
	connectTimeout = 2 * time.Second
	httpTimeout    = 10 * time.Second
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

type ifConfigResponse struct {
	IP string `json:"ip"`
}

var (
	ipInfoOrgFormat = regexp.MustCompile(`^AS(\d+)\s+(.*?)$`)
)

func doRequest(ctx context.Context, client *http.Client, url string, target interface{}) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "chore")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("cannot access endpoint: %w", err)
	}

	defer func() {
		io.Copy(io.Discard, resp.Body) //nolint: errcheck
		resp.Body.Close()
	}()

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unexpected response status code %d", resp.StatusCode)
	}

	reader := bufio.NewReader(resp.Body)
	decoder := json.NewDecoder(reader)

	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("cannot parse response: %w", err)
	}

	return nil
}

func GenerateNetwork(ctx context.Context, results chan<- string, waiters *sync.WaitGroup) {
	waiters.Add(1)

	go func() {
		defer waiters.Done()

		if _, ok := os.LookupEnv(EnvNetworkIPv4); ok {
			return
		}

		resp := ipInfoResponse{}
		if err := doRequest(ctx, network.HTTPClientV4, "https://ipinfo.io/json", &resp); err != nil {
			log.Printf("cannot request network data: %v", err)

			return
		}

		sendValue(ctx, results, EnvNetworkIPv4, resp.IP)
		sendValue(ctx, results, EnvNetworkHostname, resp.Hostname)
		sendValue(ctx, results, EnvNetworkCity, resp.City)
		sendValue(ctx, results, EnvNetworkRegion, resp.Region)
		sendValue(ctx, results, EnvNetworkCountry, resp.Country)
		sendValue(ctx, results, EnvNetworkPostal, resp.Postal)
		sendValue(ctx, results, EnvNetworkTimezone, resp.Timezone)

		asnChunks := ipInfoOrgFormat.FindStringSubmatch(resp.Org)

		switch {
		case asnChunks == nil && resp.Org != "":
			sendValue(ctx, results, EnvNetworkOrganization, resp.Org)
		case asnChunks != nil:
			sendValue(ctx, results, EnvNetworkASN, asnChunks[1])
			sendValue(ctx, results, EnvNetworkOrganization, asnChunks[2])
		}

		if lat, lon, ok := strings.Cut(resp.Loc, ","); ok {
			sendValue(ctx, results, EnvNetworkLatitude, lat)
			sendValue(ctx, results, EnvNetworkLongitude, lon)
		}
	}()
}

func GenerateNetworkIPv6(ctx context.Context, results chan<- string, waiters *sync.WaitGroup) {
	waiters.Add(1)

	go func() {
		defer waiters.Done()

		if _, ok := os.LookupEnv(EnvNetworkIPv6); ok {
			return
		}

		resp := ifConfigResponse{}
		if err := doRequest(ctx, network.HTTPClientV6, "https://ifconfig.co", &resp); err != nil {
			log.Printf("cannot get IPv6 address: %v", err)

			return
		}

		sendValue(ctx, results, EnvNetworkIPv6, resp.IP)
	}()
}
