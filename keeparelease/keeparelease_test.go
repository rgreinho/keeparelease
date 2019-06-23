package keeparelease

import (
	"strings"
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

- Make the Git operations idempotent.
- Implement all the operations required to create an environment.
- Create an Octopus client to simplify the operation using its REST API.
- Add --dry-run flag to bypass the write operations.

[//]: # (Release links)
[0.3.0]: https://github.com/shipstation/hyperion/releases/tag/0.3.0
[0.4.0]: https://github.com/shipstation/hyperion/releases/tag/0.4.0

[//]: # (Issue/PR links)
[#15]: https://github.com/shipstation/hyperion/pull/15
[#16]: https://github.com/shipstation/hyperion/pull/16
`

func TestParseChangelog00(t *testing.T) {
	title, content, err := ParseChangelog(strings.Split(changelog, "\n"))
	if err != nil {
		t.Fatalf("failed to parse the changelog: %s", err)
	}
	expected := `

  ### Added

  - A feature

  ### Fixed

  - A fix
  `

	if dedent.Dedent(title) != dedent.Dedent("0.4.0") {
		t.Fatalf("Error: title is %s, but expected is 0.4.0", title)
	}
	if dedent.Dedent(content) != dedent.Dedent(expected) {
		t.Fatalf("Error: actual is %s, but expected is %s", content, expected)
	}

}
