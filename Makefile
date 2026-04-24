default: test buildbase

all: test buildall

test:
	make -C common
	make -C drivers/tinygo

buildbase:
	./buildbase.sh

buildall: buildbase
	./buildextras.sh

  
