package keeparelease

import (
	"testing"

	"github.com/lithammer/dedent"
)

const changelog string = `
# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [[0.4.0]] - 2018-09-07

### Added

- A feature

### Fixed

- A fix

## [[0.3.0]] - 2018-09-05

### Added

- Make the Git operations idempotent. [#1]
- Implement all the operations required to create an environment. [#2]
- Add --dry-run flag to bypass the write operations.

[//]: # (Release links)
[0.3.0]: https://github.com/rgreinho/keeparelease/releases/tag/0.1.0
[0.4.0]: https://github.com/rgreinho/keeparelease/releases/tag/0.2.0

[//]: # (Issue/PR links)
[#1]: https://github.com/rgreinho/keeparelease/pull/1
[#2]: https://github.com/rgreinho/keeparelease/pull/2
`

func TestParseChangelog00(t *testing.T) {
	title, content, err := ParseChangelog(changelog)
	if err != nil {
		t.Fatalf("failed to parse the changelog: %s", err)
	}
	expectedTitle := "0.4.0"
	expectedContent := trimEdges(dedent.Dedent(`
  ### Added

  - A feature

  ### Fixed

  - A fix
  `), " \n")

	if dedent.Dedent(title) != dedent.Dedent(expectedTitle) {
		t.Fatalf("Error: title is %s, but expected is %s", title, expectedTitle)
	}
	if content != expectedContent {
		t.Fatalf("Error: content is %s, \nbut expected is %s.", content, expectedContent)
	}

}
