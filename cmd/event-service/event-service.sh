#!/bin/bash
appname=com.github.craigdfrench.event
pidpath=/tmp
repobase=$GOPATH/src/github.com/craigdfrench/event
installed_apps=`find $repobase -name '*.go' | xargs grep  -n '^package main$' | cut -d ':' -f 1-1 | xargs -n1 dirname | xargs -n1 basename`

[[ $2 == "--verbose" ]] && verbose=1
[[ $2 == "-v" ]] && verbose=1

function show_help {
cat<<HERE
event-service USAGE
$0 [cmd] [option]

where [cmd] is one of:
   restart: stop the running services and start them again
   start: if services not running, stop the services 
   stop: if services running, stop the services
   status: show the status of the services
   help: this help text

where [options]
   --verbose or -v: extra logging information
HERE
}

function build {
    cd $repobase
    go generate ./...
    go install ./...
}

function verbose {
    if [ ${verbose} ]
    then
        echo $*
    fi
}

function iterate_apps_by_runstate {
    local sense=$1 cmd=$2 accumulator=0 pidfile appcmd service
    for service in ${installed_apps[@]}; do
        pidfile=${pidpath}/${appname}.${service}.pid
        [[ -e $pidfile ]] && pidFilePresent="running" || pidFilePresent="not_running"
        if [[ $sense == $pidFilePresent ]]
        then
            verbose "iterate_apps_by_runstate executing <$cmd> <$pidfile> <$service>"
            $cmd $pidfile $service
            accumulator=$((accumulator+1))
        fi
    done
    return $accumulator
}

function launch_service {
    local app=$2 pidfile=$1
    echo Starting service $app
    $app &
    echo $! > $pidfile
    verbose Stored PID for $app in $pidfile
}

function count {
    verbose count called
}

function print_ps {
    local pidfile=$1
    ps `cat $pidfile` | tail -1
}

function print_notrunning {
    local pidfile=$1 service=$2
    echo Service $2 is not running
}

function kill_service {
    local pidfile=$1 service=$2
    echo Shutting down service $service
    kill `cat $pidfile`
    echo Removing pid file $pidfile
    rm $pidfile
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
