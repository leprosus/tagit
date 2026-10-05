.PHONY: lint test sha256 brew

BREW_FORMULA ?= ../homebrew-tap/Formula/tagit.rb

lint:
	golangci-lint run .

test:
	go test

sha256:
	@set -eu; \
		version="$$(git describe --tags --abbrev=0)"; \
		archive="$$(mktemp)"; \
		trap 'rm -f "$$archive"' EXIT; \
		trap 'exit 1' HUP INT TERM; \
		curl -fsSL -o "$$archive" "https://github.com/leprosus/tagit/archive/refs/tags/$${version}.tar.gz"; \
		checksum="$$(shasum -a 256 "$$archive")"; \
		printf '%s\n' "$${checksum%% *}"

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
