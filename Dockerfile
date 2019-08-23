# build stage
FROM golang as builder
ENV GO111MODULE=on
WORKDIR /app
RUN go mod download
COPY . .
# https://golang.org/cmd/link/ (space-separated flags to pass to the external linker)
# https://www.elwinar.com/articles/statically-link-golang-binaries
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o sd .
# final stage
FROM alpine:latest AS runtime
COPY --from=builder /app/sd /app/
EXPOSE 5000
ENTRYPOINT ["/app/sd"]

