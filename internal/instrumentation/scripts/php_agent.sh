#!/bin/sh
# Init container script for PHP auto-instrumentation.

set -e

instrumentation_src="$1"
echo "Agent instrumentation_src: $instrumentation_src"
mounted_dir="$2"
echo "Agent mounted_dir: $mounted_dir"

ls -l

cp -rf "$instrumentation_src"/* "$mounted_dir"/

extension_dir=$(cat /otel-auto-instrumentation-php-clone/extension_dir.txt)
echo "Agent extension_dir: $extension_dir"

api=$(cat /otel-auto-instrumentation-php-clone/api.txt)
echo "Agent api: $api"

standard_c_lib=$(cat /otel-auto-instrumentation-php-clone/standard_c_lib.txt)
echo "Agent standard_c_lib: $standard_c_lib"

cp -rf /autoinstrumentation/"$api"/"$standard_c_lib"/* "$extension_dir"/
