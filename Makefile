all: local

docker-image:
	docker build -t aanm/node-metrics:${VERSION} .

tests:
	go test -mod=vendor ./...

node-metrics: tests
	CGO_ENABLED=0 go build -mod=vendor -a -installsuffix cgo -o $@ ./main.go

local: node-metrics
	strip node-metrics

clean:
	rm -fr node-metrics
