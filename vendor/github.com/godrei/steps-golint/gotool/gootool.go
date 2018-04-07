package gotool

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	glob "github.com/ryanuber/go-glob"
)

func parsePackages(out string, exclude ...string) []string {
	list := []string{}
	for _, l := range strings.Split(string(out), "\n") {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}

		match := false
		for _, e := range exclude {
			if glob.Glob(e, l) {
				match = true
				break
			}
		}
		if match {
			continue
		}

		list = append(list, l)
	}
	return list
}

// ListPackages ...
func ListPackages(dir string, exclude ...string) ([]string, error) {
	cmd := command.New("go", "list", "./...").SetDir(dir)
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%s failed: %s", cmd.PrintableCommandArgs(), out)
	}

	return parsePackages(out, exclude...), nil
}
