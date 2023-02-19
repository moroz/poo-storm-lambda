GOOS := "linux"
GOARCH := "amd64"
CGO_ENABLED := "0"
FUNCTION_NAME := "poo-storm-comment-api"

dynamodb: FORCE
	docker run -d -p 8000:8000 amazon/dynamodb-local

build:
	go build -o main .

package: build
	zip main.zip main

deploy: package
	aws lambda update-function-code --function-name $(FUNCTION_NAME) --zip-file="fileb://main.zip"

clean:
	rm main main.zip

FORCE: ;
