[![Go Reference](https://pkg.go.dev/badge/github.com/catouc/semver-go.svg)](https://pkg.go.dev/github.com/catouc/semver-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/catouc/semver-go?style=flat-square)](https://goreportcard.com/report/github.com/catouc/semver-go)

# semver-go

"Get me the next patch, minor or major version of this git repository."

semver-go provides an interface to do that really easily without any dependencies for the cli. Currently, this package does support the patch, minor and major components of the [semantic versioning spec](https://semver.org).

This repository contains two things

* The `sem` package that provides an interface for SemVer versioning
* `semver` a cli that allows you to just get the next version easily from your command line or CI job shell

## Installation

### Download the release from GitHub

Pick the latest release from the [releases page](https://github.com/catouc/semver-go/releases) and download the binary for your system.
Then put it into your `PATH`.

### Nix

I'm exposing a flake currently that you can install, but I am looking to add this to nixpkgs soon so look out for that :)
Here's a sample nixpkgs overlay though:

```nix
let 
  system = "x86_64-linux";
  pkgs = import nixpkgs {
    inherit system;
    overlays = [
      (final: prev: {semver-go = semver-go.packages.${system}.semver-go;})
    ];
  };
```

## 1.0.0 notice

This project is considered feature complete from my side of things, it does everything I need it to do.
If there are any issues I will of course happily fix them up :)
I did consider adding prerelease version parsing but decided against it because I never needed it and the cases for it seem like situations where I'd not want to automatically bump versions anyway.

## CLI Usage

`semver` will keep your version prefixes if you have any. It will then print out the next version to stdout.

```sh
# git repository with the latest tag v1.0.0
$ semver patch
v1.0.1
```

If you are facing issues with a few tags not being compliant with SemVer you can use the `-i` flag to ignore parsing errors.

```sh
# git repository with a tag like v.2.0.0 and a latest "real" tag like v1.0.0
$ semver -i patch
v1.0.1
```

## Sem package usage

Find most of the detailed docs on [pkg.go.dev](https://pkg.go.dev/github.com/catouc/semver-go).

The main exposed functions are

* `GetLatestVersion` getting the latest version out of a git repository
* The `Ver` struct that provides an interface to your version and the `Next` function that increments the version in place 

## Credits

* [mycrEEpy](https://github.com/mycrEEpy) for developing this entire thing in tandem with me.
