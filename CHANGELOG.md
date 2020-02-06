# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [[1.3.0]] - 2020.02.06

## Changed

- [INTERNAL] Removes the `hub` dependency and perform all the GitHub operation using
  the `google/go-github` library.
- [INTERNAL] Improves error handling and error messages.
- Add a `--log` flag.

## Fix

- Fix tag auto-detection.

## [[1.2.2]] - 2020.02.04

### Changed

- Replace the [invoke](http://www.pyinvoke.org/) task runner with the [mage](https://magefile.org/) task runner. [#17]

### Fixed

- Fix the problem with the `--tag` flag which was not supposed to be mandatory. [#15]

## [[1.2.1]] - 2020.02.03

### Fixed

- Fix regex to take the dates starting with 2020 into consideration.

## [[1.2.0]] - 2020.01.06

### Added

- Add a flag allowing the user to specify the Changelog file to read from.

## [[1.1.2]] - 2019.12.30

### Fixed

- Fix `dist` task resulting in corrupt binaries.

## [[1.1.1]] - 2019-12-16

### Change

- Replace the `Makefile` by an `Invoke` task file.

### Fixed

- Fix the semver regex to make if work for calver as well.

## [[1.1.0]] - 2019-07-19

## [[1.0.0]] - 2019-06-23

Initial release.

[//]: # (Release links)
[1.0.0]: https://github.com/rgreinho/keeparelease/releases/tag/1.1.0
[1.1.0]: https://github.com/rgreinho/keeparelease/releases/tag/1.1.0
[1.1.1]: https://github.com/rgreinho/keeparelease/releases/tag/1.1.1
[1.1.2]: https://github.com/rgreinho/keeparelease/releases/tag/1.1.2
[1.1.2]: https://github.com/rgreinho/keeparelease/releases/tag/1.1.2
[1.2.0]: https://github.com/rgreinho/keeparelease/releases/tag/1.2.0
[1.2.1]: https://github.com/rgreinho/keeparelease/releases/tag/1.2.1
[1.3.0]: https://github.com/rgreinho/keeparelease/releases/tag/1.3.0

[//]: # (Issue/PR links)
[#12]: https://github.com/rgreinho/keeparelease/pull/12
[#15]: https://github.com/rgreinho/keeparelease/pull/15
[#17]: https://github.com/rgreinho/keeparelease/pull/17
