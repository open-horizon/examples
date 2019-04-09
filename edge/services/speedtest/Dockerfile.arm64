FROM arm64v8/python:3-alpine
WORKDIR /

# Install useful dev tools
RUN apk --no-cache --update add jq vim

# Install the SpeedTest CLI
RUN pip install speedtest-cli

# Install flask (for the REST API server)
RUN pip install Flask

# Copy over the source
COPY speedtest_server.py .

# Run the daemon
CMD python speedtest_server.py
