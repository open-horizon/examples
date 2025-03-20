ARG BUILD_IMAGE
ARG RUN_IMAGE

# Build stage
FROM $BUILD_IMAGE AS build

# Add source code and set working directory
ADD . /src
WORKDIR /src

# Set Go environment for Linux build
ENV GOOS=linux

# Set Go architecture
ARG GOARCH
ENV GOARCH=$GOARCH

# Install dependencies and build the application
RUN apk update && \
	apk --no-cache add build-base git mercurial gcc portaudio portaudio-dev pkgconfig && \
	go mod tidy && \
	go get -v && \
	go build -o /client

##################################################################################################

# Runtime stage
FROM $RUN_IMAGE

# Environment variables
ENV MQTT_BROKER=volantmq:1883 \
	MQTT_CLIENT_USER=fft-client \
	MQTT_CLIENT_PASS=client-pass \
	MQTT_CLIENT_CLIENT=fft-client \ 
	MQTT_QOS=2 \
	LOG_LEVEL=info \
	SAMPLE_RATE=48000 \
	RECORD_FRAME=3 \
	DEVICE_ID=-1

# Copy the built binary from the build stage and the start.sh script
COPY --from=build /client /client
ADD start.sh .

# Install runtime dependencies and make the start.sh script executable
RUN apk --no-cache add portaudio && \
	chmod +x start.sh

# Set the entrypoint
ENTRYPOINT ["./start.sh"]
