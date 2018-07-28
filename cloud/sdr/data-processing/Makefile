all: wat-sample

wat-sample: demo-processing/main.go watson/*/*.go wutil/wutil.go
	go build -o $@ $<

clean:
	rm -r wat-sample

.PHONY: all clean
