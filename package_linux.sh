#!/bin/bash

set -e
set -x

# Create a simple .tar.gz release of linux packages
cd -- "$(dirname -- "${BASH_SOURCE[0]}")"

make all

PACKAGE_BASE=rpngo_linux_and_pico_$(date +%Y%m%d)
TMPDIR=/tmp/$PACKAGE_BASE

rm -rf $TMPDIR

for path in $(find bin -type f -executable); do
  mkdir -p $TMPDIR/$(dirname $path)
  cp $path $TMPDIR/$path
done

for path in $(find bin -name '*.uf2'); do
  mkdir -p $TMPDIR/$(dirname $path)
  cp $path $TMPDIR/$path
done

cp -r examples $TMPDIR
cp -r img $TMPDIR
cp *.md $TMPDIR

cd /tmp
tar cvfz $PACKAGE_BASE.tar.gz $PACKAGE_BASE


