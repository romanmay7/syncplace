# syntax=docker/dockerfile:1

FROM golang:1.23-bookworm as builder
#FROM  alpine:3.5

# Set destination for COPY
WORKDIR /syncplace

# Download Go modules
COPY go.mod go.sum ./
COPY *.go ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o syncplace-srv .

# Use the official Debian slim image for a lean production container.
# https://hub.docker.com/_/debian
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM debian:bookworm-slim

WORKDIR /syncplace  # Set the working directory in the final stage!

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*


#GET wait-for-it.sh from these repos
#wget https://github.com/vishnubob/wait-for-it/blob/master/wait-for-it.sh
#wget https://github.com/eficode/wait-for/blob/master/wait-for
#chmod +x wait-for-it.sh

# Copy the binary to the production image from the builder stage.
COPY --from=builder /syncplace/syncplace-srv /syncplace/syncplace-srv
COPY --from=builder /syncplace/wait-for-it.sh /syncplace/wait-for-it.sh
RUN chmod +x /syncplace/wait-for-it.sh

EXPOSE 3100

# Run
#CMD ["/syncplace/wait-for-it.sh", "postgres:5432", "--", "/syncplace/syncplace-srv"]
#CMD ["/syncplace/syncplace-srv"]

#Run and bind to local port :  "docker run -it -p 8080:8080 my_app" , "docker run -it -p 3010:3010 syncplace-srv"
#Build and Run :"docker build -t syncplace-srv3 . && docker run -it syncplace-srv3"