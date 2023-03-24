package env

import (
	"context"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/9seconds/chore/internal/network"
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

var ipInfoOrgFormat = regexp.MustCompile(`^AS(\d+)\s+(.*?)$`)

func GenerateNetwork(
	ctx context.Context,
	results chan<- string,
	waiters *sync.WaitGroup,
	requireNetwork bool,
) {
	if !requireNetwork {
		return
	}

	waiters.Add(1)

	go func() {
		defer waiters.Done()

		if _, ok := os.LookupEnv(NetworkIPv4); ok {
			return
		}

		resp := ipInfoResponse{}

		err := network.DoJSONRequestWithClient(
			ctx,
			network.HTTPClientV4,
			"https://ipinfo.io/json",
			&resp)
		if err != nil {
			log.Printf("cannot request network data: %v", err)

			return
		}

		sendValue(ctx, results, NetworkIPv4, resp.IP)
		sendValue(ctx, results, NetworkHostname, resp.Hostname)
		sendValue(ctx, results, NetworkCity, resp.City)
		sendValue(ctx, results, NetworkRegion, resp.Region)
		sendValue(ctx, results, NetworkCountry, resp.Country)
		sendValue(ctx, results, NetworkPostal, resp.Postal)
		sendValue(ctx, results, NetworkTimezone, resp.Timezone)

		asnChunks := ipInfoOrgFormat.FindStringSubmatch(resp.Org)

		switch {
		case asnChunks == nil && resp.Org != "":
			sendValue(ctx, results, NetworkOrganization, resp.Org)
		case asnChunks != nil:
			sendValue(ctx, results, NetworkASN, asnChunks[1])
			sendValue(ctx, results, NetworkOrganization, asnChunks[2])
		}

		if lat, lon, ok := strings.Cut(resp.Loc, ","); ok {
			sendValue(ctx, results, NetworkLatitude, lat)
			sendValue(ctx, results, NetworkLongitude, lon)
		}
	}()
}

func GenerateNetworkIPv6(
	ctx context.Context,
	results chan<- string,
	waiters *sync.WaitGroup,
	requireNetwork bool,
) {
	if !requireNetwork {
		return
	}

	waiters.Add(1)

	go func() {
		defer waiters.Done()

		if _, ok := os.LookupEnv(NetworkIPv6); ok {
			return
		}

		resp := ifConfigResponse{}

		err := network.DoJSONRequestWithClient(
			ctx,
			network.HTTPClientV6,
			"https://ifconfig.co",
			&resp)
		if err != nil {
			log.Printf("cannot get IPv6 address: %v", err)

			return
		}

		sendValue(ctx, results, NetworkIPv6, resp.IP)
	}()
}
