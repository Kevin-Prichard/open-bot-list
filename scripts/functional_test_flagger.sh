#!/usr/bin/env bash

if [ -n "$1" ]
then
  GEOIP_PROVIDER="$(echo $1 | tr '[:upper:]' '[:lower:]')"
else
  GEOIP_PROVIDER="oxl"
fi

if [ -n "$2" ]
then
  DEBUG_UA="$2"
else
  DEBUG_UA='___'
fi

set -euo pipefail

cd "$(dirname "$0")/.."
BASE_DIR="$(pwd)"

export MODE_TEST=1

if [[ "$GEOIP_PROVIDER" == "ipinfo-lite" ]]
then
  PATH_GEOIP_DB='/tmp/ipinfo_lite.mmdb'
elif [[ "$GEOIP_PROVIDER" == "ipinfo" ]]
then
  PATH_GEOIP_DB='/tmp/ipinfo_asn.mmdb'
elif [[ "$GEOIP_PROVIDER" == "maxmind" ]]
then
  PATH_GEOIP_DB='/tmp/maxmind_asn.mmdb'
else
  PATH_GEOIP_DB='/tmp/oxl_geoip_asn.mmdb'
fi

PATH_OUT='/tmp/oxl-open-bot-list-out.csv'
PATH_DATA='/tmp/oxl-open-bot-list-data'
FILE_ERROR='/tmp/oxl-open-bot-list-error.txt'

rm -f "${PATH_OUT}/"* "$FILE_ERROR"

echo '### EXTRACTING TESTDATA ###'
mkdir -p "$PATH_DATA"
tar -C "$PATH_DATA" -xJvf "${BASE_DIR}/testdata/log_flagger/downloader-out.tar.xz" > /dev/null
tar -C '/tmp' -xJvf "${BASE_DIR}/testdata/log_flagger/oxl_geoip_asn.mmdb.tar.xz" > /dev/null

echo '### RUNNING ###'
build/log_flagger -csv-field-client-ip=2 \
  -csv-field-fp-ja4=4 \
  -csv-field-user-agent=5 \
  -input-file="${BASE_DIR}/testdata/log_flagger/logs.csv" \
  -data-dir="$PATH_DATA" \
  -output-file="$PATH_OUT" \
  -debug-user-agent="$DEBUG_UA" \
  -geoip-provider="$GEOIP_PROVIDER" \
  -geoip-asn-file="$PATH_GEOIP_DB" \
  -lookup-ptr

function checkFile() {
  p="$1"
  if ! [ -f "$p" ]
  then
    echo "ERROR: File missing => ${p}"
    touch "$FILE_ERROR"
  fi
  if [[ "$(cat "$p" | wc -l)" == "0" ]]
  then
    echo "ERROR: File empty => ${p}"
    touch "$FILE_ERROR"
  fi
}

function checkContent() {
  c="$1"
  if ! grep -q "$c" < "$PATH_OUT"
  then
    echo "ERROR: Expected line not found => '${c}'"
    touch "$FILE_ERROR"
  fi
}

echo ''
echo '### CHECKING ###'
checkFile "$PATH_OUT"

if [[ "$GEOIP_PROVIDER" != "oxl" ]]
then
  echo "WARNING: TESTS WILL FAIL BECAUSE GEOIP-DATA DIFFERS BETWEEN PROVIDERS!"
fi

# todo: ptr's might change and lead to test-issues

