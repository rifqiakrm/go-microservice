FROM golang:1.15.7

# Install dockerize for wait capabilities
RUN apt-get update && apt-get install -y wget
ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz

# Set working directory
WORKDIR /app

# Copy the source
COPY . .

# Clean modcache
RUN go clean -modcache

RUN git config --global url."https://{username}:{access_token}@github.com/{username}".insteadOf "https://github.com/{username}"

# Build the app
RUN go build

# Execute the app
CMD ["./go-microservice"]
