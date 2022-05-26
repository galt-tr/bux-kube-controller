# BUX: Kubernetes Controller
> Manage your [BUX](https://getbux.io) deployments using Kubernetes

[![Release](https://img.shields.io/github/release-pre/BuxOrg/bux-kube-controller.svg?logo=github&style=flat&v=1)](https://github.com/BuxOrg/bux-kube-controller/releases)
[![Build Status](https://img.shields.io/github/workflow/status/BuxOrg/bux-kube-controller/run-go-tests?logo=github&v=1)](https://github.com/BuxOrg/bux-kube-controller/actions)
[![Report](https://goreportcard.com/badge/github.com/BuxOrg/bux-kube-controller?style=flat&v=1)](https://goreportcard.com/report/github.com/BuxOrg/bux-kube-controller)
[![codecov](https://codecov.io/gh/BuxOrg/bux-kube-controller/branch/master/graph/badge.svg?v=1)](https://codecov.io/gh/BuxOrg/bux-kube-controller)
[![Mergify Status](https://img.shields.io/endpoint.svg?url=https://api.mergify.com/v1/badges/BuxOrg/bux-kube-controller&style=flat&v=1)](https://mergify.io)
[![Go](https://img.shields.io/github/go-mod/go-version/BuxOrg/bux-kube-controller?v=1)](https://golang.org/)
<br>
[![Gitpod Ready-to-Code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/BuxOrg/bux-kube-controller)
[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg?style=flat)](https://github.com/RichardLitt/standard-readme)
[![Makefile Included](https://img.shields.io/badge/Makefile-Supported%20-brightgreen?=flat&logo=probot)](Makefile)
[![Sponsor](https://img.shields.io/badge/sponsor-BuxOrg-181717.svg?logo=github&style=flat&v=1)](https://github.com/sponsors/BuxOrg)
[![Donate](https://img.shields.io/badge/donate-bitcoin-ff9900.svg?logo=bitcoin&style=flat&v=1)](https://gobitcoinsv.com/#sponsor?utm_source=github&utm_medium=sponsor-link&utm_campaign=bux-kube-controller&utm_term=bux-kube-controller&utm_content=bux-kube-controller)

<br/>

## Table of Contents
- [What is BUX?](#what-is-bux)
- [Installation](#installation)
- [Configuration](#configuration)
- [Examples & Tests](#examples--tests)
- [Benchmarks](#benchmarks)
- [Code Standards](#code-standards)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

<br/>

## What is BUX?
[Read more about BUX](https://getbux.io)

<br/>

## Prerequisites
If you want to use redis, this controller assumed you have installed the [redis
operator](https://github.com/spotahome/redis-operator). In addition, if you
with to use MongoDB this controller assumes you have installed the [mongo
community operator](https://github.com/mongodb/mongodb-kubernetes-operator).

In addition, if you wish to deploy bux to a domain with SSL, this controller
assumes you have installed the [nginx ingress
controller](https://kubernetes.github.io/ingress-nginx/) as well as [cert
manager](https://cert-manager.io/).

<br/>


## Installation

**bux-kube-controller** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).

To install the CRDs and deploy the controller to the cluster, run:
```shell script
make deploy
```

### Run controller locally

To run the controller locally for development, first install the CRDs:
```bash
make install
```

Then run the controller:
```bash
make run
```

<br/>

## Configuration

There is a limited set of configurable spec fields on the Bux CR. Eventually,
we will enable the ability to set the entire bux config in the CR, but for now
the following list are the available parameters:

| Key           | Type     | Description                                 |
|---------------|----------|---------------------------------------------|
| configuration | `Object` | Bux configuration settings                  |
| domain        | `string` | Domain to deploy bux to                     |
| clusterIssuer | `string` | Name of cluster issuer object for SSL certs |
| console       | `bool`   | Enable bux-console provisioning             |

<details>
<summary><strong><code>Repository Features</code></strong></summary>
<br/>

This repository was created using [MrZ's `go-template`](https://github.com/mrz1836/go-template#about)

#### Built-in Features
- Continuous integration via [GitHub Actions](https://github.com/features/actions)
- Build automation via [Make](https://www.gnu.org/software/make)
- Dependency management using [Go Modules](https://github.com/golang/go/wiki/Modules)
- Code formatting using [gofumpt](https://github.com/mvdan/gofumpt) and linting with [golangci-lint](https://github.com/golangci/golangci-lint) and [yamllint](https://yamllint.readthedocs.io/en/stable/index.html)
- Unit testing with [testify](https://github.com/stretchr/testify), [race detector](https://blog.golang.org/race-detector), code coverage [HTML report](https://blog.golang.org/cover) and [Codecov report](https://codecov.io/)
- Releasing using [GoReleaser](https://github.com/goreleaser/goreleaser) on [new Tag](https://git-scm.com/book/en/v2/Git-Basics-Tagging)
- Dependency scanning and updating thanks to [Dependabot](https://dependabot.com) and [Nancy](https://github.com/sonatype-nexus-community/nancy)
- Security code analysis using [CodeQL Action](https://docs.github.com/en/github/finding-security-vulnerabilities-and-errors-in-your-code/about-code-scanning)
- Automatic syndication to [pkg.go.dev](https://pkg.go.dev/) on every release
- Generic templates for [Issues and Pull Requests](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/configuring-issue-templates-for-your-repository) in Github
- All standard Github files such as `LICENSE`, `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, and `SECURITY.md`
- Code [ownership configuration](.github/CODEOWNERS) for Github
- All your ignore files for [vs-code](.editorconfig), [docker](.dockerignore) and [git](.gitignore)
- Automatic sync for [labels](.github/labels.yml) into Github using a pre-defined [configuration](.github/labels.yml)
- Built-in powerful merging rules using [Mergify](https://mergify.io/)
- Welcome [new contributors](.github/mergify.yml) on their first Pull-Request
- Follows the [standard-readme](https://github.com/RichardLitt/standard-readme/blob/master/spec.md) specification
- [Visual Studio Code](https://code.visualstudio.com) configuration with [Go](https://code.visualstudio.com/docs/languages/go)
- (Optional) [Slack](https://slack.com), [Discord](https://discord.com) or [Twitter](https://twitter.com) announcements on new Github Releases
- (Optional) Easily add [contributors](https://allcontributors.org/docs/en/bot/installation) in any Issue or Pull-Request

</details>

<details>
<summary><strong><code>Package Dependencies</code></strong></summary>
<br/>

- [stretchr/testify](https://github.com/stretchr/testify)
</details>

<details>
<summary><strong><code>Library Deployment</code></strong></summary>
<br/>

Releases are automatically created when you create a new [git tag](https://git-scm.com/book/en/v2/Git-Basics-Tagging)!

If you want to manually make releases, please install GoReleaser:

[goreleaser](https://github.com/goreleaser/goreleaser) for easy binary or library deployment to Github and can be installed:
- **using make:** `make install-releaser`
- **using brew:** `brew install goreleaser`

The [.goreleaser.yml](.goreleaser.yml) file is used to configure [goreleaser](https://github.com/goreleaser/goreleaser).

<br/>

### Automatic Releases on Tag Creation (recommended)
Automatic releases via [Github Actions](.github/workflows/release.yml) from creating a new tag:
```shell
make tag version=1.2.3
```

<br/>

### Manual Releases (optional)
Use `make release-snap` to create a snapshot version of the release, and finally `make release` to ship to production (manually).

<br/>

</details>

<details>
<summary><strong><code>Makefile Commands</code></strong></summary>
<br/>

View all `makefile` commands
```shell script
make help
```

List of all current commands:
```text
bench            Run all benchmarks in the Go application
build            Build manager binary.
build-go         Build the Go application (locally)
clean-mods       Remove all the Go mod cache
controller-gen   Download controller-gen locally if necessary.
coverage         Shows the test coverage
deploy           Deploy controller to the K8s cluster specified in ~/.kube/config.
diff             Show the git diff
docker-build     Build docker image with the manager.
docker-push      Push docker image with the manager.
envtest          Download envtest-setup locally if necessary.
fmt              Run go fmt against code.
generate         Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
generate         Runs the go generate command in the base of the repo
godocs           Sync the latest tag with GoDocs
install          Install CRDs into the K8s cluster specified in ~/.kube/config.
install-all-contributors  Installs all contributors locally
install-go       Install the application (Using Native Go)
install-releaser  Install the GoReleaser application
kustomize        Download kustomize locally if necessary.
lint             Run the golangci-lint application (install if not found)
manifests        Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
release          Full production release (creates release in Github)
release-snap     Test the full release (build binaries)
release-test     Full production test release (everything except deploy)
replace-version  Replaces the version in HTML/JS (pre-deploy)
run              Run a controller from your host.
tag              Generate a new tag and push (tag version=0.0.0)
tag-remove       Remove a tag if found (tag-remove version=0.0.0)
tag-update       Update an existing tag to current commit (tag-update version=0.0.0)
test             Runs lint and ALL tests
test-ci          Runs all tests via CI (exports coverage)
test-ci-no-race  Runs all tests via CI (no race) (exports coverage)
test-ci-short    Runs unit tests via CI (exports coverage)
test-no-lint     Runs just tests
test-short       Runs vet, lint and tests (excludes integration tests)
test-unit        Runs tests and outputs coverage
undeploy         Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
uninstall        Uninstall CRDs from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
update           Update all project dependencies
update-contributors  Regenerates the contributors html/list
update-linter    Update the golangci-lint package (macOS only)
update-releaser  Update the goreleaser application
vet              Run go vet against code.
```
</details>

<br/>

## Examples & Tests
All unit tests and [examples](examples) run via [Github Actions](https://github.com/BuxOrg/bux-kube-controller/actions) and
uses [Go version 1.16.x](https://golang.org/doc/go1.16). View the [configuration file](.github/workflows/run-tests.yml).

<br/>

Run all tests (including integration tests)
```shell script
make test
```

<br/>

Run tests (excluding integration tests)
```shell script
make test-short
```

<br/>

## Benchmarks
Run the Go benchmarks:
```shell script
make bench
```

<br/>

## Code Standards
Read more about this Go project's [code standards](.github/CODE_STANDARDS.md).

<br/>

## Usage
Checkout all the [examples](examples)!

<br/>

## Contributing
View the [contributing guidelines](.github/CONTRIBUTING.md) and follow the [code of conduct](.github/CODE_OF_CONDUCT.md).

<br/>

### How can I help?
All kinds of contributions are welcome :raised_hands:!
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:.
You can also support this project by [becoming a sponsor on GitHub](https://github.com/sponsors/mrz1836) :clap:
or by making a [**bitcoin donation**](https://gobitcoinsv.com/#sponsor?utm_source=github&utm_medium=sponsor-link&utm_campaign=bux-kube-controller&utm_term=bux-kube-controller&utm_content=bux-kube-controller) to ensure this journey continues indefinitely! :rocket:

[![Stars](https://img.shields.io/github/stars/BuxOrg/bux-kube-controller?label=Please%20like%20us&style=social)](https://github.com/BuxOrg/bux-kube-controller/stargazers)

<br/>

### Contributors ✨
Thank you to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="https://github.com/galt-tr"><img src="https://avatars.githubusercontent.com/u/64976002?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Dylan</b></sub></a><br /><a href="#infra-galt-tr" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a> <a href="#maintenance-galt-tr" title="Maintenance">🚧</a> <a href="https://github.com/BuxOrg/bux-kube-controller/commits?author=galt-tr" title="Code">💻</a></td>
    <td align="center"><a href="https://mrz1818.com"><img src="https://avatars.githubusercontent.com/u/3743002?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Mr. Z</b></sub></a><br /><a href="#infra-mrz1836" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a> <a href="https://github.com/BuxOrg/bux-kube-controller/commits?author=mrz1836" title="Code">💻</a> <a href="#maintenance-mrz1836" title="Maintenance">🚧</a> <a href="#security-mrz1836" title="Security">🛡️</a></td>
  </tr>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

> This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification.

<br/>

## License

[![License](https://img.shields.io/github/license/BuxOrg/bux-kube-controller.svg?style=flat&v=1)](LICENSE)
