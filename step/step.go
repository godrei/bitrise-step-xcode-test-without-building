package step

import (
	"fmt"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/godrei/bitrise-step-xcode-test-without-building/xcodebuild"
	"github.com/kballard/go-shellquote"
)

type Input struct {
	Xctestrun         string `env:"xctestrun,required"`
	Destination       string `env:"destination,required"`
	XcodebuildOptions string `env:"xcodebuild_options"`
}

type Config struct {
	Xctestrun         string
	Destination       string
	XcodebuildOptions []string
}

type Result struct {
}

type Step struct {
	logger         log.Logger
	inputParser    stepconf.InputParser
	outputEnvStore env.Repository
	xcodebuild     xcodebuild.Xcodebuild
}

func New(logger log.Logger, inputParser stepconf.InputParser, outputEnvStore env.Repository, xcodebuild xcodebuild.Xcodebuild) Step {
	return Step{
		logger:         logger,
		inputParser:    inputParser,
		outputEnvStore: outputEnvStore,
		xcodebuild:     xcodebuild,
	}
}

func (s Step) ProcessConfig() (*Config, error) {
	var input Input
	if err := s.inputParser.Parse(&input); err != nil {
		return nil, err
	}

	stepconf.Print(input)
	s.logger.Println()

	xcodebuildOptions, err := shellquote.Split(input.XcodebuildOptions)
	if err != nil {
		return nil, fmt.Errorf("provided xcodebuild options (%s) are not valid CLI parameters: %w", input.XcodebuildOptions, err)
	}

	return &Config{
		Xctestrun:         input.Xctestrun,
		Destination:       input.Destination,
		XcodebuildOptions: xcodebuildOptions,
	}, nil
}

func (s Step) InstallDependencies() error {
	return nil
}

func (s Step) Run(config Config) (*Result, error) {
	if err := s.xcodebuild.TestWithoutBuilding(config.Xctestrun, config.Destination, config.XcodebuildOptions...); err != nil {
		return nil, err
	}
	return &Result{}, nil
}

func (s Step) ExportOutputs(result Result) error {
	return nil
}
