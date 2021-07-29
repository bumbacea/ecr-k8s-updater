FROM golang:1.16
WORKDIR /go/src/github.com/alexbumbacea/ecr-k8s-updater
COPY go.* ./
RUN go mod download
COPY main.go ./
RUN  CGO_ENABLED=0 GOOS=linux go build -o app .