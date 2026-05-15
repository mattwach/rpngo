# Changelog

## 1.0.1 - Unreleased

Bug Fixes

- Fixed hang and improved error message when the `.rpngo` startup file executes
  with an error.

Features

- Added `split` function for splitting a string
- Added more examples to `examples/` directory

## 1.0.0 - `rpngo_linux_and_pico_20260515`

This is the initial version released with the new release process.  It contains
the following main files:

- `bin/`: Linux binaries
  - `frpn-x86_64.AppImage`: Graphical UI calculator (uses Fyne)
  - `mprn`: minimal calculator
  - `nrpn-x86_64.AppImage`: ncurses-based calculator
- `pico_rp2040_uf2/`: firmware for a Pi Pico microcontroller
  - `rpngo_ili9341_pico_20260515.uf2`: Support for a pico + ili9341 LCD
  - `rpngo_picocalc_pico_20260515.uf2`: PicoCalc firmware
  - `rpngo_serialonly_pico_20260515.uf2`: Minimal USB serial communications.
    Requires only the Pico.
- `pico_rp2350_uf2/`: All the same firmware files, but for the Pico 2
- `examples/`: Example programs for the calculator
- `README.md`, `USER_GUIDE.md`, `CHANGELOG.md`: Documentation


