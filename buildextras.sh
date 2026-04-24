#!/bin/bash
set -e
set -x
cd bin
base=$(pwd)
cd $base/minimal/mrpn && go build
cd $base/tinygo/serialonly && tinygo build -target pico
cd $base/tinygo/ili9341 && make build


