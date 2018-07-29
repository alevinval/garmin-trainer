package provider

import (
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/alevinval/trainer/internal/adapter"
	"github.com/alevinval/trainer/internal/trainer"
)

type (
	adapterConstructor func(data []byte) (provider trainer.ActivityProvider, err error)
)

var (
	// ErrExtensionNotSupported is returned when the file extension is not supported
	ErrExtensionNotSupported = errors.New("unrecognized file extension")

	extToAdapter = map[string]adapterConstructor{
		".gpx": adapter.Gpx,
		".fit": adapter.Fit,
	}
)

// File reads a file content and returns an Activity.
func File(name string) (actvity *trainer.Activity, err error) {
	ext, isGzip := isGzip(name)
	if !isExtSupported(ext) {
		return nil, ErrExtensionNotSupported
	}
	r, err := getReaderForFile(name, isGzip)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	provider, err := extToAdapter[ext](data)
	if err != nil {
		return nil, err
	}
	return buildActivity(provider, data, name), nil
}

// isGzip returns the extension of a file and whether its zipped or not.
func isGzip(name string) (ext string, isGzip bool) {
	if strings.HasSuffix(name, ".gz") {
		nameWithoutGz := strings.TrimSuffix(name, ".gz")
		ext := path.Ext(nameWithoutGz)
		return ext, true
	}
	return path.Ext(name), false
}

// isExtSupported returns if an adapter exists for the given extension.
func isExtSupported(ext string) (supported bool) {
	_, ok := extToAdapter[ext]
	return ok
}

// getReaderForFile returns a reader for a file, supports zipped files.
func getReaderForFile(name string, isGzip bool) (r io.Reader, err error) {
	r, err = os.Open(name)
	if err != nil {
		return nil, err
	}
	if !isGzip {
		return r, err
	}
	return gzip.NewReader(r)
}

// buildActivity returns an activity with metadata reflecting the source
// of the activity.
func buildActivity(provider trainer.ActivityProvider, data []byte, fileName string) *trainer.Activity {
	metadata := provider.Metadata()
	metadata.DataSource = trainer.DataSource{
		Type: trainer.FileDataSource,
		Name: fileName,
	}
	datapoints := provider.DataPoints()
	return trainer.NewActivity(data, metadata, datapoints)
}