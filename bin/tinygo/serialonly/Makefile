TARGET ?= pico
DATE = $(shell date +"%Y%m%d")
RELEASE = rpngo_serialonly_$(TARGET)_$(DATE).uf2

build:
	tinygo build -target=$(TARGET) -scheduler=tasks -o $(RELEASE)

flash:
	tinygo flash -target=$(TARGET) -scheduler=tasks
	

allocs:
	tinygo build -target=$(TARGET) -print-allocs='.*' 2>&1

