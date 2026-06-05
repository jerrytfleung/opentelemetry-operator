#!/bin/sh
# Clone-container init script for PHP auto-instrumentation.
# Runs as a clone of the application container so it can detect PHP API version,
# PHP extension directory and glibc / musl standard C library. Copies the result
# into a shared volume for the agent init container to pick up.
#
set -e

DIR="/otel-auto-instrumentation-php-clone"

# Always wrap variables in double quotes to handle paths with spaces safely
if [ -d "$DIR" ]; then
    echo "1 Success: $DIR exists and is a directory."
else
    echo "1 Error: $DIR does not exist."
fi

thread_safety=non-zts
if [ "enabled" = "$(php -i | grep "^Thread Safety => " | awk '{print $4}')" ]; then
    thread_safety=zts
fi
echo "$thread_safety" > /otel-auto-instrumentation-php-clone/thread_safety.txt

api=$(php -i | grep "^PHP API => " | awk '{print $4}')
echo "$api" > /otel-auto-instrumentation-php-clone/api.txt

# check if alpine
standard_c_lib=glibc
if [ -f /etc/alpine-release ]; then
    standard_c_lib=musl
fi
echo "$standard_c_lib" > /otel-auto-instrumentation-php-clone/standard_c_lib.txt

t=$(cat /otel-auto-instrumentation-php-clone/thread_safety.txt)
echo "Read thread_safety: $t"

a=$(cat /otel-auto-instrumentation-php-clone/api.txt)
echo "Read api: $a"

s=$(cat /otel-auto-instrumentation-php-clone/standard_c_lib.txt)
echo "Read standard_c_lib: $s"

# Always wrap variables in double quotes to handle paths with spaces safely
if [ -d "$DIR" ]; then
    echo "2 Success: $DIR exists and is a directory."
else
    echo "2 Error: $DIR does not exist."
fi
