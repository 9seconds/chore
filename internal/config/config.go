package config

import (
	"fmt"
	"io"

	"github.com/9seconds/chore/internal/git"
)

type Config struct {
	Description string
	Git         git.AccessMode
	Network     bool
	Parameters  map[string]Parameter
	Flags       map[string]Flag
}

func Parse(reader io.Reader) (Config, error) { //nolint: cyclop
	raw, err := parseRaw(reader)
	if err != nil {
		return Config{}, err
	}

	gitMode, err := git.GetAccessMode(raw.Git)
	if err != nil {
		return Config{}, fmt.Errorf("cannot parse git access mode: %w", err)
	}

	conf := Config{
		Description: raw.Description,
		Network:     raw.Network,
		Git:         gitMode,
		Parameters:  make(map[string]Parameter),
		Flags:       make(map[string]Flag),
	}

	for name, param := range raw.Flags {
		conf.Flags[name] = NewFlag(param.Description, param.Required)
	}

	for name, param := range raw.Parameters {
		name := NormalizeName(name)

		var (
			value Parameter
			err   error
		)

		switch param.Type {
		case ParameterInteger:
			value, err = NewInteger(param.Description, param.Required, param.Spec)
		case ParameterString:
			value, err = NewString(param.Description, param.Required, param.Spec)
		case ParameterFloat:
			value, err = NewFloat(param.Description, param.Required, param.Spec)
		case ParameterURL:
			value, err = NewURL(param.Description, param.Required, param.Spec)
		case ParameterEmail:
			value, err = NewEmail(param.Description, param.Required, param.Spec)
		case ParameterEnum:
			value, err = NewEnum(param.Description, param.Required, param.Spec)
		case ParameterBase64:
			value, err = NewBase64(param.Description, param.Required, param.Spec)
		case ParameterHex:
			value, err = NewHex(param.Description, param.Required, param.Spec)
		case ParameterHostname:
			value, err = NewHostname(param.Description, param.Required, param.Spec)
		case ParameterMac:
			value, err = NewMac(param.Description, param.Required, param.Spec)
		case ParameterJSON:
			value, err = NewJSON(param.Description, param.Required, param.Spec)
		case ParameterXML:
			value, err = NewXML(param.Description, param.Required, param.Spec)
		case ParameterUUID:
			value, err = NewUUID(param.Description, param.Required, param.Spec)
		case ParameterDirectory:
			value, err = NewDirectory(param.Description, param.Required, param.Spec)
		case ParameterFile:
			value, err = NewFile(param.Description, param.Required, param.Spec)
		case ParameterSemver:
			value, err = NewSemver(param.Description, param.Required, param.Spec)
		case ParameterDatetime:
			value, err = NewDatetime(param.Description, param.Required, param.Spec)
		case ParameterGit:
			value, err = NewGit(param.Description, param.Required, param.Spec, git.Get)
		default:
			return conf, fmt.Errorf("unknown parameter type %s for parameter %s", param.Type, name)
		}

		if err != nil {
			return conf, fmt.Errorf("cannot initialize parameter %s: %w", name, err)
		}

		conf.Parameters[name] = value
	}

	return conf, nil
}
