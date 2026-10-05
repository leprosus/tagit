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

Local builds report `dev`. To embed a release version, build with:

```sh
go build -ldflags "-X github.com/leprosus/tagit/command.releaseVersion=v1.2.3" -o tagit .
```

Homebrew embeds the formula's release version when building the executable.
Run `tagit ver` from any directory to print the running executable's version;
this command does not require Git.

Run it from a Git repository:

```sh
tagit          # no semantic tags yet: creates v0.0.1; otherwise increments patch
tagit ver      # prints the executable's release version, or dev for local builds
tagit patch    # v0.0.1 -> v0.0.2
tagit minor    # v0.0.1 -> v0.1.0
tagit major    # v0.0.1 -> v1.0.1
tagit list     # lists valid SemVer tags from oldest to newest
tagit list 10  # lists the latest 10 valid SemVer tags
tagit set v1.2.3 # creates the v1.2.3 tag
tagit del v1.2.3 # deletes the v1.2.3 tag when it exists
tagit del-all v1.2.3 # deletes the tag from all remotes, then locally
```

Outside a Git repository, or with an invalid `set` version, `tagit` prints
usage help and exits successfully. An invalid `del` version and every other
error write the error and usage help to standard error, then exit with status
`1`. Deleting a tag that does not exist completes without changes.

Errors while checking the repository, including a missing Git executable or
repository access failures, also exit with status `1` and write to standard error.

The command considers only stable tags in the exact `v{MAJOR}.{MINOR}.{PATCH}` form,
chooses the greatest version, creates the next tag locally, and prints it. It
requires the Git command-line client. Tag creation is local and does not push tags
to a remote; `del-all` performs remote deletion.

`del-all` checks and deletes the exact tag on every configured remote's push URL
(including multiple push URLs). It reports each successfully processed remote
on standard output. Tags that are already absent count as successfully processed.
After all remotes succeed, it deletes the local tag and reports the local result.
With no remotes configured, it only deletes the local tag.

If a remote fails, `del-all` continues with the remaining remotes, reports the
failures on standard error, and exits with status `1`. The local tag is retained;
successful remote deletions are not rolled back. Run the command again after
resolving the failure to complete deletion. If a remote has separate fetch and
push URLs, only its push destinations are modified.
