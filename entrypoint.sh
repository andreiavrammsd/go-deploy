#!/bin/bash

source ~/.profile

cmd="$@"

if [ "$cmd" == "build" ]
then
    cd /deploy/src/deploy
    rm -f deploy
    go get
else
    # Execute the command sent via Dockerfile
    exec "$@"
fi
