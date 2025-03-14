ARG BUILD_IMAGE
ARG RUN_IMAGE

# Build stage
FROM $BUILD_IMAGE AS build

# Add source code and set working directory
ADD fft_client/ /src/fft_client
ADD fft_server/ /src/fft_server

WORKDIR /src/fft_server

# Set Go environment for Linux build
ENV GOOS=linux

# Set Go architecture
ARG GOARCH
ENV GOARCH=$GOARCH

# Install dependencies and build the application
RUN apk --no-cache add git && \
	go get -d -v && \
	go build -o /server

##################################################################################################

# Runtime stage
FROM $RUN_IMAGE

# Environment variables
ENV MQTT_BROKER=volantmq:1883 \
	MQTT_SERVER_USER=fft-server \
	MQTT_SERVER_PASS=server-pass \
	MQTT_SERVER_CLIENT=fft-server \
	MQTT_RESULTS_TOPIC=results \
	MQTT_QOS=2 \
	SAMPLE_RATE=48000 \
	LOG_LEVEL=info \
	NFFT=8192 \
	PEAKS_LIMIT=2 \
	PEAKS_THRESHOLD=0.25 \
	FREQS_THRESHOLD=0.5

# Copy the built binary from the build stage and the start.sh script
COPY --from=build /server /server
ADD ./fft_server/start.sh .

# Install runtime dependencies and make the start.sh script executable
RUN chmod +x start.sh

# Set the entrypoint
ENTRYPOINT ["./start.sh"]