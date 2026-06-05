#!/bin/sh
# Init container script for PHP auto-instrumentation.

set -e

instrumentation_src="$1"
echo "Agent instrumentation_src: $instrumentation_src"
mounted_dir="$2"
echo "Agent mounted_dir: $mounted_dir"

DIR="/autoinstrumentation/20240924"
# Always wrap variables in double quotes to handle paths with spaces safely
if [ -d "$DIR" ]; then
    echo "3 Success: $DIR exists and is a directory."
else
    echo "3 Error: $DIR does not exist."
fi

INST_DIR="/otel-auto-instrumentation-php"
# Always wrap variables in double quotes to handle paths with spaces safely
if [ -d "$INST_DIR" ]; then
    echo "3 Success: $INST_DIR exists and is a directory."
else
    echo "3 Error: $INST_DIR does not exist."
fi

CLONE_DIR="/otel-auto-instrumentation-php-clone"
# Always wrap variables in double quotes to handle paths with spaces safely
if [ -d "$CLONE_DIR" ]; then
    echo "3 Success: $CLONE_DIR exists and is a directory."
else
    echo "3 Error: $CLONE_DIR does not exist."
fi

thread_safety=$(cat /otel-auto-instrumentation-php-clone/thread_safety.txt)
echo "Agent thread_safety: $thread_safety"

api=$(cat /otel-auto-instrumentation-php-clone/api.txt)
echo "Agent api: $api"

standard_c_lib=$(cat /otel-auto-instrumentation-php-clone/standard_c_lib.txt)
echo "Agent standard_c_lib: $standard_c_lib"

cp -rf "$instrumentation_src"/"$api"/"$standard_c_lib"/"$thread_safety"/* "$mounted_dir"/
cp -rf "$instrumentation_src"/opentelemetry.ini "$mounted_dir"/
