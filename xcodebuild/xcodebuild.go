package xcodebuild

import (
	"os"

	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/log"
)

type Xcodebuild interface {
	TestWithoutBuilding(xctestrun, destination string) error
}

type xcodebuild struct {
	commandFactory command.Factory
	logger         log.Logger
}

func New(commandFactory command.Factory, logger log.Logger) Xcodebuild {
	return xcodebuild{
		commandFactory: commandFactory,
		logger:         logger,
	}
}

func (x xcodebuild) TestWithoutBuilding(xctestrun, destination string) error {
	cmd := x.commandFactory.Create("xcodebuild", []string{"test-without-building", "-xctestrun", xctestrun, "-destination", destination}, &command.Opts{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	x.logger.TDonef(cmd.PrintableCommandArgs())
	return cmd.Run()
}
