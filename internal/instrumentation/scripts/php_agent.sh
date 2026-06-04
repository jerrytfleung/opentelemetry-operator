#!/bin/sh
# Init container script for PHP auto-instrumentation.

set -e
set -x

instrumentation_src="$1"
echo "instrumentation_src: $instrumentation_src"
mounted_dir="$2"
echo "mounted_dir: $mounted_dir"

ls -l

cp -rf "$instrumentation_src"/* "$mounted_dir"/

extension_dir=$(</otel-auto-instrumentation-php-clone/extension_dir.txt)
echo "$extension_dir"

api=$(</otel-auto-instrumentation-php-clone/api.txt)
echo "$api"

standard_c_lib=$(</otel-auto-instrumentation-php-clone/standard_c_lib.txt)
echo "$standard_c_lib"

cp -rf /autoinstrumentation/"$api"/"$standard_c_lib"/* "$extension_dir"/
