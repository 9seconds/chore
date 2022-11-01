package env

const (
	ChoreDir = "chore"
)

const (
	Prefix             = "CHORE_"
	EnvArgPrefix       = Prefix + "ARG_"
	EnvNetworkPrefix   = Prefix + "NETWORK_"
	EnvStartedAtPrefix = Prefix + "STARTED_AT_"

	EnvNamespace     = Prefix + "NAMESPACE"
	EnvCaller        = Prefix + "CALLER"
	EnvCallerPath    = Prefix + "CALLER_PATH"
	EnvPersistentDir = Prefix + "PERSISTENT_DIR"
	EnvTempDir       = Prefix + "TMP_DIR"
	EnvCorrelateId   = Prefix + "CORRELATE_ID"
	EnvRunId         = Prefix + "RUN_ID"
	EnvCacheId       = Prefix + "CACHE_ID"
	EnvMachineId     = Prefix + "MACHINE_ID"

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
	EnvNetworkOrganization = EnvNetworkPrefix + "ORGANIZATION"
	EnvNetworkPostal       = EnvNetworkPrefix + "POSTAL"
	EnvNetworkTimezone     = EnvNetworkPrefix + "TIMEZONE"
	EnvNetworkLatitude     = EnvNetworkPrefix + "LATITUDE"
	EnvNetworkLongitude    = EnvNetworkPrefix + "LONGITUDE"
)
