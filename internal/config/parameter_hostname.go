package config

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/9seconds/chore/internal/network"
)

const ParameterHostname = "hosthame"

var (
	errIncorrectHostname = errors.New("incorrect hostname")
	errNoHostnameRecords = errors.New("no dns records")

	hostnameRegexp = regexp.MustCompile(
		// https://www.rfc-editor.org/rfc/rfc1123#page-13
		// https://github.com/go-playground/validator/blob/1e8c614c2a5449c8537d78a155c214d1dd50b030/regexes.go#L51
		`^([a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,62}){1}(\.[a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,62})*?$`,
	)
	fqdnRegexp = regexp.MustCompile(
		// https://www.rfc-editor.org/rfc/rfc1123#page-13
		// https://github.com/go-playground/validator/blob/1e8c614c2a5449c8537d78a155c214d1dd50b030/regexes.go#L52
		`^([a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,62})(\.[a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,62})*?(\.[a-zA-Z]{1}[a-zA-Z0-9]{0,62})\.?$`,
	)
)

type paramHostname struct {
	mixinStringLength

	required bool
	isFQDN   bool
	resolve  bool
	re       *regexp.Regexp
}

func (p paramHostname) Required() bool {
	return p.required
}

func (p paramHostname) Type() string {
	return ParameterHostname
}

func (p paramHostname) String() string {
	return fmt.Sprintf(
		"required=%t, is_fqdn=%t, resolve=%t, re=%v, %s",
		p.required,
		p.isFQDN,
		p.resolve,
		p.re,
		p.mixinStringLength)
}

func (p paramHostname) Validate(ctx context.Context, value string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := p.mixinStringLength.Validate(value); err != nil {
		return err
	}

	switch {
	case p.isFQDN && !fqdnRegexp.MatchString(value):
		return errIncorrectHostname
	case !p.isFQDN && !hostnameRegexp.MatchString(value):
		return errIncorrectHostname
	case p.re != nil && !p.re.MatchString(value):
		return fmt.Errorf("hostname does not match %s", p.re)
	case p.resolve:
		values, err := network.DNSResolver.LookupHost(ctx, value)

		switch {
		case err != nil:
			return fmt.Errorf("cannot resolve dns records: %w", err)
		case len(values) == 0:
			return errNoHostnameRecords
		}
	}

	return nil
}

func NewHostname(required bool, spec map[string]string) (Parameter, error) {
	param := paramHostname{
		required: required,
	}

	if mixin, err := makeMixinStringLength(spec, "min_length", "max_length"); err == nil {
		param.mixinStringLength = mixin
	} else {
		return nil, err
	}

	if value, err := parseBool(spec, "is_fqdn"); err == nil {
		param.isFQDN = value
	} else {
		return nil, err
	}

	if value, err := parseBool(spec, "resolve"); err == nil {
		param.resolve = value
	} else {
		return nil, err
	}

	if value, err := parseRegexp(spec, "regexp"); err == nil {
		param.re = value
	} else {
		return nil, err
	}

	return param, nil
}
