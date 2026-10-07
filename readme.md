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

Run it from a Git repository:

```sh
tagit                 # no semantic tags yet: creates v0.0.1; otherwise increments patch
tagit ver             # prints the executable's release version, or dev for local builds
tagit history         # shows detailed history since the greatest local version tag
tagit history v1.2.3  # shows detailed history since the specified local tag
tagit changes         # lists commits since the greatest local version tag
tagit changes v1.2.3  # lists commits since the specified local tag
tagit latest          # prints the greatest local version tag
tagit show v1.2.3     # shows the local tag's commit, commit date, and tag message
tagit --help          # shows help successfully, from any directory
tagit help list       # shows help for a command; tagit list --help also works
tagit patch           # v0.0.1 -> v0.0.2
tagit minor           # v0.0.1 -> v0.1.0
tagit major           # v0.0.1 -> v1.0.1
tagit list            # lists valid SemVer tags from oldest to newest
tagit list 10         # lists the latest 10 valid SemVer tags
tagit set v1.2.3      # creates the v1.2.3 tag
tagit del v1.2.3      # deletes the v1.2.3 tag when it exists
tagit del-all v1.2.3  # deletes the tag from all remotes, then locally
```

Commands use positional arguments. Ctrl+C cancels active Git subprocesses.

`tagit show v1.2.3` prints the full commit hash, commit author (name and email), committer date in
ISO 8601 format, and the annotated tag message. Lightweight tags display
`Message: (none)`. The date belongs to the commit, not the tag creation time.
The command reads only local tags and does not modify the repository. A missing
tag or invalid argument produces an error with exit status `1` and empty stdout.

`tagit latest` compares local tags numerically in the exact `vMAJOR.MINOR.PATCH`
format, including `v0.0.0`, regardless of tag creation dates or the current branch.
It does not fetch remote tags or modify the repository. If no valid version tags
exist, it leaves stdout empty, reports an error on stderr, and exits with status `1`.

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

`tagit changes [vMAJOR.MINOR.PATCH]` lists commits reachable from `HEAD` but
not from the selected local tag, excluding the tagged commit itself. Without an
argument it selects the greatest local SemVer tag, as `tagit latest` does; if
none exist it lists the entire history reachable from `HEAD`. The selected tag
may be on another branch: the same reachability rule applies.

Each commit uses the same block format as `history` and `show`: full commit
hash, author, committer date, and the first line of the message on the same
line as `Message:`, newest first (Git date order). An empty range produces no output and succeeds. Invalid
arguments, a missing explicit tag, or a Git error produce exit status 1.
The command only reads the local repository.

`tagit history [vMAJOR.MINOR.PATCH]` selects the same range and uses the same output as `changes`. It
prints a separate block for each commit: full hash, author name and email,
committer date in ISO 8601, and the first line of the commit message. Blocks are separated
by a blank line. Lightweight and annotated tags select the same commit range;
tag annotations are not commit messages. Empty ranges succeed without output.

Example output:

```text
Commit: 9f8e7d6c5b4a32100123456789abcdef01234567
Author: Ivan Petrov <ivan@example.com>
Date: 2026-10-07T15:42:10+02:00
Message: Add changes command

Commit: 123456789abcdef0123456789abcdef012345678
Author: Ivan Petrov <ivan@example.com>
Date: 2026-10-07T14:18:03+02:00
Message: Fix tag validation
```

`show` also displays the commit author; its `Message` remains the annotated
tag's message, or `(none)` for a lightweight tag.

All three commands use the same field order: `Commit`, `Author`, `Date`,
`Message`. Each message is displayed on one line; subsequent lines are omitted.
For `show`, this is the first line of the annotated tag message, or `(none)`
for a lightweight tag.
