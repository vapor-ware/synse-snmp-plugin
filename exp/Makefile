cover: test
	go tool cover -html=coverage.out

test:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...