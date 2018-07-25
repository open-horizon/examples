all: msghub-producer msghub-consumer

%: %.go
	go build $<

clean:
	rm -r msghub-producer msghub-consumer

.PHONY: all clean