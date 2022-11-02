#!/usr/bin/env bash

source ./style_info.cfg
source ./path_info.cfg
source ./function.sh

service_name=${service_filename[0]}
  list=$(cat $config_path | grep -w apiPort | awk -F '[:]' '{print $NF}')
  list_to_string $list
  for j in ${ports_array}; do
    port=$(ss -tunlp| grep excel | awk '{print $5}' | grep -w ${j} | awk -F '[:]' '{print $NF}')
    if [[ ${port} -ne ${j} ]]; then
      echo -e ${YELLOW_PREFIX}${i}${COLOR_SUFFIX}${RED_PREFIX}"service does not start normally,not initiated port is "${COLOR_SUFFIX}${YELLOW_PREFIX}${j}${COLOR_SUFFIX}
      echo -e ${RED_PREFIX}"please check ../logs/Excel-Props.log "${COLOR_SUFFIX}
      exit -1
    else
      echo -e ${j}${GREEN_PREFIX}" port has been listening,belongs service is "apiPort${COLOR_SUFFIX}
    fi
  done





#check=$(ps aux | grep -w ./${cron_task_name} | grep -v grep | wc -l)
#if [ $check -ge 1 ]; then
#  echo -e ${GREEN_PREFIX}"none  port has been listening,belongs service is openImCronTask"${COLOR_SUFFIX}
#else
#  echo -e ${RED_PREFIX}"cron_task_name service does not start normally"${COLOR_SUFFIX}
#        echo -e ${RED_PREFIX}"please check ../logs/openIM.log "${COLOR_SUFFIX}
#      exit -1
#fi
#
echo -e ${YELLOW_PREFIX}"all services launch success"${COLOR_SUFFIX}
