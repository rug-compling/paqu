#!/bin/bash

if [ "`curl http://127.0.0.1:8067/up 2> /dev/null`" = "up" ]
then
    exit
fi

cd `dirname $0`

ln -s lock.$$ lock
if [ "`readlink lock`" != lock.$$ ]
then
    echo Getting lock failed
    exit
fi

killall pqserve
sleep 2
killall -KILL pqserve
sleep 2
echo Restarting pqserve
cat pqserve.out
../bin/pqserve &> pqserve.out &

rm -f lock
