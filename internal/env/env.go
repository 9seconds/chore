package env

import "strings"

const (
	Prefix              = "CHORE_"
	ParameterPrefix     = Prefix + "P_"
	ParameterPrefixList = Prefix + "PL_"
	FlagPrefix          = Prefix + "F_"
	PathPrefix          = Prefix + "PATH_"
	NetworkPrefix       = Prefix + "NETWORK_"
	StartedAtPrefix     = Prefix + "STARTED_AT_"
	IDPrefix            = Prefix + "ID_"
	IDChainPrefix       = Prefix + "CHAIN_ID_"
	OSPrefix            = Prefix + "OS_"
	OSVersionPrefix     = OSPrefix + "VERSION_"
	GitPrefix           = Prefix + "GIT_"
	UserPrefix          = Prefix + "USER_"

	Debug     = Prefix + "DEBUG"
	MachineID = Prefix + "MACHINE_ID"

	Namespace = Prefix + "NAMESPACE"
	Caller    = Prefix + "CALLER"
	Self      = Prefix + "SELF"
	Slug      = Prefix + "SLUG"

	PathCaller = PathPrefix + "CALLER"
	PathData   = PathPrefix + "DATA"
	PathCache  = PathPrefix + "CACHE"
	PathState  = PathPrefix + "STATE"
	PathTemp   = PathPrefix + "TEMP"

	IDRun           = IDPrefix + "RUN"
	IDChainRun      = IDChainPrefix + "RUN"
	IDIsolated      = IDPrefix + "ISOLATED"
	IDChainIsolated = IDChainPrefix + "ISOLATED"

	OSType         = OSPrefix + "TYPE"
	OSArch         = OSPrefix + "ARCH"
	OSID           = OSPrefix + "ID"
	OSCodename     = OSPrefix + "CODENAME"
	OSVersion      = OSPrefix + "VERSION"
	OSVersionMajor = OSVersionPrefix + "MAJOR"
	OSVersionMinor = OSVersionPrefix + "MINOR"

	StartedAtRFC3339    = StartedAtPrefix + "RFC3339"
	StartedAtUnix       = StartedAtPrefix + "UNIX"
	StartedAtYear       = StartedAtPrefix + "YEAR"
	StartedAtYearDay    = StartedAtPrefix + "YEAR_DAY"
	StartedAtDay        = StartedAtPrefix + "DAY"
	StartedAtMonth      = StartedAtPrefix + "MONTH"
	StartedAtMonthStr   = StartedAtPrefix + "MONTH_STR"
	StartedAtHour       = StartedAtPrefix + "HOUR"
	StartedAtMinute     = StartedAtPrefix + "MINUTE"
	StartedAtSecond     = StartedAtPrefix + "SECOND"
	StartedAtNanosecond = StartedAtPrefix + "NANOSECOND"
	StartedAtTimezone   = StartedAtPrefix + "TIMEZONE"
	StartedAtOffset     = StartedAtPrefix + "OFFSET"
	StartedAtWeekday    = StartedAtPrefix + "WEEKDAY"
	StartedAtWeekdayStr = StartedAtPrefix + "WEEKDAY_STR"

	NetworkIPv4         = NetworkPrefix + "IPV4"
	NetworkIPv6         = NetworkPrefix + "IPV6"
	NetworkHostname     = NetworkPrefix + "HOSTNAME"
	NetworkCity         = NetworkPrefix + "CITY"
	NetworkRegion       = NetworkPrefix + "REGION"
	NetworkCountry      = NetworkPrefix + "COUNTRY"
	NetworkASN          = NetworkPrefix + "ASN"
	NetworkOrganization = NetworkPrefix + "ORGANIZATION"
	NetworkPostal       = NetworkPrefix + "POSTAL"
	NetworkTimezone     = NetworkPrefix + "TIMEZONE"
	NetworkLatitude     = NetworkPrefix + "LATITUDE"
	NetworkLongitude    = NetworkPrefix + "LONGITUDE"

	Hostname     = Prefix + "HOSTNAME"
	HostnameFQDN = Prefix + "HOSTNAME_FQDN"

	GitReference       = GitPrefix + "REFERENCE"
	GitReferenceShort  = GitPrefix + "REFERENCE_SHORT"
	GitReferenceType   = GitPrefix + "REFERENCE_TYPE"
	GitCommitHash      = GitPrefix + "COMMIT_HASH"
	GitCommitHashShort = GitPrefix + "COMMIT_HASH_SHORT"
	GitIsDirty         = GitPrefix + "IS_DIRTY"

	UserUID  = UserPrefix + "UID"
	UserGID  = UserPrefix + "GID"
	UserName = UserPrefix + "NAME"
)

func ParameterName(name string) string {
	return ParameterPrefix + strings.ToUpper(name)
}

func ParameterNameList(name string) string {
	return ParameterPrefixList + strings.ToUpper(name)
}

func FlagName(name string) string {
	return FlagPrefix + strings.ToUpper(name)
}
