FROM golang:1.24 as builder
WORKDIR /go/src/github.com/alexbumbacea/ecr-k8s-updater
COPY go.* ./
RUN go mod download
COPY main.go ./
RUN  CGO_ENABLED=0 GOOS=linux go build -o /opt/app .
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /opt
COPY --from=builder /opt/app ./
CMD ["./app"]
