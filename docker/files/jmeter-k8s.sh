#!/usr/bin/env bash
# Description: Run Controller or Workes nodes
# Author: Diego Rodrigues diegofull at gmail.com
#               2022-02-26
# Version: 0.1 : start

# Global variables
nbInjectors=1
controller=0
worker=0
pathJMX="/mnt/e/Diego/Projetos/jmeter-k8s-starterkit/docker/files/scenario"

function help_usage(){
cat <<EOF
$(basename $0) -j -n -c|-w  [-m|-r|-i]
  -j filename.jmx
  -n namespace
  -c run in controller mode
  -w run in worker mode
  -m flag to copy fragmented jmx present in scenario/project/module if you use include controller and external test fragment
  -i <injectorNumber> to scale workes pods to the desired number of JMeter injectors
  -r flag to enable report generation at the end of the test

Examples:
Run in controller mode with 2 Worker and report enabled
    $(basename $0) -j my_file.jmx -n default -c -i 2 -r

Run in worker mode with 2 Worker and report enabled
    $(basename $0) -j my_file.jmx -n default -w -i 2 -r

EOF
exit 0
}

function exit_system(){
  msg="${1:-}"
  shift
  exitCode="${1:-0}"
  echo -e "${msg}"
  exit ${exitCode}
}

function controller_start(){

    echo "teste"

}

function worker_start(){
  echo "worker 1"
}

function validation_options(){

  # Check if only one mode is selected
  [ "${controller}" -eq 1 -a "${worker}" -eq 1  ] && exit_system "Only one mode is accepted per time" 1

  # Check if one mode is selected
  [ "${controller}" -eq 0 -a "${worker}" -eq 0  ] && exit_system "ONE mode should be selected" 1

  # Check if namespace is defined
  [ -z "${namespace}" ] && exit_system "The NAMESPACE is required" 1

  # Check if JMX file is defined
  [ -z "${jmx}" ] && exit_system "The JMX is required" 1

  # Check if JMX file exist in the path
  [[ -e "${pathJMX}/${jmx}" ]] || exit_system "JMX file does not exist ${pathJMX}/${jmx}" 1

}

function main(){

  # Check if arguments were sent
  [[ -z "${@}" ]] && help_usage

  while getopts 'i:mj:hcwrn:' option;
  do
    case $option in
      n) namespace=${OPTARG} ;;
      c) controller=1 ;;
      w) worker=1 ;;
      m) module=1 ;;
      r) enable_report=1 ;;
      j) jmx=${OPTARG} ;;
      i) nbInjectors=${OPTARG} ;;
      h) help_usage ;;
      *) help_usage ;;
    esac
  done

  # Check ARGS logics
  validation_options

  # Entering in the mode selected
  [[ "${controller}" -eq 1 ]] && controller_start

  [[ "${worker}" -eq 1 ]] && worker_start


}

main $@
