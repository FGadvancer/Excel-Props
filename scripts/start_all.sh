#!/usr/bin/env bash
#Include shell font styles and some basic information
source ./style_info.cfg
source ./path_info.cfg
source ./function.sh
ulimit -n 200000

service_name=excel_props_api
list1=$(cat $config_path | grep apiPort | awk -F '[:]' '{print $NF}')

list_to_string $list1
api_ports=($ports_array)

#Check if the service exists
#If it is exists,kill this process
check=$(ps aux | grep -w ./${service_name} | grep -v grep | wc -l)
if [ $check -ge 1 ]; then
  oldPid=$(ps aux | grep -w ./${service_name} | grep -v grep | awk '{print $2}')
    kill -9 ${oldPid}
fi
#Waiting port recycling
sleep 1
cd ${msg_gateway_binary_root}
for ((i = 0; i < ${#ws_ports[@]}; i++)); do
  nohup ./${service_name}  >>../logs/Excel-Props.log 2>&1 &
done

#Check launched service process
sleep 3
check=$(ps aux | grep -w ./${service_name} | grep -v grep | wc -l)
allPorts=""
if [ $check -ge 1 ]; then
  allNewPid=$(ps aux | grep -w ./${service_name} | grep -v grep | awk '{print $2}')
  for i in $allNewPid; do
    ports=$(netstat -netulp | grep -w ${i} | awk '{print $4}' | awk -F '[:]' '{print $NF}')
      allPorts=${allPorts}"$ports "
  done
  echo -e ${SKY_BLUE_PREFIX}"SERVICE START SUCCESS"${COLOR_SUFFIX}
  echo -e ${SKY_BLUE_PREFIX}"SERVICE_NAME: "${COLOR_SUFFIX}${YELLOW_PREFIX}${service_name}${COLOR_SUFFIX}
  echo -e ${SKY_BLUE_PREFIX}"PID: "${COLOR_SUFFIX}${YELLOW_PREFIX}${allNewPid}${COLOR_SUFFIX}
  echo -e ${SKY_BLUE_PREFIX}"LISTENING_PORT: "${COLOR_SUFFIX}${YELLOW_PREFIX}${allPorts}${COLOR_SUFFIX}
else
  echo -e ${YELLOW_PREFIX}${service_name}${COLOR_SUFFIX}${RED_PREFIX}"SERVICE START ERROR, PLEASE CHECK openIM.log"${COLOR_SUFFIX}
fi
