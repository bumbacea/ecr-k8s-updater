module github.com/alexbumbacea/ecr-k8s-updater

go 1.16

require (
	github.com/aws/aws-sdk-go-v2/config v1.5.0
	github.com/aws/aws-sdk-go-v2/service/ecr v1.4.1
	github.com/docker/docker v25.0.6+incompatible
	github.com/joho/godotenv v1.3.0
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	golang.org/x/net v0.36.0 // indirect
	gotest.tools/v3 v3.0.3 // indirect
	k8s.io/api v0.21.3
	k8s.io/apimachinery v0.21.3
	k8s.io/client-go v0.21.3
)
