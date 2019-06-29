package keeparelease

import (
	"errors"
	"io/ioutil"
	"regexp"
	"strings"
)

// re is the regex matching a semver in a markdown H2 header.
var re = regexp.MustCompile(`^\#\#\s+\[\[{1,2}(?P<full_version>(?P<major>(?:0|[1-9][0-9]*))\.(?P<minor>(?:0|[1-9][0-9]*))\.(?P<patch>(?:0|[1-9][0-9]*))(\-(?P<prerelease>(?:0|[1-9A-Za-z-][0-9A-Za-z-]*)(\.(?:0|[1-9A-Za-z-][0-9A-Za-z-]*))*))?(\+(?P<build>[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*))?)\]{1,2}\s+-?\s+(?P<date>.*)`)

// ParseChangelog parses a ChangeLog respecting the Keep A Changelog format.
// Returns the title of the last release as well as its content.
func ParseChangelog(changelog string) (title, content string, err error) {
	changelogLines := strings.Split(changelog, "\n")
	releaseInfo := make([]string, 1)
	foundStart := false
	for _, line := range changelogLines {

		// Look for a release line.
		if re.MatchString(line) {
			if !foundStart {

				// Indicate we found the first line of the last release information.
				foundStart = true
				match := re.FindStringSubmatch(line)
				result := make(map[string]string)
				for i, name := range re.SubexpNames() {
					if i != 0 && name != "" {
						result[name] = match[i]
					}
				}

				// Add title line.
				title = result["full_version"]
				continue

			} else {
				// We reached the last line of the last release information.
				break
			}
		}
		if foundStart {
			// Ensure the first line is followed by a blank line.
			if len(releaseInfo) == 1 && line != "" {
				releaseInfo = append(releaseInfo, "")
			}

			// Append the content.
			releaseInfo = append(releaseInfo, line)
		}
	}
	if len(releaseInfo) == 0 {
		return "", "", errors.New("could not extract release information")
	}

	return title, strings.Join(releaseInfo, "\n"), nil
}

// ReadChangelog reads the changelog file.
// Returns the title of the last release as well as its content.
func ReadChangelog(file string) (title, content string, err error) {
	if file == "" {
		file = "CHANGELOG.md"
	}
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return "", "", err
	}

	changes := string(dat)
	title, content, err = ParseChangelog(changes)
	if err != nil {
		return "", "", err
	}
	return title, content, nil
}
