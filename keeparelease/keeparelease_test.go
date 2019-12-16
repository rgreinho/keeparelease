package keeparelease

import (
	"testing"

	"github.com/lithammer/dedent"
)

const semverChangelog string = `
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

const calverChangelog string = `
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Calendar Versioning](https://calver.org/).

## [Unreleased]

## Added

- Add a mechanism to configure the common tools using environment variables.

## [19.12.1]

### Added

- Add the following features to docgen-nodes-reboot:
  - Ability to skip the n first nodes.
  - Ability to target a specific node by ID.
  - Ability to exit on failure instead of moving forward to the next step.

### Changed

- Replace all the output statements by logging statements associated with the appropriate loggin level.

### Fixed

- Fix the AWS EC2 rate limit issue with docgen-nodes-reboot.
`

func TestParseChangelog00(t *testing.T) {
	testcases := []struct {
		changelog       string
		expectedTitle   string
		expectedContent string
	}{
		{semverChangelog, "0.4.0", trimEdges(dedent.Dedent(`
    ### Added

    - A feature

    ### Fixed

    - A fix
    `), " \n")},
		{calverChangelog, "19.12.1", trimEdges(dedent.Dedent(`
    ### Added

    - Add the following features to docgen-nodes-reboot:
      - Ability to skip the n first nodes.
      - Ability to target a specific node by ID.
      - Ability to exit on failure instead of moving forward to the next step.

    ### Changed

    - Replace all the output statements by logging statements associated with the appropriate loggin level.

    ### Fixed

    - Fix the AWS EC2 rate limit issue with docgen-nodes-reboot.
    `), " \n")},
	}

	for _, tc := range testcases {
		title, content, err := ParseChangelog(tc.changelog)
		if err != nil {
			t.Fatalf("failed to parse the changelog: %s", err)
		}

		if dedent.Dedent(title) != dedent.Dedent(tc.expectedTitle) {
			t.Fatalf("title is %q, but expected is %q", title, tc.expectedTitle)
		}
		if content != tc.expectedContent {
			t.Fatalf("content is %q, \nbut expected is %q.", content, tc.expectedContent)
		}
	}
}
