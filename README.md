# Keeparelease

Publish beautiful GitHub releases with "keep a changelog".

## Goal

This tool was made to simplify creating beatiful releases on GitHub for projects following the
"[Keep A Changelog](https://keepachangelog.com/en/1.0.0/)" specification to write their Changelogs.

It aims to make it easy and convenient to publish the relevant Changelog section as the release description and
to attach the related assets.

Here is an example:

![keeparelease example](docs/img/keep-a-release.png)

## Installation

Simply download the binary for your platform from [the release page](https://github.com/rgreinho/keeparelease/releases).

## Usage

Creating a release and attaching all binaries:

```bash
keeparelease -t 1.1.1 \
  -a keeparelease-1.1.1-darwin-amd64 \
  -a keeparelease-1.1.1-linux-amd64 \
  -a keeparelease-1.1.1-windows-amd64
```

Print the changelog for the latest release:

```bash
keeparelease -x
```
