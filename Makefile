race := 1

.PHONY: test
test: lint
	mkdir -p target
	go test -coverprofile target/coverage.out -race=$(race) -count=10 ./... -tags=$(tags) || exit 1; \
	go tool cover -html=target/coverage.out -o target/coverage.html

.PHONY: lint
lint:
	golint -set_exit_status ./...
	go vet ./...