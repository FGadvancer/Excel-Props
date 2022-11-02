#!/usr/bin/env bash
#fixme This script is to stop the service

source ./style_info.cfg
source ./path_info.cfg

service_name=excel_props_api

  #Check whether the service exists

  check=$(ps aux | grep -w ./${service_name} | grep -v grep | wc -l)
if [ $check -ge 1 ]; then
  oldPid=$(ps aux | grep -w ./${service_name} | grep -v grep | awk '{print $2}')
      echo -e "${SKY_BLUE_PREFIX}Killing service:${service_name} pid:${oldPid}${COLOR_SUFFIX}"
      kill -9 ${oldPid}
      echo -e "${SKY_BLUE_PREFIX}service:${service_name} was killed ${COLOR_SUFFIX}"
fi

