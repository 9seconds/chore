package config

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/9seconds/chore/internal/network"
)

const ParameterIP = "ip"

var (
	errInvalidIP           = errors.New("invalid IP address")
	errNotInAllowedSubnets = errors.New("cannot find address in allowed subnets")
)

type paramIP struct {
	required         bool
	resolve          bool
	allowedSubnets   []*net.IPNet
	forbiddenSubnets []*net.IPNet
}

func (p paramIP) Required() bool {
	return p.required
}

func (p paramIP) Type() string {
	return ParameterIP
}

func (p paramIP) String() string {
	return fmt.Sprintf(
		"required=%t, resolve=%t, allowedSubnets=%v, forbiddenSubnets=%v",
		p.required,
		p.resolve,
		p.allowedSubnets,
		p.forbiddenSubnets)
}

func (p paramIP) Validate(ctx context.Context, value string) error {
	addr := net.ParseIP(value)
	if addr == nil {
		return errInvalidIP
	}

	for _, subnet := range p.forbiddenSubnets {
		if subnet.Contains(addr) {
			return fmt.Errorf("address blacklisted in %s", subnet)
		}
	}

	found := true

	for _, subnet := range p.allowedSubnets {
		found = false

		if subnet.Contains(addr) {
			found = true
			break
		}
	}

	if !found {
		return errNotInAllowedSubnets
	}

	if p.resolve {
		if _, err := network.DNSResolver.LookupAddr(ctx, addr.String()); err != nil {
			return fmt.Errorf("cannot do reverse lookup: %w", err)
		}
	}

	return nil
}

func NewIP(required bool, spec map[string]string) (Parameter, error) {
	param := paramIP{
		required: required,
	}

	for _, v := range parseCSV(spec["allowed_subnets"]) {
		_, subnet, err := net.ParseCIDR(v)

		if err != nil {
			return nil, fmt.Errorf("%s is incorrect subnet: %w", v, err)
		}

		param.allowedSubnets = append(param.allowedSubnets, subnet)
	}

	for _, v := range parseCSV(spec["forbidden_subnets"]) {
		_, subnet, err := net.ParseCIDR(v)

		if err != nil {
			return nil, fmt.Errorf("%s is incorrect subnet: %w", v, err)
		}

		param.forbiddenSubnets = append(param.forbiddenSubnets, subnet)
	}

	if resolve, err := parseBool(spec, "resolve"); err == nil {
		param.resolve = resolve
	} else {
		return nil, err
	}

	return param, nil
}
