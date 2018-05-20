.PHONY: test deps

test: deps
	go test ./pkg/...

deps: vendor
	go get github.com/hashicorp/vault/api
	go get github.com/hashicorp/consul/api
	go get github.com/aws/aws-sdk-go/aws
	go get github.com/aws/aws-sdk-go/service

vendor: Gopkg.lock
	dep ensure -vendor-only
