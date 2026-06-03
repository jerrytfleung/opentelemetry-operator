#!/bin/sh
# Init container script for PHP auto-instrumentation.

set -e
set -x

echo "PHP auto-instrumentation init container started with args: $*"

mounted_dir="$1"

echo "Copying auto-instrumentation files to mounted directory: $mounted_dir"

cp -r /autoinstrumentation/opentelemetry.ini "$mounted_dir"
cp -r /autoinstrumentation/version.txt "$mounted_dir"
#
#extension_dir=$(php -i | grep "^extension_dir" | awk '{print $5}')
#
#api=$(php -i | grep "^PHP API => " | awk '{print $4}')
#
## check if alpine
#standard_c_lib=glibc
#if [ -f /etc/alpine-release ]; then
#    standard_c_lib=musl
#fi
#
#cp -r /autoinstrumentation/"$api"/"$standard_c_lib"/* "$extension_dir"

echo "PHP auto-instrumentation init container completed successfully"
