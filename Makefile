main_package_path=./...
build_dir=./bin

init:
	@go get ${main_package_path}

run:
	@go run ./cmd/server/main.go

build: test
	@mkdir -p ${build_dir} && \
		GOOS=linux \
		GOARCH=amd64 \
		CGO_ENABLED=0 \
		go build -o ${build_dir} ${main_package_path}

build-run: build
	@${build_dir}/server

build-ci: test-verbose
	@mkdir -p ${build_dir} && \
		GOOS=linux \
		GOARCH=amd64 \
		CGO_ENABLED=0 \
		go build -o ${build_dir} ${main_package_path}

clean:
	@echo "Cleaning up..."
	@rm -rf ${build_dir}
	@go clean -testcache -cache

test:
	@CGO_ENABLED=0 go test ${main_package_path}

test-verbose:
	@CGO_ENABLED=0 go test -v -coverprofile=coverage.out ${main_package_path}
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out

