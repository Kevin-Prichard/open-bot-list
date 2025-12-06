#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")/.."

export MODE_TEST=1

PATH_RUN='/tmp/oxl-open-bot-list-run'
PATH_OUT='/tmp/oxl-open-bot-list-out'
FILE_ERROR='/tmp/oxl-open-bot-list-error.txt'

RUN_FILES_EXIST=(
  'match_ip_net__overall.csv' 'match_ip_net_ai.csv' 'match_ptr_crawler.csv'
  'match_user_agent__overall.csv' 'match_user_agent_crawler.csv'
  'iplist_crawler_src_net_crawler_microsoft_bing_0' 'iplist_crawler_src_net_crawler_google_common_0'
  'iplist_vpn_src_net_vpn_apple_privacyrelay_0' 'iplist_proxy_src_net_proxy_tor_0'
  'iplist_cdn_src_net_cdn_bunnyway_1' 'iplist_ai_src_net_crawler_openai_aidata_0'
)

OUT_FILES_EXIST=(
  'fingerprint_script.lst' 'fingerprint_script.map' 'fingerprint_script_tls_ja4.lst'
  'http_user_agent_ai_sub.lst' 'http_user_agent_ai_sub.map'
  'src_net_crawler_google_common_all.lst' 'src_net_crawler_google_common_net4.lst'
  'src_net_crawler_google_common_net6.lst' 'src_net_crawler_google_special_net4.lst'
  'src_net_crawler_microsoft_bing_net4.lst'
  'src_net_proxy_tor_net4.lst'
  'src_net_vpn_apple_privacyrelay_net4.lst' 'src_net_vpn_apple_privacyrelay_net6.lst'
)

rm -f "${PATH_RUN}/"* "${PATH_OUT}/"* "$FILE_ERROR"

build/downloader -output-dir="$PATH_OUT" -runtime-dir="$PATH_RUN"

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

for f in "${RUN_FILES_EXIST[@]}"
do
  checkFile "${PATH_RUN}/${f}"
done

for f in "${OUT_FILES_EXIST[@]}"
do
  checkFile "${PATH_OUT}/${f}"
done

if [ -f "$FILE_ERROR" ]
then
  echo "### FAILED ###"
  exit 1
else
  echo "### SUCCESS ###"
fi
