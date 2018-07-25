all: sdr-db

%: %.go
	go build $<

clean:
	rm -r sdr-db

.PHONY: all clean
