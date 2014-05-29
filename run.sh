#!/bin/bash

cd `dirname $0`

if [ "`curl http://127.0.0.1:8067/up 2> /dev/null`" != "up" ]
then
    killall pqserve
    sleep 2
    killall -KILL pqserve
    sleep 2
    echo Restarting pqserve
    cat pqserve.out
    pqserve/pqserve &> pqserve.out &
fi
