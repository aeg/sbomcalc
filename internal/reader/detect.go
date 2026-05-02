package reader

import (
	"encoding/json"
	"fmt"
	"io"
)

func Detect(r io.Reader) (Info, error) {
	dec := json.NewDecoder(r)
	if err := expectObject(dec); err != nil {
		return Info{}, err
	}

	var spdxVersion string
	var bomFormat string
	var specVersion string

	for dec.More() {
		token, err := dec.Token()
		if err != nil {
			return Info{}, err
		}
		key, ok := token.(string)
		if !ok {
			return Info{}, fmt.Errorf("invalid top-level JSON object")
		}
		switch key {
		case "spdxVersion":
			if err := dec.Decode(&spdxVersion); err != nil {
				return Info{}, err
			}
		case "bomFormat":
			if err := dec.Decode(&bomFormat); err != nil {
				return Info{}, err
			}
		case "specVersion":
			if err := dec.Decode(&specVersion); err != nil {
				return Info{}, err
			}
		default:
			if err := skipValue(dec); err != nil {
				return Info{}, err
			}
		}
	}

	switch spdxVersion {
	case "SPDX-2.2":
		return Info{Format: FormatSPDXJSON, Version: "2.2"}, nil
	case "SPDX-2.3":
		return Info{Format: FormatSPDXJSON, Version: "2.3"}, nil
	}

	if bomFormat == "CycloneDX" {
		switch specVersion {
		case "1.5", "1.6", "1.7":
			return Info{Format: FormatCycloneDXJSON, Version: specVersion}, nil
		}
	}

	return Info{}, fmt.Errorf("unsupported SBOM format")
}
