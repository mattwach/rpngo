#!/bin/bash
set -e
set -x
cd bin
base=$(pwd)
cd $base/minimal/mrpn && go build
cd $base/tinygo/serialonly && make build
cd $base/tinygo/ili9341 && make build

cd $base/tinygo/serialonly && make build TARGET=pico2
cd $base/tinygo/ili9341 && make build TARGET=pico2
cd $base/tinygo/picocalc && make build TARGET=pico2


