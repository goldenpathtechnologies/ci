# Contributing Guidelines
Thank you for showing interest in contributing to `ci`. Before contributing, please
make yourself familiar with the [Code of Conduct](CODE_OF_CONDUCT.md). The Code of
Conduct governs all aspects of this project, including but not limited to variable
names, names of branches, comments, personal interactions, commit messages, and
file names. I'm relatively new to Go (I started using Go not long before
[this commit](https://github.com/goldenpathtechnologies/ci/commit/40a24213c12cb9d2cd9ce7b4aba37db6776ed75d)),
so if you find anything here that is inaccurate or missing, please feel free to open
a corresponding issue or discussion.

## Table of contents

- [Creating issues](#creating-issues)
- [Creating pull requests](#creating-pull-requests)
- [Software development process](#software-development-process)
  - [Repository information](#repository-information)
    - [Branches](#branches)
    - [Changelog](#changelog)
  - [Conventional commits](#conventional-commits)
    - [Breaking changes](#breaking-changes)
  - [Coding standards and style](#coding-standards-and-style)
  - [Testing and coverage](#testing-and-coverage)
    - [Writing tests](#writing-tests)
    - [Running tests](#running-tests)
  - [Building the project](#building-the-project)
  - [Development installation](#development-installation)

## Creating issues
There isn't a lot of ceremony around creating issues for `ci`. It's enough to follow
these recommendations:

- Be clear and concise. Don't write a novel when a few sentences will suffice.
- Do **not** submit security vulnerabilities to the issue tracker. Email me at [daryl@goldenpath.ca](mailto:daryl@goldenpath.ca) instead.
- Use relevant labels and feel free to suggest new ones that are helpful to the project.
- Avoid duplicate issues by searching for similar ones first. Add comments to existing issues if you can provide new and/or helpful information.
- Use [Discussions](https://github.com/goldenpathtechnologies/ci/discussions) for all other inquiries, such as ideas, shout outs, general questions, etc. If I really like an idea discussed here, it may be added to the issue backlog.
- I may create issue templates in the future, feel free to use them or not when available.

## Creating pull requests
You may want to skip creating an issue and simply submit a pull request if you know
how to fix a particular issue. This is perfectly fine as long as the following
guidelines are followed:

- Ensure your branches are from `dev` and not `main` or any other branch. The PR will otherwise be rejected.
- Ensure your fork is kept up to date with upstream changes. Use `git pull --rebase origin dev` by default when updating from upstream to avoid cluttering the commit history.
- Use [Conventional Commits](#conventional-commits) format when naming the pull request. The PR title is the first line of the commit, and the PR body is the remaining sections. This helps with automating releases.
- Ensure all non-squashed commits in the PR branch follow the Conventional Commits [spec](https://www.conventionalcommits.org/en/v1.0.0/). Otherwise, depending on the situation, the PR may be rejected, or I may squash the PR commits once the body contains a description of all relevant changes.
- Link issues related to the PR. Depending on permissions, this is something I may do instead.
- Review the [coding standards](#coding-standards-and-style) to ensure consistency throughout the project. 

## Software development process

### Repository information

#### Branches
`ci` has two primary branches, `main` and `dev`. The `HEAD` of the `main` branch is
synonymous with the current release. The `HEAD` of the `dev` branch is edge, but
not necessarily the latest development prerelease. The branching methodology used for
`ci` is similar to Gitflow, although I'm currently not using any tools to accomplish
this.

No changes are ever made directly to `main`. There are only two acceptable methods
of updating this branch:

1. Merges from release branches of `dev` (e.g. `release/*`)
2. Merges from hotfix branches of `main` (e.g. `hotfix/*`)

Enforcement of the above rules is currently not automated nor configured. `main` must
always build and must also maintain a code coverage threshold of 80% (see the section
on [Testing and coverage](#testing-and-coverage)). Since the `main` branch's purpose
is to contain the latest public release, it should not be branched off of for
development due to the risk of colliding with changes in the `dev` branch or of
duplicating code already present there. Any pull request from a branch of `main` that
is not for a hotfix will be rejected. Hotfix branches must be merged to `dev` at the
same time they are merged to `main`.

The `dev` branch is where all other feature branches get branched from. These branches
should have the following format whether they are for features for bug fixes:
`feature/<snake-case_description>`. The name of the branch should be as relevant to
the change as possible. For bug fix branches, the snake-case description should
contain the word 'fix' somewhere to indicate its intended purpose. `dev` can not be
merged into other than by pull requests or directly by project administrators.
Non-hotfix pull requests will only be accepted into `dev` and from branches of it.

A `test` branch may sometimes be present, alongside versioned tags labeled with
the eponymous release channel (e.g. `v1.0.0-test.1`). I usually create `test` to
test GitHub workflows and releases. `test` is a prototype branch that never gets
merged, and is deleted along with its tags and releases when my tests are complete.

#### Changelog
`CHANGELOG.md` is automatically generated and updated during CI/CD and should not be
modified otherwise. The changelog for `dev` is maintained separately from that of
`main`. That is, the `CHANGELOG.md` in `main` will never get changes merged from
`dev`. However, I'm still currently working on a way to do this automatically. For
now, I'll be checking out the `main` version of `CHANGELOG.md` when creating release
branches.

### Conventional commits
`ci` follows the [Conventional Commits spec](https://www.conventionalcommits.org/en/v1.0.0/)
for all commits in the repository. It is recommended to use a tool like
[commitizen](https://github.com/commitizen-tools/commitizen) to automate the creation
of properly formatted Conventional Commits. It is also acceptable to use multiple
`-m` flags on `git commit` to achieve the same result.

#### Breaking changes
Any pull request that contains breaking changes may either be rejected, or shelved
until the next major version is ready for release. Breaking changes are detected
either by the presence of keywords in the body of any commit (e.g. `BREAKING CHANGE`,
`BREAKING CHANGES`, or `BREAKING`) or determined via peer review.

### Coding standards and style
Please adhere as closely as possible to [Effective Go](https://go.dev/doc/effective_go)
coding standards. Admittedly, I let GoLand do most of this work for me as I code, and
therefore recommend it as the best IDE to use for this project. Additionally, here are
some basic guidelines to follow while working on `ci`:

- Readability is more important than brevity.
- Treat all typos as bugs, regardless of file type, test results, or build status.
- Use comments wisely. Leave TODO comments for anything you may not get to right away, or a problem someone else may be more fit to solve.
- Functions should be kept as small as possible and broken up in several parts when necessary. However, it's more important that the function does one thing well.
- Inline functions can be used provided they are simple. When they increase in complexity, they should be defined outside of the function they're used in and tested.

These rough guidelines are subject to change as the project evolves.

### Testing and coverage

#### Writing tests
Any code contributed must be accompanied by tests with 80% coverage being maintained
throughout the project. While coverage is important, tests should also be relevant.
That is, we may sometimes need more tests than will satisfy the coverage threshold
to ensure all scenarios are handled. However, if a function is too simple to test
(e.g. no implication logic and calling a single third-party function), then writing
tests for it is optional.

I deviate a bit from Go test naming conventions with the following format:

```go
// Format:
// Test_<struct or parent group name>_<func name>_<statement/description>

// Example:
func Test_DirectoryList_getDetailsText_DoesNotReturnOutputFromPreviousCall(t *testing.T) {...}
```

In this way, tests also act as a spec that describes in human language how  `ci` is
supposed to operate. My approach to testing in Go is not set in stone, and I'm open to
suggestions on how I can test better.

#### Running tests
If possible, tests should be run in both Windows and Linux. This is best accomplished
with GoLand [run targets](https://blog.jetbrains.com/go/2021/04/29/what-are-run-targets-and-how-to-run-code-anywhere/)
on a Windows machine that has Windows Subsystem for Linux (WSL) installed. Regardless,
tests for both platforms are run on each push and some PR events
(see [test workflow](https://github.com/goldenpathtechnologies/ci/blob/main/.github/workflows/test.yml)).

When not using GoLand, run the following commands to test `ci`:

```bash
# Test with coverage output in Linux
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# View coverage analysis
go tool cover -html coverage.out
```

```powershell
# Test with coverage output in Windows (PowerShell)
go test -v -race -coverprofile coverage.out -covermode=atomic ./...

# View coverage analysis
go tool cover -html coverage.out
```

### Building the project
Building is relatively straightforward. Just run the command for your OS environment
from the root repo directory and customize variables as needed:

#### Linux (Bash)

```bash
go build -gcflags="all=-N -l" -ldflags \
"-X 'main.BuildVersion=0.0.0-dev.0' -X 'main.BuildDate=2021-10-14T15:15:00Z' -X 'main.BuildOwner1=Daryl G. Wright' -X 'main.BuildOwner2=Golden Path Technologies Inc.'" \
-tags forceposix \
-o ./bin/ci
```

#### Windows (PowerShell)

```powershell
go build -gcflags="all=-N -l" -ldflags `
"-X 'main.BuildVersion=0.0.0-dev.0' -X 'main.BuildDate=2021-10-14T15:15:00Z' -X 'main.BuildOwner1=Daryl G. Wright' -X 'main.BuildOwner2=Golden Path Technologies Inc.'" `
-tags forceposix `
-o ./bin/ci.exe
```

One thing to note is that the `main.BuildDate` variable must be in RFC3339 format.
Otherwise, the build will fail.

### Development installation
During development, you may want to install and uninstall `ci` from your working
directory. You can run the following scripts from the root repo directory or the
scripts directory:

#### Linux (Bash)

 ```bash
# From root repo directory
./scripts/install.sh
./scripts/uninstall.sh

# From scripts directory
./install.sh
./uninstall.sh
 ```

#### Windows (PowerShell)

```powershell
# From root repo directory
.\scripts\install.ps1
.\scripts\uninstall.ps1

# From scripts directory
.\install.ps1
.\uninstall.ps1
```
