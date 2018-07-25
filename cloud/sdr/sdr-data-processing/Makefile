all: stt-sample

stt-sample: demo-processing/main.go watson/stt/stt.go
	go build -o $@ $<

clean:
	rm -r stt-sample

.PHONY: all clean
