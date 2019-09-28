#!/bin/bash 

APPNAME=com.github.craigdfrench.event-service
PIDPATH=/tmp
REPOBASE=$GOPATH/src/github.com/craigdfrench/event-service
MAKEFILE_DIRS=( run-scripts contracts web-ui )
GO_INSTALL_DIRS=( storage web )

function build_makefiles {
    for dir in ${MAKEFILE_DIRS[@]}; do 
        cd $REPOBASE/$dir
        echo Executing in $REPOBASE/$dir:  makefile 
        make 
    done
}

function build_goinstallables {
    cd $REPOBASE
    for dir in ${GO_INSTALL_DIRS[@]}; do
        echo Executing: go install $REPOBASE/$dir/event-$dir.go  
        go install $REPOBASE/$dir/event-$dir.go     
    done
}

function run_goinstallables {
    cd $REPOBASE
    for dir in ${GO_INSTALL_DIRS[@]}; do
        echo Starting service event-$dir  
        event-$dir &
        echo $! > ${PIDPATH}/${APPNAME}.event-$dir.pid
        echo ${PIDPATH}/${APPNAME}.event-$dir.pid 
    done
}

function check_goinstallables {
    servicecount=0
    for dir in ${GO_INSTALL_DIRS[@]}; do
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
function build_service {
    build_goinstallables
    build_makefiles
}

function start_service {
    build_service
    run_goinstallables
}

function stop_service {
    check_goinstallables kill_service
}

function print_ps {
    ps `cat $1`
}

function kill_service {
    echo Shutting down service $2
    kill `cat $1`
    echo Removing pid file $1
    rm $1
}

case "$1" in
status)
    check_goinstallables print_ps 
    ;;
start) 
    check_goinstallables
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
    check_goinstallables
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
    check_goinstallables
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
esac

exit 0
