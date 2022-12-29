package config

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"regexp"

	"github.com/9seconds/chore/internal/network"
)

const ParameterURL = "url"

type paramURL struct {
	baseParameter

	resolve  bool
	scheme   string
	domainRE *regexp.Regexp
	pathRE   *regexp.Regexp
	userRE   *regexp.Regexp
}

func (p paramURL) Type() string {
	return ParameterURL
}

func (p paramURL) String() string {
	return fmt.Sprintf(
		"%q (required=%t, scheme=%s, domain_re=%v, path_re=%v, user_re=%v, resolve=%t)",
		p.description,
		p.required,
		p.scheme,
		p.domainRE,
		p.pathRE,
		p.userRE,
		p.resolve)
}

func (p paramURL) Validate(ctx context.Context, value string) error { //nolint: cyclop
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	parsed, err := url.Parse(value)

	switch {
	case err != nil:
		return fmt.Errorf("incorrect url: %w", err)
	case p.scheme != "" && p.scheme != parsed.Scheme:
		return fmt.Errorf("incorrect scheme %s", parsed.Scheme)
	case p.domainRE != nil && !p.domainRE.MatchString(parsed.Host):
		return fmt.Errorf("incorrect host %s", parsed.Host)
	case p.pathRE != nil && !p.pathRE.MatchString(parsed.Path):
		return fmt.Errorf("incorrect path %s", parsed.Path)
	case p.userRE != nil && !p.userRE.MatchString(parsed.User.Username()):
		return fmt.Errorf("incorrect user '%s'", parsed.User.Username())
	case p.resolve:
		address := parsed.Host

		if _, _, err := net.SplitHostPort(address); err != nil {
			switch parsed.Scheme {
			case "http":
				address = net.JoinHostPort(address, "80")
			case "https":
				address = net.JoinHostPort(address, "443")
			default:
				return fmt.Errorf("unknown scheme %s", parsed.Scheme)
			}
		}

		conn, err := network.NetDialer.DialContext(ctx, "tcp", address)
		if err != nil {
			return fmt.Errorf("cannot dial to %s: %w", parsed.Host, err)
		}

		conn.Close()
	}

	return nil
}

func NewURL(description string, required bool, spec map[string]string) (Parameter, error) {
	param := paramURL{
		baseParameter: baseParameter{
			required:    required,
			description: description,
		},
		scheme: spec["scheme"],
	}

	if parsed, err := parseBool(spec, "resolve"); err == nil {
		param.resolve = parsed
	} else {
		return nil, err
	}

	if parsed, err := parseRegexp(spec, "domain_re"); err == nil {
		param.domainRE = parsed
	} else {
		return nil, err
	}

	if parsed, err := parseRegexp(spec, "path_re"); err == nil {
		param.pathRE = parsed
	} else {
		return nil, err
	}

	if parsed, err := parseRegexp(spec, "user_re"); err == nil {
		param.userRE = parsed
	} else {
		return nil, err
	}

	return param, nil
}
