[![Go Reference](https://pkg.go.dev/badge/github.com/Deichindianer/semver-go.svg)](https://pkg.go.dev/github.com/Deichindianer/semver-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/Deichindianer/semver-go?style=flat-square)](https://goreportcard.com/report/github.com/Deichindianer/semver-go)

# semver-go

"Get me the next patch, minor or major version of this git repository."

semver-go provides an interface to do that really easily without any dependencies for the cli. Currently this package does only support the patch, minor and major components of the [semantic versioning spec](https://semver.org).

This repository contains two things

* The `sem` package that provides an interface for SemVer versioning
* `semver` a cli that allows you to just get the next version easily from your command line or CI job shell

## CLI Usage

`semver` will keep your version prefixes if you have any. It will then print out the next version to stdout.

```sh
# git repository with the latest tag v1.0.0
$ semver patch
v1.0.1
```

## Sem package usage

Find most of the detailed docs on [pkg.go.dev](https://pkg.go.dev/github.com/Deichindianer/semver-go).

The main exposed functions are

* `GetLatestVersion` getting the latest version out of a git repository
* The `Ver` struct that provides an interface to your version and the `Next` function that increments the version in place 

## Credits

* [mycrEEpy](https://github.com/mycrEEpy) for developing this entire thing in tandem with me.
