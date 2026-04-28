#!/bin/bash
set -e
set -x
cd bin
base=$(pwd)
cd $base/ncurses/nrpn && go build
cd $base/fyne/frpn && go build
cd $base/tinygo/picocalc && make build


