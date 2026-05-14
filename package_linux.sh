#!/bin/bash

set -e
set -x

# Create a simple .tar.gz release of linux packages
cd -- "$(dirname -- "${BASH_SOURCE[0]}")"

make clean
make all

PACKAGE_BASE=rpngo_linux_and_pico_$(date +%Y%m%d)
TMPDIR=/tmp/$PACKAGE_BASE

rm -rf $TMPDIR

#for path in $(find bin -type f -executable); do
#  mkdir -p $TMPDIR/$(dirname $path)
#  cp $path $TMPDIR/$path
#done

PICO_DIR=$TMPDIR/pico_rp2040_uf2
mkdir -p $PICO_DIR
for path in $(find bin -name '*_pico_*.uf2'); do
  cp $path $PICO_DIR
done

PICO_DIR=$TMPDIR/pico2_rp2350_uf2
mkdir -p $PICO_DIR
for path in $(find bin -name '*_pico2_*.uf2'); do
  cp $path $PICO_DIR
done

cp -r examples $TMPDIR
cp -r img $TMPDIR
cp *.md $TMPDIR

cd /tmp
tar cvfz $PACKAGE_BASE.tar.gz $PACKAGE_BASE


