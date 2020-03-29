############################
# STEP 1 build executable binary
############################
FROM golang:1.13 AS builder

# Install git.
# Git is required for fetching the dependencies.
# RUN apk update && apk add --no-cache git make build-base

WORKDIR /build

ENV GOOS=linux
ENV GOARCH=amd64

COPY . ./
RUN go mod download
RUN go install github.com/djimenez/iconv-go

# Build the binary.
RUN go build -o ./main

# EXPOSE 2222
# CMD ["/build/main"]

############################
# STEP 2 build a small image
############################
FROM scratch

# # Copy our static executable.
COPY --from=builder /build/main /
# # Run the binary.
EXPOSE 2222
CMD ["/main"]
