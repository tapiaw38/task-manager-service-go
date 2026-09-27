.DEFAULT_GOAL := help
.PHONY: help welcome install-deps run run-dev test test-cover cover cover-html gen-mocks mocks clean fmt vet lint run-doc docs

welcome:
	@printf "\033[1;33mTask Manager Service\033[0m\n\n"

install-deps: welcome
	@echo "Installing Delve debugger..."
	@go install github.com/go-delve/delve/cmd/dlv@latest
	@echo "Installing Air live reloader..."
	@go install github.com/air-verse/air@latest

run: welcome
	@go run ./cmd/api/

run-dev: welcome
	@air

test: welcome
	@go test ./... -race -coverpkg=./... -coverprofile=coverage.out -covermode=atomic
	@./build/covignore.sh coverage.out .covignore

test-cover: test
	@go tool cover -func=coverage.out | tail -1

cover: test
	@go tool cover -func=coverage.out | tail -1

cover-html: test
	@go tool cover -html=coverage.out

gen-mocks: welcome
	@mockgen -source=internal/adapters/datasources/repositories/task/repository.go -destination=internal/adapters/datasources/repositories/task/mocks/repository.go -package=mocks
	@for usecase in create list get complete delete; do \
		mockgen -source=internal/usecases/task/$$usecase.go -destination=internal/usecases/task/mocks/$$usecase.go -package=mocks; \
	done

mocks: gen-mocks

clean: welcome
	@go clean -cache
	@rm -f coverage.out build-errors.log
	@rm -rf build tmp

fmt:
	@gofmt -w ./cmd ./internal

vet:
	@go vet ./...

lint: fmt vet

docs:
	@npx -y docsify-cli serve ./docs/guide

run-doc: docs

help: welcome
	@printf '%s\n' 'clean cover cover-html docs fmt gen-mocks install-deps lint mocks run run-dev run-doc test test-cover vet'
