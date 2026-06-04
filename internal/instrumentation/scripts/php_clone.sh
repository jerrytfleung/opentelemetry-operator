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

extension_dir=$(php -i | grep "^extension_dir" | awk '{print $5}')
echo "$extension_dir" > /otel-auto-instrumentation-php-clone/extension_dir.txt

api=$(php -i | grep "^PHP API => " | awk '{print $4}')
echo "$api" > /otel-auto-instrumentation-php-clone/api.txt

# check if alpine
standard_c_lib=glibc
if [ -f /etc/alpine-release ]; then
    standard_c_lib=musl
fi
echo "$standard_c_lib" > /otel-auto-instrumentation-php-clone/standard_c_lib.txt

e=$(cat /otel-auto-instrumentation-php-clone/extension_dir.txt)
echo "Read extension_dir: $e"

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
