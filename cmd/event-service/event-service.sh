#!/bin/bash -x 

APPNAME=com.github.craigdfrench.event
PIDPATH=/tmp
REPOBASE=$GOPATH/src/github.com/craigdfrench/event
GO_INSTALL_DIRS=( daemon web )
INSTALLED_APPS=`find . -name '*.go' | xargs grep  -n '^package main$' | cut -d ':' -f 1-1 | xargs -n1 dirname | xargs -n1 basename`
VERBOSE=$2
function show_help {
    verbose verbose text here verbose set to $VERBOSE
    echo ${INSTALLED_APPS} are installed apps
    echo $0 cmd 
    echo where cmd is
    echo restart
    echo start
    echo stop
    echo status
    echo help
}

function build_goinstallables {
    cd $REPOBASE
    go generate ./...
    go install ./...
}

function run_apps {
    app=$1
    echo Starting service $app  
    app &
    echo $! > ${PIDPATH}/${APPNAME}.$app.pid
    verbose ${PIDPATH}/${APPNAME}.$app.pid 
}

function verbose {
    if [ ${VERBOSE} ] 
    then
        echo $*
    fi
}

function run_goinstallables {
    cd $REPOBASE
    iterate_apps run_apps 
}

function iterate_apps {
    action=$1; parameter=$2; accumulator=$3 
    verbose iterate_apps called with $action
    for appcmd in ${INSTALLED_APPS[@]}; do
        verbose "action <$action> appcmd <$appcmd> param <$parameter> acc<$accumulator>" 
        ${action} $appcmd $parameter $accumulator
        accumulator=$?
        verbose accumulator is $accumulator
    done
    return $accumulator
}

function running_cmd {
    pidfile=${PIDPATH}/${APPNAME}.$1.pid
    appcmd=$1;cmd=$2;acc=$3
    verbose "inside running_cmd cmd <$cmd> appcmd <$appcmd> acc <$acc>" 
    if [ -e $pidfile ]
    then    
        acc=$((acc+1))   
        if [ $cmd ] 
        then
            $cmd $pidfile $appcmd 
        fi
    fi
    return $acc
}

function not_running_cmd {
    pidfile=${PIDPATH}/${APPNAME}.$1.pid
    appcmd=$1
    cmd=$2
    acc=$3
    verbose "inside not_running_cmd cmd <$cmd> appcmd <$appcmd>  acc <$acc>" 

    if [ -e $pidfile ]
    then
        verbose "not_running_cmd" found $pidfile running
    else    
        acc=$((acc+1))   
        if [ $cmd ] 
        then
            $cmd $pidfile $appcmd 
        fi
    fi
    return $acc
}

function noop {
    verbose inside noop
}

function count_running {
    iterate_apps running_cmd noop 
    return $? 
}

function action_running { 
    action=$1
    iterate_apps running_cmd $action
    return $?     
}

function action_notrunning {
    action=$1
    iterate_apps not_running_cmd $action
    return $?     
}

function check_goinstallables {
    servicecount=0
    for dir in ${INSTALLED_APPS[@]}; do
        echo Checking for service event-$dir  
        if [ -e ${PIDPATH}/${APPNAME}.event-$dir.pid ]
        then
            if [ $# -gt 0 ]
            then
                $1 ${PIDPATH}/${APPNAME}.event-$dir.pid event-$dir
            fi     
            servicecount=$((servicecount+1))
        fi
        echo ${PIDPATH}/${APPNAME}.event-$dir.pid 
    done
    return $servicecount
}



function start_service {
    build_goinstallables
    run_goinstallables
}

function stop_service {
    action_running kill_service  
}

function print_ps {
    ps `cat $1` | tail -1
}

function print_notrunning {
    echo Service $2 is not running
}

function kill_service {
    echo Shutting down service $2
    kill `cat $1`
    echo Removing pid file $1
    rm $1
}

case "$1" in
status)    
    action_running print_ps
    echo $? services running
    action_notrunning print_notrunning
    echo $? services not_running
    ;;
start) 
    count_running 
    if [ $? -gt 0 ]
    then 
        echo "Event Services running."
        echo "Cannot start."
        exit 0
    else
        start_service
    fi
    ;;
stop)
    count_running
    if [ $? -eq 0 ]
    then 
        echo "Event Services not running."
        echo "Cannot stop."
        exit 0
    else
        stop_service
    fi
    ;;
restart)
    count_running
    if [ $? -eq 0 ]
    then 
        echo "Event Services not running."
        echo "Cannot restart."
        exit 0
    else
        stop_service
        start_service
    fi
    ;;
help)
    show_help
    ;;
*)
    echo ERROR: Unknown command
    show_help
    ;;
esac

exit 0
