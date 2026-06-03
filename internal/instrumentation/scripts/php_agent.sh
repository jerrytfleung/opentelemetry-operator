#!/bin/sh
# Init container script for PHP auto-instrumentation.

set -e
set -x

instrumentation_src="$1"
echo "instrumentation_src: $instrumentation_src"
mounted_dir="$2"
echo "mounted_dir: $mounted_dir"

ls -l
