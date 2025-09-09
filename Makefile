main_package_path=./...
build_dir=./bin

init:
	@go get ${main_package_path}

run:
	@go run ./cmd/server/main.go

build:
	@make test
	@mkdir -p ${build_dir} && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${build_dir} ${main_package_path}

clean:
	@echo "Cleaning up..."
	@rm -rf ${build_dir}

test:
	CGO_ENABLED=0 go test ${main_package_path}

test-verbose:
	CGO_ENABLED=0 go test -v ${main_package_path}

