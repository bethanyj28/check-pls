build:
	docker build -t bethanyj28/check-pls .
run:
	docker run --rm -p 8080:8080 bethanyj28/check-pls
test:
	go test ./...
vendor:
	go mod vendor && go mod tidy
