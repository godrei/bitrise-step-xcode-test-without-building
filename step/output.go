package step

import "github.com/bitrise-io/go-steputils/output"

type OutputExporter interface {
	ZipAndExportOutput(artifact string, destinationZipPth, envKey string) error
}

type outputExporter struct {
}

func NewOutputExporter() OutputExporter {
	return outputExporter{}
}

func (e outputExporter) ZipAndExportOutput(artifact string, destinationZipPth, envKey string) error {
	return output.ZipAndExportOutput([]string{artifact}, destinationZipPth, envKey)
}
