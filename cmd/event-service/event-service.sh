#!/bin/bash   

APPNAME=com.github.craigdfrench.event
PIDPATH=/tmp
REPOBASE=$GOPATH/src/github.com/craigdfrench/event
INSTALLED_APPS=`find $REPOBASE -name '*.go' | xargs grep  -n '^package main$' | cut -d ':' -f 1-1 | xargs -n1 dirname | xargs -n1 basename`
[[ $2 == "--verbose" ]] && VERBOSE=1
[[ $2 == "-v" ]] && VERBOSE=1 

echo VERBOSE is $VERBOSE
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

function build {
    cd $REPOBASE
    go generate ./...
    go install ./...
}

function run_apps {
    app=$2;pidfile=$1
    echo Starting service $app  
    $app &
    echo $! > $pidfile
    verbose Stored PID for $app in $pidfile 
}

function verbose {
    if [ ${VERBOSE} ] 
    then
        echo $*
    fi
}

function iterate_apps {
    action=$1; accumulator=0
    verbose iterate_apps called with $action
    for appcmd in ${INSTALLED_APPS[@]}; do
        verbose "iterate_apps:${action} $appcmd $accumulator ${@:2}"
        ${action} $appcmd $accumulator "${@:2}"
        accumulator=$?
        verbose accumulator=$accumulator
    done
    return $accumulator
}

# app_callbacks are called with 2 parameters: path_to_pidfile_for_component name_of_component
# iterate_apps_by_runstate [running|notrunning] app_callback ...(additional parameters)
function iterate_apps_by_runstate {
    runstate=$1
    iterate_apps conditional_cmd $runstate "${@:2}"
    return $?
}

# To use, call iterate_apps conditional_cmd [running|not_running] 
function conditional_cmd {
    pidfile=${PIDPATH}/${APPNAME}.$1.pid
    appcmd=$1;accumulator=$2;cmd=$4;sense=$3
    [[ -e $pidfile ]] && pidFilePresent="running" || pidFilePresent="not_running" 
    verbose echo "inside <${@}> conditional_cmd cmd <$cmd> appcmd <$appcmd> accumulator <$accumulator> sense <$sense> pidFilePresent <$pidFilePresent>"

    if [[ $sense == $pidFilePresent ]]
    then    
        verbose "conditional_cmd executing <$cmd> <$pidfile> <$appcmd> "
        $cmd $pidfile $appcmd
        accumulator=$((accumulator+1))
    fi 
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

function count {
    verbose inside count
}

function count_running {
    iterate_apps running_cmd noop 
    return $? 
}

function action_running { 
    action=$1
    iterate_apps conditional_cmd 0 $action running
    return $?     
}

function action_notrunning {
    action=$1
    iterate_apps_by_runstate not_running $action
    return $?     
}

function start_service {
    build_goinstallables
    iterate_apps conditional_cmd 0 running run_apps
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
    iterate_apps_by_runstate running print_ps
    echo $? services running
    iterate_apps_by_runstate not_running print_notrunning
    echo $? services not_running
    ;;
build) 
    build
    ;;
start) 
    iterate_apps_by_runstate running count 
    if [ $? -gt 0 ]
    then 
        echo "Event Services already running."
        exit 1
    else
        iterate_apps_by_runstate not_running run_apps
    fi
    ;;
stop)
    iterate_apps_by_runstate running count
    if [ $? -eq 0 ]
    then 
        echo "Event Services not running."
        exit 1
    else
        iterate_apps_by_runstate running kill_service
    fi
    ;;
restart)
    iterate_apps_by_runstate running count
    if [ $? -eq 0 ]
    then 
        echo "Event Services not running."
        exit 1
    else
        iterate_apps_by_runstate running kill_service
        iterate_apps_by_runstate not_running run_apps
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
