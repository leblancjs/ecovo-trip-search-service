FROM golang:alpine AS build

ARG TARGETARCH

ENV PROJECT_NAME=azure.com/ecovo/trip-search-service
ENV BINARY_NAME=trip-search-service

# Install dependencies for the build along with trusted certificates
RUN apk --no-cache add git
RUN apk --no-cache add ca-certificates

# Force the project to use Go mod for dependencies
ENV GO111MODULE=on

# Copy the project files
COPY . $GOPATH/src/${PROJECT_NAME}
WORKDIR $GOPATH/src/${PROJECT_NAME}/cmd

# Manage dependencies
RUN go mod download

# Build the project
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -ldflags '-w -s' -o /bin/${BINARY_NAME}

# Test the project
RUN CGO_ENABLED=0 go test ../...

# Expose port
EXPOSE 8080/tcp

# Start a new container from scratch to keep only the compiled binary
FROM scratch AS run

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /bin/${BINARY_NAME} /bin/${BINARY_NAME}

CMD ["/bin/trip-search-service"]
ENTRYPOINT ["/bin/trip-search-service"]