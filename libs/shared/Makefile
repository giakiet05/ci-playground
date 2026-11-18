# CI contract (github.com/giakiet05/ci-templates): `lint` and `test`.

.PHONY: lint test

## lint: fail on unformatted files, then go vet
lint:
	@unformatted=$$(gofmt -l .); if [ -n "$$unformatted" ]; then echo "gofmt needed: $$unformatted"; exit 1; fi
	go vet ./...

## test: unit tests with the race detector
test:
	go test -race -count=1 ./...
