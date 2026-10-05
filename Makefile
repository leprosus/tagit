.PHONY: lint test sha256 brew

BREW_FORMULA ?= ../homebrew-tap/Formula/tagit.rb
define BREW_TEMPLATE
class Tagit < Formula
  desc "To create and increment semantic Git tags"
  homepage "https://github.com/leprosus/tagit"
  url "https://github.com/leprosus/tagit/archive/refs/tags/@VERSION@.tar.gz"
  sha256 "@SHA256@"
  license "MIT"

  depends_on "go" => :build
  depends_on "git"

  def install
    system "go", "build",
      *std_go_args(ldflags: "-s -w -X github.com/leprosus/tagit/command.releaseVersion=v#{version}"), "."
  end

  test do
    assert_equal "v#{version}\n", shell_output("#{bin}/tagit ver")
    system "git", "init"
    assert_match "unknown version increment",
      shell_output("#{bin}/tagit invalid 2>&1", 1)
  end
end
endef
export BREW_TEMPLATE

lint:
	golangci-lint run ./...

test:
	go test --race ./...

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
		printf '%s\n' "$$version" | grep -Eq '^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$$'; \
		test -f "$(BREW_FORMULA)"; \
		checksum="$$( $(MAKE) sha256 )"; \
		printf '%s\n' "$$checksum" | grep -Eq '^[0-9a-f]{64}$$'; \
		formula_tmp="$$(mktemp "$(BREW_FORMULA).XXXXXX")"; \
		trap 'rm -f "$$formula_tmp"' EXIT; \
		trap 'exit 1' HUP INT TERM; \
		cp -p "$(BREW_FORMULA)" "$$formula_tmp"; \
		printf '%s\n' "$$BREW_TEMPLATE" | \
			sed -e "s|@VERSION@|$$version|g" -e "s|@SHA256@|$$checksum|g" > "$$formula_tmp"; \
		mv -f "$$formula_tmp" "$(BREW_FORMULA)"
