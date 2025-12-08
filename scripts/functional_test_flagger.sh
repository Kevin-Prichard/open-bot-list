#!/usr/bin/env bash

if [ -n "$1" ]
then
  DEBUG_UA="$1"
else
  DEBUG_UA='___'
fi

set -euo pipefail

cd "$(dirname "$0")/.."
BASE_DIR="$(pwd)"

export MODE_TEST=1

PATH_OUT='/tmp/oxl-open-bot-list-out.csv'
PATH_DATA='/tmp/oxl-open-bot-list-data'
FILE_ERROR='/tmp/oxl-open-bot-list-error.txt'

rm -f "${PATH_OUT}/"* "$FILE_ERROR"

echo '### EXTRACTING TESTDATA ###'
mkdir -p "$PATH_DATA"
tar -C "$PATH_DATA" -xJvf "${BASE_DIR}/testdata/log_flagger/downloader-out.tar.xz" > /dev/null

echo '### RUNNING ###'
build/log_flagger -csv-field-client-ip=2 \
  -csv-field-fp-ja4=4 \
  -csv-field-user-agent=5 \
  -input-file="${BASE_DIR}/testdata/log_flagger/logs.csv" \
  -data-dir="$PATH_DATA" \
  -output-file="$PATH_OUT" \
  -debug-user-agent="$DEBUG_UA"


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

echo '### CHECKING ###'
checkFile "$PATH_OUT"

# Googlebot
checkContent '2025-11-12T13:52:11.243000+01:00,17d5b.test.oxl.app,66.249.77.224,200,t13d181300_e8a523a41297_43ade6aba3df,Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html),bot|bot_crawler|crawler_search|crawler_verified|fingerprint_crawler|org_google'
# Google-Storebot
checkContent '2025-11-12T13:52:11.168000+01:00,7c8f2.test.oxl.app,66.249.83.65,200,t13d181300_e8a523a41297_43ade6aba3df,"Mozilla/5.0 (X11; Linux x86_64; Storebot-Google/1.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",bot|bot_crawler|crawler_user|crawler_verified|fingerprint_crawler|org_google'
# Petalbot
checkContent '2025-11-12T13:52:11.222000+01:00,3ebcc.test.oxl.app,114.119.156.140,200,t13d311100_e8f1e7e78f70_d41ae481755e,"Mozilla/5.0 (Linux; Android 7.0;) AppleWebKit/537.36 (KHTML, like Gecko) Mobile Safari/537.36 (compatible; PetalBot;+https://webmaster.petalsearch.com/site/petalbot)",bot|bot_crawler'
# Amazonbot
checkContent '2025-11-12T13:52:11.207000+01:00,3ebcc.test.oxl.app,52.205.141.124,200,t12d430700_1ce71f0edbb1_269097225c9f,"Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; Amazonbot/0.1; +https://developer.amazon.com/support/amazonbot) Chrome/119.0.6045.214 Safari/537.36",bot|bot_crawler|crawler_ai_data|org_amazon'
# Bingbot
checkContent '2025-11-12T13:52:12.216000+01:00,b2466.test.oxl.app,52.167.144.187,200,t13d2212h2_231e334592e8_36bf25f296df,"Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm) Chrome/116.0.1938.76 Safari/537.36",bot|bot_crawler|crawler_search|crawler_verified|fingerprint_crawler|org_microsoft'
# Script-Bot Fingerprint-Detection
checkContent '2025-11-12T13:52:11.240000+01:00,eb84e.test.oxl.app,116.204.8.140,200,t13i140900_cbb2034c60b8_e7c285222651,This is a very serious browser!,bot|bot_script|fingerprint_script'
# Script-Bot User-Agent Detection
checkContent '2025-11-12T13:52:11.144000+01:00,7c8f2.test.oxl.app,213.55.221.232,200,q13d0311h3_55b375c5d22e_f2a83c8e78ae,Python requests someversion,bot|bot_script'

if [ -f "$FILE_ERROR" ]
then
  echo "### FAILED ###"
  exit 1
else
  echo "### SUCCESS ###"
fi
