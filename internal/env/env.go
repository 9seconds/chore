package env

import (
	"strings"
)

const (
	DebugEnabled = "1"
)

const (
	Prefix             = "CHORE_"
	EnvParameterPrefix = Prefix + "P_"
	EnvFlagPrefix      = Prefix + "F_"
	EnvPathPrefix      = Prefix + "PATH_"
	EnvNetworkPrefix   = Prefix + "NETWORK_"
	EnvStartedAtPrefix = Prefix + "STARTED_AT_"
	EnvIDPrefix        = Prefix + "ID_"
	EnvIDChainPrefix   = Prefix + "CHAIN_ID_"
	EnvOSPrefix        = Prefix + "OS_"
	EnvOSVersionPrefix = EnvOSPrefix + "VERSION_"
	EnvGitPrefix       = Prefix + "GIT_"
	EnvUserPrefix      = Prefix + "USER_"

	Debug = Prefix + "DEBUG"

	EnvNamespace = Prefix + "NAMESPACE"
	EnvCaller    = Prefix + "CALLER"
	EnvSelf      = Prefix + "SELF"
	EnvSlug      = Prefix + "SLUG"

	EnvPathCaller = EnvPathPrefix + "CALLER"
	EnvPathData   = EnvPathPrefix + "DATA"
	EnvPathCache  = EnvPathPrefix + "CACHE"
	EnvPathState  = EnvPathPrefix + "STATE"
	EnvPathTemp   = EnvPathPrefix + "TEMP"

	EnvIDRun           = EnvIDPrefix + "RUN"
	EnvIDChainRun      = EnvIDChainPrefix + "RUN"
	EnvIDIsolated      = EnvIDPrefix + "ISOLATED"
	EnvIDChainIsolated = EnvIDChainPrefix + "ISOLATED"

	EnvMachineID = Prefix + "MACHINE_ID"

	EnvOSType         = EnvOSPrefix + "TYPE"
	EnvOSArch         = EnvOSPrefix + "ARCH"
	EnvOSID           = EnvOSPrefix + "ID"
	EnvOSCodename     = EnvOSPrefix + "CODENAME"
	EnvOSVersion      = EnvOSPrefix + "VERSION"
	EnvOSVersionMajor = EnvOSVersionPrefix + "MAJOR"
	EnvOSVersionMinor = EnvOSVersionPrefix + "MINOR"

	EnvStartedAtRFC3339    = EnvStartedAtPrefix + "RFC3339"
	EnvStartedAtUnix       = EnvStartedAtPrefix + "UNIX"
	EnvStartedAtYear       = EnvStartedAtPrefix + "YEAR"
	EnvStartedAtYearDay    = EnvStartedAtPrefix + "YEAR_DAY"
	EnvStartedAtDay        = EnvStartedAtPrefix + "DAY"
	EnvStartedAtMonth      = EnvStartedAtPrefix + "MONTH"
	EnvStartedAtMonthStr   = EnvStartedAtPrefix + "MONTH_STR"
	EnvStartedAtHour       = EnvStartedAtPrefix + "HOUR"
	EnvStartedAtMinute     = EnvStartedAtPrefix + "MINUTE"
	EnvStartedAtSecond     = EnvStartedAtPrefix + "SECOND"
	EnvStartedAtNanosecond = EnvStartedAtPrefix + "NANOSECOND"
	EnvStartedAtTimezone   = EnvStartedAtPrefix + "TIMEZONE"
	EnvStartedAtOffset     = EnvStartedAtPrefix + "OFFSET"
	EnvStartedAtWeekday    = EnvStartedAtPrefix + "WEEKDAY"
	EnvStartedAtWeekdayStr = EnvStartedAtPrefix + "WEEKDAY_STR"

	EnvNetworkIPv4         = EnvNetworkPrefix + "IPV4"
	EnvNetworkIPv6         = EnvNetworkPrefix + "IPV6"
	EnvNetworkHostname     = EnvNetworkPrefix + "HOSTNAME"
	EnvNetworkCity         = EnvNetworkPrefix + "CITY"
	EnvNetworkRegion       = EnvNetworkPrefix + "REGION"
	EnvNetworkCountry      = EnvNetworkPrefix + "COUNTRY"
	EnvNetworkASN          = EnvNetworkPrefix + "ASN"
	EnvNetworkOrganization = EnvNetworkPrefix + "ORGANIZATION"
	EnvNetworkPostal       = EnvNetworkPrefix + "POSTAL"
	EnvNetworkTimezone     = EnvNetworkPrefix + "TIMEZONE"
	EnvNetworkLatitude     = EnvNetworkPrefix + "LATITUDE"
	EnvNetworkLongitude    = EnvNetworkPrefix + "LONGITUDE"

	EnvHostname     = Prefix + "HOSTNAME"
	EnvHostnameFQDN = Prefix + "HOSTNAME_FQDN"

	EnvGitReference       = EnvGitPrefix + "REFERENCE"
	EnvGitReferenceShort  = EnvGitPrefix + "REFERENCE_SHORT"
	EnvGitReferenceType   = EnvGitPrefix + "REFERENCE_TYPE"
	EnvGitCommitHash      = EnvGitPrefix + "COMMIT_HASH"
	EnvGitCommitHashShort = EnvGitPrefix + "COMMIT_HASH_SHORT"
	EnvGitIsDirty         = EnvGitPrefix + "IS_DIRTY"

	EnvUserUID  = EnvUserPrefix + "UID"
	EnvUserGID  = EnvUserPrefix + "GID"
	EnvUserName = EnvUserPrefix + "NAME"
)

func ParameterName(name string) string {
	return EnvParameterPrefix + strings.ToUpper(name)
}

func FlagName(name string) string {
	return EnvFlagPrefix + strings.ToUpper(name)
}
