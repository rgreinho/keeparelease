package keeparelease

import (
	"errors"
	"io/ioutil"
	"regexp"
	"strings"
)

const regexTitle string = `^\#\#\s+\[{1,2}(?P<full_version>(?P<major>(?:0|[1-9][0-9]*))\.(?P<minor>(?:0|[1-9][0-9]*))\.(?P<patch>(?:0|[1-9][0-9]*))(\-(?P<prerelease>(?:0|[1-9A-Za-z-][0-9A-Za-z-]*)(\.(?:0|[1-9A-Za-z-][0-9A-Za-z-]*))*))?(\+(?P<build>[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*))?)\]{1,2}\s+-?\s+(?P<date>.*)`

// ParseChangelog parses a ChangeLog respecting the Keep A Changelog format.
func ParseChangelog(changelog []string) (string, string, error) {
	re := regexp.MustCompile(regexTitle)
	title := ""
	releaseInfo := make([]string, 1)
	foundStart := false
	for _, line := range changelog {

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
func ReadChangelog() (string, string, error) {
	dat, err := ioutil.ReadFile("CHANGELOG.md")
	if err != nil {
		return "", "", err
	}
	content := string(dat)
	title, info, err := ParseChangelog(strings.Split(content, "\n"))
	if err != nil {
		return "", "", err
	}
	return title, info, nil
}
