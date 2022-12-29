package config

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	"github.com/9seconds/chore/internal/network"
)

const ParameterEmail = "email"

var (
	errMalformedEmail = errors.New("incorrect email")
	errNoMXRecords    = errors.New("no mx records were found")

	emailRegexp = regexp.MustCompile(
		// https://html.spec.whatwg.org/multipage/input.html#valid-e-mail-address
		"^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
	)
)

type paramEmail struct {
	baseParameter

	resolve  bool
	domainRE *regexp.Regexp
	nameRE   *regexp.Regexp
}

func (p paramEmail) Type() string {
	return ParameterEmail
}

func (p paramEmail) String() string {
	return fmt.Sprintf(
		"%q (required=%t, resolve=%t, domain_re=%v, name_re=%v)",
		p.description,
		p.required,
		p.resolve,
		p.domainRE,
		p.nameRE)
}

func (p paramEmail) Validate(ctx context.Context, value string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	name, domain, found := strings.Cut(value, "@")
	parsed, err := mail.ParseAddress(value)

	switch {
	case err != nil, !found, parsed.Address != value, !emailRegexp.MatchString(value):
		return errMalformedEmail
	case p.domainRE != nil && !p.domainRE.MatchString(domain):
		return fmt.Errorf("domain %s does not match %s", domain, p.domainRE)
	case p.nameRE != nil && !p.nameRE.MatchString(name):
		return fmt.Errorf("name %s does not match %s", name, p.nameRE)
	case p.resolve:
		values, err := network.DNSResolver.LookupMX(ctx, domain)

		switch {
		case err != nil:
			return fmt.Errorf("cannot resolve MX records of the domain: %w", err)
		case len(values) == 0:
			return errNoMXRecords
		}
	}

	return nil
}

func NewEmail(description string, required bool, spec map[string]string) (Parameter, error) {
	param := paramEmail{
		baseParameter: baseParameter{
			required:    required,
			description: description,
		},
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

	if parsed, err := parseRegexp(spec, "name_re"); err == nil {
		param.nameRE = parsed
	} else {
		return nil, err
	}

	return param, nil
}
