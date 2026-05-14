default: test buildbase

all: test buildall

clean:
	rm $(shell find bin -type f -executable)
	rm $(shell find bin -name '*.uf2')

test:
	make -C common
	make -C drivers/tinygo

buildbase:
	./buildbase.sh

buildall: buildbase
	./buildextras.sh

  
