SOURCES := $(wildcard *.go)

bin/sprout: $(SOURCES)
	@ go build -o bin/sprout .

clean:
	@ rm -rf bin
