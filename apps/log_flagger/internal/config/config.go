package config

import (
	"os"
	"time"
)

const (
	GEOIP_PROVIDER_OXL         = "oxl"
	GEOIP_PROVIDER_IPINFO_LITE = "ipinfo-lite"
	GEOIP_PROVIDER_IPINFO      = "ipinfo"
	GEOIP_PROVIDER_MAXMIND     = "maxmind"
)

var (
	PATH_OPEN_BOT_DATA        = ""
	PATH_OUTPUT               = ""
	PATH_INPUT                = ""
	CSV_FIELD_CLIENT_IP       = 0
	CSV_FIELD_FP_JA4          = 1
	CSV_FIELD_UA              = 2
	LOOKUP_PTR                = true
	DNS_RESOLVE_TIMEOUT       = 200 * time.Millisecond
	GEOIP_PROVIDER            = "oxl"
	PATH_GEOIP_ASN_DB         = ""
	DEBUG                     = false
	MODE_TEST                 = os.Getenv("MODE_TEST") == "1"
	DEBUG_UA                  = "___"
	SUPPORTED_GEOIP_PROVIDERS = []string{GEOIP_PROVIDER_IPINFO_LITE, GEOIP_PROVIDER_IPINFO, GEOIP_PROVIDER_MAXMIND, GEOIP_PROVIDER_OXL}
)

const VERSION = 1.2
