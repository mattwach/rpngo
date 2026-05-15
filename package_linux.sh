#!/bin/bash

set -e
set -x

create_app_run() {
  local exec_path=$1
  local app_dir=$2
  cat <<EOF > $app_dir/AppRun
#!/bin/sh
SELF=\$(dirname "\$(readlink -f "\$0")")
export PATH="\$SELF/usr/bin:\$PATH"
export LD_LIBRARY_PATH="\$SELF/usr/lib:\$LD_LIBRARY_PATH"
exec "\$SELF/usr/bin/$(basename $exec_path)" "\$@"
EOF
  chmod +x $app_dir/AppRun
}

copy_icon_file() {
  local exec_path=$1
  local app_dir=$2
  local bin_name=$(basename $exec_path)
  local icon_dir=$app_dir/usr/share/icons/hicolor/256x256/apps
  mkdir -p $icon_dir
  cp img/rpngo_icon.png $icon_dir/$bin_name.png
  cp img/rpngo_icon.png $app_dir/$bin_name.png
}

create_desktop_file() {
  local exec_path=$1
  local app_dir=$2
  local bin_name=$(basename $exec_path)
cat <<EOF > $app_dir/rpngo.desktop
[Desktop Entry]
Name=$bin_name
Exec=$bin_name
Icon=$bin_name
Type=Application
Categories=Utility;
EOF
}

make_app_image() {
  local base_dir=$1
  local exec_path=$2
  local dir=/tmp/appimage
  rm -rf $dir
  local image_name=$(basename $exec_path)-x86_64.AppDir
  local app_dir=$dir/$image_name
  local bin_dir=$app_dir/usr/bin
  mkdir -p $bin_dir 
  cp $exec_path $bin_dir
  create_app_run $exec_path $app_dir
  copy_icon_file $exec_path $app_dir
  create_desktop_file $exec_path $app_dir
  pushd $dir
  linuxdeploy \
      --appdir $image_name \
      --output appimage \
      --desktop-file $image_name/rpngo.desktop
  cp $(basename $exec_path)-x86_64.AppImage $base_dir
  pushd
}

make_linux_apps() {
  local base_dir=$1
  mkdir -p $base_dir

  # mrpn is a special case that needs no dependencies
  mv bin/minimal/mrpn/mrpn $base_dir

  for path in $(find bin -type f -executable); do
    make_app_image $base_dir $path
  done
}

make_pico_uf2() {
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
}

cd -- "$(dirname -- "${BASH_SOURCE[0]}")"

make clean
make all

PACKAGE_BASE=rpngo_linux_and_pico_$(date +%Y%m%d)
TMPDIR=/tmp/$PACKAGE_BASE

rm -rf $TMPDIR
mkdir -p $TMPDIR

cp -r examples $TMPDIR
cp -r img $TMPDIR
cp *.md $TMPDIR

make_pico_uf2
make_linux_apps $TMPDIR/bin

cd /tmp
tar cvfz $PACKAGE_BASE.tar.gz $PACKAGE_BASE


