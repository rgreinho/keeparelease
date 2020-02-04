package keeparelease

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"

	"github.com/lithammer/dedent"
)

// re is the regex matching a semver in a markdown H2 header.
var re = regexp.MustCompile(`\#\#\s+\[{1,2}(?P<full_version>(?P<Major>0|\d+)\.(?P<Minor>0|\d+)\.(?P<Patch>0|\d+)(?P<PreReleaseTagWithSeparator>-(?P<PreReleaseTag>(?:[a-z-][\da-z-]+|[\da-z-]+[a-z-][\da-z-]*|0|[1-9]\d*)(?:\.(?:[a-z-][\da-z-]+|[\da-z-]+[a-z-][\da-z-]*|0|[1-9]\d*))*))?(?P<BuildMetadataWithSeparator>\+(?P<BuildMetadata>[\da-z-]+(?:\.[\da-z-]+)*))?)`)

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

	trimmed := trimEdges(dedent.Dedent(strings.Join(releaseInfo, "\n")), " \n")
	return title, trimmed, nil
}

// ReadChangelog reads the changelog file.
// Returns the title of the last release as well as its content.
func ReadChangelog(file string) (title, content string, err error) {
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

func trimEdges(s, cutset string) string {
	trimmed := strings.TrimLeft(s, " \n")
	trimmed = strings.TrimRight(trimmed, " \n")
	return trimmed
}

// GetTag retrieves the latest tag from a git repository.
func GetTag() (string, error) {
	cmd := exec.Command("git", "describe")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
