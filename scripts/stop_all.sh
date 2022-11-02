#!/usr/bin/env bash
#fixme This script is to stop the service

source ./style_info.cfg
source ./path_info.cfg

service_name=${service_filename[0]}

  #Check whether the service exists
  name="ps -aux |grep -w ${service_name} |grep -v grep"
  count="${name}| wc -l"
  if [ $(eval ${count}) -gt 0 ]; then
    pid="${name}| awk '{print \$2}'"
    echo -e "${SKY_BLUE_PREFIX}Killing service:${service_name} pid:$(eval $pid)${COLOR_SUFFIX}"
    #kill the service that existed
    kill -9 $(eval $pid)
    echo -e "${SKY_BLUE_PREFIX}service:${service_name} was killed ${COLOR_SUFFIX}"
  fi
