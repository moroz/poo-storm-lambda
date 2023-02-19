GOOS := "linux"
GOARCH := "amd64"
CGO_ENABLED := "0"

dynamodb: FORCE
	docker run -d -p 8000:8000 amazon/dynamodb-local

build:
	go build -o main .

package: build
	zip main.zip main

FORCE: ;
