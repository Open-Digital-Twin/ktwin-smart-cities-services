FROM golang:alpine AS build

ARG SERVICE_NAME

WORKDIR /build

# First install library dependencies
# (These are expensive to download and change rarely;
# doing this once up front saves time on rebuilds)
COPY go.mod go.sum ./
RUN go mod download

# Copy the whole application tree in
COPY . .

# Build the specific component we want to run
RUN go build -o service ./cmd/${SERVICE_NAME}

# Final runtime image:
FROM alpine
# Get the built binary
COPY --from=build /build/service /usr/bin
# And set it as the main container command
CMD ["service"]