package env

const (
	ChoreDir = "chore"
)

const (
	Prefix             = "CHORE_"
	EnvArgPrefix       = Prefix + "ARG_"
	EnvPathPrefix      = Prefix + "PATH_"
	EnvNetworkPrefix   = Prefix + "NETWORK_"
	EnvStartedAtPrefix = Prefix + "STARTED_AT_"
	EnvIdPrefix        = Prefix + "ID_"
	EnvIdChainPrefix   = Prefix + "CHAIN_ID_"

	EnvNamespace = Prefix + "NAMESPACE"
	EnvCaller    = Prefix + "CALLER"

	EnvPathCaller  = EnvPathPrefix + "CALLER"
	EnvPathData    = EnvPathPrefix + "DATA"
	EnvPathCache   = EnvPathPrefix + "CACHE"
	EnvPathState   = EnvPathPrefix + "STATE"
	EnvPathRuntime = EnvPathPrefix + "RUNTIME"
	EnvPathTemp    = EnvPathPrefix + "TEMP"

	EnvIdUnique        = EnvIdPrefix + "UNIQUE"
	EnvIdChainUnique   = EnvIdChainPrefix + "UNIQUE"
	EnvIdIsolated      = EnvIdPrefix + "ISOLATED"
	EnvIdChainIsolated = EnvIdChainPrefix + "ISOLATED"

	EnvMachineId = Prefix + "MACHINE_ID"

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
)
