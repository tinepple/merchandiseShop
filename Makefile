test:
	go test ./...

test-storage:
	go test ./internal/storage/... -tags=integration

test-all:
	test test-storage

mocks:
	mockgen -source=internal/handler/deps.go -destination='internal/handler/mocks/mocks.go' -package=mocks
