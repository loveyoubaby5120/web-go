all:
	go build ./... && \
	go run cmd/main/main.go
.PHONY: all
