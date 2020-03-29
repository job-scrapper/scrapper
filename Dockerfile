############################
# STEP 1 build executable binary
############################
FROM golang:1.13

WORKDIR /build

ENV GOOS=linux
ENV GOARCH=amd64

COPY . ./
RUN go mod download
RUN go install github.com/djimenez/iconv-go

# Build the binary.
RUN go build -o ./main

EXPOSE 2222
CMD ["/build/main"]
