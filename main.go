package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
)

// Config ...
type Config struct {
	Packages string `env:"packages,required"`
}

func installedInPath(name string) bool {
	cmd := exec.Command("which", name)
	outBytes, err := cmd.Output()
	return err == nil && strings.TrimSpace(string(outBytes)) != ""
}

func failf(format string, args ...interface{}) {
	log.Errorf(format, args...)
	os.Exit(1)
}

func main() {
	var cfg Config
	if err := stepconf.Parse(&cfg); err != nil {
		failf("Error: %s\n", err)
	}
	stepconf.Print(cfg)

	if !installedInPath("errcheck") {
		cmd := command.New("go", "get", "-u", "github.com/kisielk/errcheck")

		log.Infof("\nInstalling errcheck")
		log.Donef("$ %s", cmd.PrintableCommandArgs())

		if out, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
			failf("Failed to install errcheck: %s", out)
		}
	}

	packages := strings.Split(cfg.Packages, ",")

	log.Infof("\nRunning errcheck...")

	for _, p := range packages {
		cmd := command.NewWithStandardOuts("errcheck", "-asserts=true", "-blank=true", "-verbose", p)

		log.Printf("$ %s", cmd.PrintableCommandArgs())

		if err := cmd.Run(); err != nil {
			failf("errcheck failed: %s", err)
		}
	}
}
