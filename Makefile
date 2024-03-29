docker_run_with_volume = docker run --rm -v "$$(pwd)":/usr/src/pmongo -w /usr/src/pmongo golang:1.21

# go fmt
fmt:
	$(docker_run_with_volume) gofmt -s -w .

# go fmt list files affected
fmt_list:
	$(docker_run_with_volume) gofmt -s -l .

# Run tests w/ docker services
test_ci:
	sh test/entrypoint.sh

test: clean
	docker-compose -f $(compose_file) up --exit-code-from go-test

clean:
	docker rm mongodb -f || true
	docker rm go-test -f || true

compose_file  = test/docker-compose.yml