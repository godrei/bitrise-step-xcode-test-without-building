package step

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/godrei/bitrise-step-xcode-test-without-building/xcodebuild"
	"github.com/kballard/go-shellquote"
)

const (
	testResultBundleKey       = "BITRISE_XCRESULT_PATH"
	zippedTestResultBundleKey = "BITRISE_XCRESULT_ZIP_PATH"
)

type Input struct {
	Xctestrun         string `env:"xctestrun,required"`
	Destination       string `env:"destination,required"`
	XcodebuildOptions string `env:"xcodebuild_options"`
	DeployDir         string `env:"BITRISE_DEPLOY_DIR"`
}

type Config struct {
	Xctestrun         string
	Destination       string
	XcodebuildOptions []string
	DeployDir         string
}

type Result struct {
	TestOutputDir string
	DeployDir     string
}

type Step struct {
	logger         log.Logger
	inputParser    stepconf.InputParser
	xcodebuild     xcodebuild.Xcodebuild
	outputEnvStore env.Repository
	outputExporter OutputExporter
}

func New(logger log.Logger, inputParser stepconf.InputParser, xcodebuild xcodebuild.Xcodebuild, outputEnvStore env.Repository, outputExporter OutputExporter) Step {
	return Step{
		logger:         logger,
		inputParser:    inputParser,
		xcodebuild:     xcodebuild,
		outputEnvStore: outputEnvStore,
		outputExporter: outputExporter,
	}
}

func (s Step) ProcessConfig() (*Config, error) {
	var input Input
	if err := s.inputParser.Parse(&input); err != nil {
		return nil, err
	}

	stepconf.Print(input)

	xcodebuildOptions, err := shellquote.Split(input.XcodebuildOptions)
	if err != nil {
		return nil, fmt.Errorf("provided xcodebuild options (%s) are not valid CLI parameters: %w", input.XcodebuildOptions, err)
	}

	return &Config{
		Xctestrun:         input.Xctestrun,
		Destination:       input.Destination,
		XcodebuildOptions: xcodebuildOptions,
		DeployDir:         input.DeployDir,
	}, nil
}

func (s Step) InstallDependencies() error {
	return nil
}

func (s Step) Run(config Config) (*Result, error) {
	s.logger.Println()
	s.logger.Infof("Running tests:")

	result := &Result{
		DeployDir: config.DeployDir,
	}

	outputDir, err := s.xcodebuild.TestWithoutBuilding(config.Xctestrun, config.Destination, config.XcodebuildOptions...)
	result.TestOutputDir = outputDir

	if err != nil {
		var exerr *exec.ExitError
		if errors.As(err, &exerr) {
			return result, fmt.Errorf("failing tests (exist status %v)", exerr.ExitCode())
		} else {
			return result, fmt.Errorf("test execute failed: %w", err)
		}
	}

	s.logger.TDonef("Passing tests")
	return result, err
}

func (s Step) ExportOutputs(result Result) error {
	s.logger.Println()
	s.logger.Infof("Exporting outputs:")

	if result.TestOutputDir != "" {
		if err := s.outputEnvStore.Set(testResultBundleKey, result.TestOutputDir); err != nil {
			s.logger.Warnf("Failed to export: %s: %s", testResultBundleKey, err)
		}
		s.logger.Donef("%s: %s", testResultBundleKey, result.TestOutputDir)

		if result.DeployDir != "" {
			xcresultZipPath := filepath.Join(result.DeployDir, filepath.Base(result.TestOutputDir)+".zip")
			if err := s.outputExporter.ZipAndExportOutput(result.TestOutputDir, xcresultZipPath, zippedTestResultBundleKey); err != nil {
				s.logger.Warnf("Failed to export: %s: %s", zippedTestResultBundleKey, err)
			}
			s.logger.Donef("%s: %s", zippedTestResultBundleKey, xcresultZipPath)
		}
	}
	return nil
}
