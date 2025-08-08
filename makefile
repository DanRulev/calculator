run:
	go run ./cmd

test:
	go test ./... -v -cover

.PHONY: run test