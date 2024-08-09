#!/bin/bash
# DEBUG
if [ "$DEBUG" == "true" ]; then
    /usr/bin/goosefs-cli2api -d
else
    /usr/bin/goosefs-cli2api
fi