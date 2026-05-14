FS ?= littlefs
BLOCKDEV ?= machineflash
TARGET ?= pico
DATE = $(shell date +"%Y%m%d")
RELEASE = rpngo_ili9341_$(TARGET)_$(DATE).uf2

build:
	tinygo build -tags=$(FS),$(BLOCKDEV) -target=$(TARGET) -scheduler=tasks -o $(RELEASE)

flash:
	tinygo flash -tags=$(FS),$(BLOCKDEV) -target=$(TARGET) -scheduler=tasks
	

allocs:
	tinygo build -tags=$(FS),$(BLOCKDEV) -target=$(TARGET) -print-allocs='.*' 2>&1
