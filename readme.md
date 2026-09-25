# tagit

`tagit` creates and increments semantic git tags in the current repository.

## Installation

Install with Homebrew:

```sh
brew tap leprosus/tap
brew trust leprosus/tap
brew install tagit
```

`brew trust` explicitly permits Homebrew to execute formulas from the
third-party `leprosus/tap` repository.

## Usage

Build the command:

```sh
go build -o tagit .
```

Run it from a Git repository:

```sh
tagit          # no semantic tags yet: creates v0.0.1; otherwise increments patch
tagit patch    # v0.0.1 -> v0.0.2
tagit minor    # v0.0.1 -> v0.1.0
tagit major    # v0.0.1 -> v1.0.1
tagit list     # lists valid SemVer tags from oldest to newest
tagit list 10  # lists the latest 10 valid SemVer tags
tagit set v1.2.3 # creates the v1.2.3 tag
tagit del v1.2.3 # deletes the v1.2.3 tag when it exists
```

Outside a Git repository, or with an invalid `set` version, `tagit` prints
usage help and exits successfully. An invalid `del` version and every other
error write the error and usage help to standard error, then exit with status
`1`. Deleting a tag that does not exist completes without changes.

The command considers only stable tags in the exact `v{MAJOR}.{MINOR}.{PATCH}` form,
chooses the greatest version, creates the next tag locally, and prints it. It
requires the Git command-line client and does not push tags to a remote.
