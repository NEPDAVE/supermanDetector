# build stage
# using ubuntu to ensure base images uses Glibc not Musl to help avoid dependency runtime errors
FROM ubuntu:18.04 as builder

# update and upgrade ubuntu
RUN apt-get update -y -q && apt-get upgrade -y -q

# install curl
RUN DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends -y -q curl build-essential ca-certificates git

# download go 1.12
RUN curl -s https://storage.googleapis.com/golang/go1.12.2.linux-amd64.tar.gz| tar -v -C /usr/local -xz

# create GOPATH
ENV PATH $PATH:/usr/local/go/bin

# enable go mod
ENV GO111MODULE=on

# build image
WORKDIR /app/

COPY go.mod .
COPY go.sum .
COPY GeoLite2-City.mmdb .

RUN go mod download

COPY . .

#using external linker for build to account for sqlite driver and GeoLite2 dependencies
#https://golang.org/cmd/link/ (space-separated flags to pass to the external linker)
#https://www.elwinar.com/articles/statically-link-golang-binaries
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o sd .

# final stage
FROM scratch
WORKDIR /app/
COPY --from=builder /app/sd /app/GeoLite2-City.mmdb /app/
EXPOSE 5000
ENTRYPOINT ["/app/sd"]
