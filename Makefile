mockgen:
	mockgen -source=./internal/app/domain/ports.go -destination=./internal/app/domain/tests/mocks/ports.go

tests: mockgen
	go test ./... -test.v

clean:
	rm -rf build