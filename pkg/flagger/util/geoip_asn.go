package util

// base-code source: https://github.com/O-X-L/geoip-lookup-service

import (
	"fmt"
	"git.oxl.at/open-bot-list/pkg/flagger/config"
	"github.com/oschwald/maxminddb-golang"
	"net"
	"reflect"
	"strings"
)

// IPInfo schema: https://github.com/ipinfo/sample-database/
type IPINFO_LITE struct {
	ASN  string `maxminddb:"asn"`
	Name string `maxminddb:"as_name"`
}

type IPINFO_ASN struct {
	ASN  string `maxminddb:"asn"`
	Name string `maxminddb:"name"`
}

// MaxMind schema: https://github.com/maxmind/MaxMind-DB/tree/main/source-data
type MAXMIND_ASN struct {
	ASN  uint64 `maxminddb:"autonomous_system_number"`
	Name string `maxminddb:"autonomous_system_organization"`
}

// OXL GEOIP-ASN schema: https://github.com/O-X-L/geoip-asn/blob/latest/schema | https://github.com/O-X-L/geoip-asn/blob/latest/example
type OXL_ASN struct {
	ASN          int `maxminddb:"asn"`
	Organization struct {
		Name string `maxminddb:"name"`
	} `maxminddb:"organization"`
}

func getMapValue(dataStructure interface{}, name string) interface{} {
	return reflect.Indirect(
		reflect.ValueOf(&dataStructure),
	).Elem().Interface().(map[string]interface{})[name]
}

func getMapValueStr(dataStructure interface{}, name string) string {
	return getMapValue(dataStructure, name).(string)
}

func lookupAsnIpInfoLite(ip net.IP) (string, string) {
	data := lookupBase(ip, IPINFO_LITE{})
	if data == nil {
		return "", ""
	}
	asn := strings.TrimPrefix(getMapValueStr(data, "asn"), "AS")
	as_name := getMapValueStr(data, "as_name")
	return asn, as_name
}

func lookupAsnIpInfo(ip net.IP) (string, string) {
	data := lookupBase(ip, IPINFO_ASN{})
	if data == nil {
		return "", ""
	}
	asn := strings.TrimPrefix(getMapValueStr(data, "asn"), "AS")
	as_name := getMapValueStr(data, "as_name")
	return asn, as_name
}

func lookupAsnMaxMind(ip net.IP) (string, string) {
	data := lookupBase(ip, MAXMIND_ASN{})
	if data == nil {
		return "", ""
	}
	asn := fmt.Sprintf("%d", getMapValue(data, "autonomous_system_number"))
	as_name := getMapValueStr(data, "autonomous_system_organization")
	return asn, as_name
}

func lookupAsnOXL(ip net.IP) (string, string) {
	data := lookupBase(ip, OXL_ASN{})
	if data == nil {
		return "", ""
	}
	asn := fmt.Sprintf("%d", getMapValue(data, "asn"))
	as_name := getMapValueStr(getMapValue(data, "organization"), "name")
	return asn, as_name
}

func lookupBase(ip net.IP, dataStructure interface{}) interface{} {
	db, err := maxminddb.Open(config.PATH_GEOIP_ASN_DB)
	if err != nil {
		if config.MODE_TEST {
			fmt.Println("GeoIP: open error")
		}
		return nil
	}
	defer db.Close()

	err = db.Lookup(ip, &dataStructure)
	if err != nil {
		if config.MODE_TEST {
			fmt.Println("GeoIP: lookup error")
		}
		return nil
	}
	return dataStructure
}

func LookupGeoIPASN(clientIP string) (string, string) {
	var asn string
	var as_name string

	ip := net.ParseIP(clientIP)

	defer func() {
		if err := recover(); err != nil {
			if config.MODE_TEST {
				fmt.Printf("GeoIP: unexpected error: %v\n", err)
			}
		}
	}()

	switch config.GEOIP_PROVIDER {
	case config.GEOIP_PROVIDER_IPINFO:
		asn, as_name = lookupAsnIpInfo(ip)
	case config.GEOIP_PROVIDER_IPINFO_LITE:
		asn, as_name = lookupAsnIpInfoLite(ip)
	case config.GEOIP_PROVIDER_MAXMIND:
		asn, as_name = lookupAsnMaxMind(ip)
	default:
		asn, as_name = lookupAsnOXL(ip)
	}

	return asn, as_name
}
