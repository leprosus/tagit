lint:
	golangci-lint run .

test:
	go test

sha256:
	@version="$$(git describe --tags --abbrev=0)"; \
		curl -fsSL "https://github.com/leprosus/tagit/archive/refs/tags/$${version}.tar.gz" | shasum -a 256 | awk '{print $$1}'
