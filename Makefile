.PHONY: all build clean run

all: build

build:
	go build

clean:
	rm -f dt

run: build
	./dt
