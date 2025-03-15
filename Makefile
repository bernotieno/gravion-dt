.PHONY: all build clean run install

all: build

build:
	@go build

clean:
	rm -f dt

run: build
	./dt

install: build
	sudo mv dt /usr/local/bin/
