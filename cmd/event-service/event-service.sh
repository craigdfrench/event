#!/bin/bash
appname=com.github.craigdfrench.event
pidpath=/tmp
repobase=$GOPATH/src/github.com/craigdfrench/event
installed_apps=`find $repobase -name '*.go' | xargs grep  -n '^package main$' | cut -d ':' -f 1-1 | xargs -n1 dirname | xargs -n1 basename`

[[ $2 == "--verbose" ]] && VERBOSE=1
[[ $2 == "-v" ]] && VERBOSE=1

function show_help {
    echo ${installed_apps} are installed apps
    echo $0 cmd
    echo where cmd is
    echo restart
    echo start
    echo stop
    echo status
    echo help
}

function build {
    cd $repobase
    go generate ./...
    go install ./...
}

function verbose {
    if [ ${VERBOSE} ]
    then
        echo $*
    fi
}

function iterate_apps {
    local action=$1; local accumulator=0
    for appcmd in ${installed_apps[@]}; do
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
    local runstate=$1
    iterate_apps conditional_cmd $runstate "${@:2}"
    return $?
}

# To use, call iterate_apps conditional_cmd [running|not_running]
function conditional_cmd {
    local pidfile=${pidpath}/${appname}.$1.pid
    local appcmd=$1;local accumulator=$2;local cmd=$4;local sense=$3
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

function launch_service {
    local app=$2;local pidfile=$1
    echo Starting service $app
    $app &
    echo $! > $pidfile
    verbose Stored PID for $app in $pidfile
}

function count {
    verbose count called
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
            iterate_apps_by_runstate not_running launch_service
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
            iterate_apps_by_runstate not_running launch_service
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
