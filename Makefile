.PHONY: lint test sha256 brew

BREW_FORMULA ?= ../homebrew-tap/Formula/tagit.rb

lint:
	golangci-lint run .

test:
	go test

sha256:
	@version="$$(git describe --tags --abbrev=0)"; \
		curl -fsSL "https://github.com/leprosus/tagit/archive/refs/tags/$${version}.tar.gz" | shasum -a 256 | awk '{print $$1}'

brew:
	@set -eu; \
		version="$$(git describe --tags --abbrev=0)"; \
		checksum="$$( $(MAKE) sha256 )"; \
		printf '%s\n' "$$checksum" | grep -Eq '^[0-9a-f]{64}$$'; \
		test -f "$(BREW_FORMULA)"; \
		grep -q '^  url "' "$(BREW_FORMULA)"; \
		grep -q '^  sha256 "' "$(BREW_FORMULA)"; \
		sed -i.bak \
			-e "s|^  url \".*\"|  url \"https://github.com/leprosus/tagit/archive/refs/tags/$${version}.tar.gz\"|" \
			-e "s|^  sha256 \".*\"|  sha256 \"$${checksum}\"|" \
			"$(BREW_FORMULA)"; \
		rm -f "$(BREW_FORMULA).bak"
