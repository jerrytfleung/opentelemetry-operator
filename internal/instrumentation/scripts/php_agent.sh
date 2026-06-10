#!/bin/sh
# Init container script for PHP auto-instrumentation.
# Runs in the opentelemetry-auto-instrumentation-php init container (the one with the
# instrumentation image), after the opentelemetry-auto-instrumentation-clone-php init
# container has staged the PHP API version, PHP extension directory, glibc / musl
# standard C library and thread safety into the shared volume.
set -e

# Inputs:
#   $1 - Instrumentation source directory containing subdirectories for each API version and standard C library variant, with the compiled agent extensions inside (e.g. /autoinstrumentation/20240924/glibc/non-zts).
#   $2 - Directory where the agent extensions should be copied to (e.g. /otel-auto-instrumentation-php).

instrumentation_src="$1"
mounted_dir="$2"

thread_safety=$(cat /otel-auto-instrumentation-php-clone/thread_safety.txt)
api=$(cat /otel-auto-instrumentation-php-clone/api.txt)
standard_c_lib=$(cat /otel-auto-instrumentation-php-clone/standard_c_lib.txt)

cp -rf "$instrumentation_src"/"$api"/"$standard_c_lib"/"$thread_safety"/* "$mounted_dir"/
cp -rf "$instrumentation_src"/* "$mounted_dir"/