# Googlebot
checkContent '2025-11-12T13:52:11.243000+01:00,17d5b.test.oxl.app,66.249.77.224,200,t13d181300_e8a523a41297_43ade6aba3df,Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html),bot|bot_crawler|crawler_search|crawler_verified|fingerprint_crawler|fingerprint_crawler_tls_ja4|org_google|ptr_crawler|ptr_crawler_googlebot|src_as_name_Google-LLC|src_asn_15169|src_asn_crawler|src_asn_hosting|src_hosting'
# Google-Storebot
checkContent '2025-11-12T13:52:11.168000+01:00,7c8f2.test.oxl.app,66.249.83.65,200,t13d181300_e8a523a41297_43ade6aba3df,"Mozilla/5.0 (X11; Linux x86_64; Storebot-Google/1.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",bot|bot_crawler|crawler_user|crawler_verified|fingerprint_crawler|fingerprint_crawler_tls_ja4|org_google|ptr_crawler|ptr_crawler_googlebot|src_as_name_Google-LLC|src_asn_15169|src_asn_crawler|src_asn_hosting|src_hosting'
# Petalbot
checkContent '2025-11-12T13:52:11.222000+01:00,3ebcc.test.oxl.app,114.119.156.140,200,t13d311100_e8f1e7e78f70_d41ae481755e,"Mozilla/5.0 (Linux; Android 7.0;) AppleWebKit/537.36 (KHTML, like Gecko) Mobile Safari/537.36 (compatible; PetalBot;+https://webmaster.petalsearch.com/site/petalbot)",bot|bot_crawler|ptr_crawler|ptr_crawler_huawei_petalsearch|src_as_name_Huawei-Clouds|src_asn_136907|src_asn_hosting|src_hosting'
# Amazonbot
checkContent '2025-11-12T13:52:13.606000+01:00,eb84e.test.oxl.app,54.197.114.76,200,t12d430700_1ce71f0edbb1_269097225c9f,"Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; Amazonbot/0.1; +https://developer.amazon.com/support/amazonbot) Chrome/119.0.6045.214 Safari/537.36",bot|bot_crawler|crawler_ai_data|org_amazon|ptr_crawler|ptr_crawler_amazon_main|src_as_name_Amazon.com-Inc.|src_asn_14618|src_asn_hosting|src_hosting'
# Bingbot
checkContent '2025-11-12T13:52:12.216000+01:00,b2466.test.oxl.app,52.167.144.187,200,t13d2212h2_231e334592e8_36bf25f296df,"Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm) Chrome/116.0.1938.76 Safari/537.36",bot|bot_crawler|crawler_search|crawler_verified|fingerprint_crawler|fingerprint_crawler_tls_ja4|org_microsoft|ptr_crawler|ptr_crawler_bingbot|src_as_name_Microsoft-Corporation|src_asn_8075|src_asn_crawler|src_asn_hosting|src_hosting'
# Script-Bot Fingerprint-Detection
checkContent '2025-11-12T13:52:11.240000+01:00,eb84e.test.oxl.app,116.204.8.140,200,t13i140900_cbb2034c60b8_e7c285222651,This is a very serious browser!,bot|bot_script|fingerprint_script|fingerprint_script_tls_ja4|ptr_hosting|ptr_hosting_huawei_cloud|src_as_name_Huawei-Clouds|src_asn_55990|src_asn_hosting|src_hosting|src_ip_hosting'
# Script-Bot User-Agent Detection
checkContent '2025-11-12T13:52:11.144000+01:00,7c8f2.test.oxl.app,213.55.221.232,200,q13d0311h3_55b375c5d22e_f2a83c8e78ae,Python requests someversion,bot|bot_script|src_as_name_SALT-MOBILE-S.A|src_asn_15796|src_asn_isp|src_isp'
# Bot over Apple-Privacy-Relay
checkContent '2025-11-12T13:52:17.554000+01:00,b2466.test.oxl.app,172.226.108.65,301,t13d2013h2_a09f3c656075_7f0f34a4126d,"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_1) AppleWebKit/601.2.4 (KHTML, like Gecko) Version/9.0.1 Safari/601.2.4 facebookexternalhit/1.1 Facebot Twitterbot/1.0",bot|bot_crawler|crawler_socialmedia|crawler_user|org_meta|org_vpn_apple|ptr_hosting|ptr_hosting_akamai|src_as_name_Akamai-Technologies-Inc.|src_asn_36183|src_asn_hosting|src_hosting|src_ip_hosting|src_ip_vpn'
# ASN heavily used by VPN-Providers
checkContent '2025-11-12T13:52:11.105000+01:00,b2466.test.oxl.app,185.132.187.239,301,t13d1812h1_85036bcba153_d41ae481755e,"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Brave/1.63.162 Chrome/122.0.6261.111 Safari/537.36",src_as_name_Internet-Utilities-Europe-and-Asia-Limited|src_asn_206092|src_asn_hosting|src_asn_vpn|src_hosting'

echo ''
echo "CHECK OUTPUT: less ${PATH_OUT}"
echo ''

if [ -f "$FILE_ERROR" ]
then
  echo "### FAILED ###"
  exit 1
else
  echo "### SUCCESS ###"
fi
