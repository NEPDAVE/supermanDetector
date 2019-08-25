# build stage
#FROM frolvlad/alpine-glibc
#FROM golang as builder

#FROM alpine-pkg-glibc
#FROM golang:alpine as builder
#FROM blang/golang-alpine as builder


FROM ubuntu:18.04 as builder

# Update and upgrade repo
RUN apt-get update -y -q && apt-get upgrade -y -q

# Install tools we might need
RUN DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends -y -q curl build-essential ca-certificates git

# Download Go 1.2.2 and install it to /usr/local/go
RUN curl -s https://storage.googleapis.com/golang/go1.12.2.linux-amd64.tar.gz| tar -v -C /usr/local -xz

# Let's people find our Go binaries
ENV PATH $PATH:/usr/local/go/bin

ENV GO111MODULE=on

WORKDIR /app/

COPY go.mod .
COPY go.sum .
COPY GeoLite2-City.mmdb .

RUN go mod download

COPY . .

#RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o sd .
#RUN CGO_ENABLED=0 GOOS=linux go build -o sd .
#RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o sd .
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o sd .
#RUN CGO_ENABLED=1 go build -a -installsuffix cgo -o sd .

# final stage
#FROM ubuntu:18.04
FROM scratch
WORKDIR /app/
COPY --from=builder /app/sd /app/GeoLite2-City.mmdb /app/
#COPY GeoLite2-City.mmdb .
EXPOSE 5000
ENTRYPOINT ["/app/sd"]

