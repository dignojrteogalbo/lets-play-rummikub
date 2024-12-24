fmt:
	go fmt ./...

test: fmt
	go test -v -coverprofile cover.out ./internal/...
	go tool cover -html cover.out -o cover.html
	open cover.html

run:
	go run ./main.go