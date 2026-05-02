package reader

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/eiji/sbomcalc/internal/model"
)

type Format int

const (
	FormatSPDXJSON Format = iota
	FormatCycloneDXJSON
)

type Info struct {
	Format  Format
	Version string
}

type ComponentFunc func(model.ComponentRecord) error

func ScanFile(path string, fn ComponentFunc) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := Detect(file)
	if err != nil {
		return err
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	switch info.Format {
	case FormatSPDXJSON:
		return scanSPDXJSON(file, path, fn)
	case FormatCycloneDXJSON:
		return scanCycloneDXJSON(file, path, fn)
	default:
		return fmt.Errorf("unsupported SBOM format")
	}
}

func ScanKeySet(path string, level model.Level) (model.KeySet, error) {
	set := model.NewKeySet()
	err := ScanFile(path, func(record model.ComponentRecord) error {
		key := model.KeyFor(record, level)
		if key.Name != "" {
			set.Add(key)
		}
		return nil
	})
	return set, err
}

func skipValue(dec *json.Decoder) error {
	var raw json.RawMessage
	return dec.Decode(&raw)
}

func expectObject(dec *json.Decoder) error {
	token, err := dec.Token()
	if err != nil {
		return err
	}
	delim, ok := token.(json.Delim)
	if !ok || delim != '{' {
		return fmt.Errorf("top-level JSON value must be an object")
	}
	return nil
}
