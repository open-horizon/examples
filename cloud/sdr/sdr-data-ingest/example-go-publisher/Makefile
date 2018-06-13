all: msghub-pubsync msghub-consumer

%: %.go
	go build $<

clean:
	rm -r msghub-pubsync msghub-consumer

.PHONY: all clean