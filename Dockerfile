# build stage
FROM golang as builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o sd .

# final stage
FROM scratch
COPY --from=builder /app/sd /app/
EXPOSE 5000
ENTRYPOINT ["/app/sd"]
