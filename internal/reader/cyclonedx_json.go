package reader

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/aeg/sbomcalc/internal/model"
)

type cdxComponent struct {
	Name     string          `json:"name"`
	Version  string          `json:"version"`
	PURL     string          `json:"purl"`
	Supplier json.RawMessage `json:"supplier"`
	Licenses []cdxLicense    `json:"licenses"`
	Hashes   []cdxHash       `json:"hashes"`
}

type cdxLicense struct {
	License    cdxLicenseInfo `json:"license"`
	Expression string         `json:"expression"`
}

type cdxLicenseInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type cdxHash struct {
	Alg     string `json:"alg"`
	Content string `json:"content"`
}

func scanCycloneDXJSON(r io.Reader, source string, fn ComponentFunc) error {
	dec := json.NewDecoder(r)
	if err := expectObject(dec); err != nil {
		return err
	}

	for dec.More() {
		token, err := dec.Token()
		if err != nil {
			return err
		}
		key, ok := token.(string)
		if !ok {
			return fmt.Errorf("invalid CycloneDX JSON object")
		}
		if key != "components" {
			if err := skipValue(dec); err != nil {
				return err
			}
			continue
		}
		return scanCycloneDXComponents(dec, source, fn)
	}
	return nil
}

func scanCycloneDXComponents(dec *json.Decoder, source string, fn ComponentFunc) error {
	token, err := dec.Token()
	if err != nil {
		return err
	}
	delim, ok := token.(json.Delim)
	if !ok || delim != '[' {
		return fmt.Errorf("CycloneDX JSON components must be an array")
	}

	for dec.More() {
		var component cdxComponent
		if err := dec.Decode(&component); err != nil {
			return err
		}
		record := cdxComponentRecord(component, source)
		if record.Name == "" {
			continue
		}
		if err := fn(record); err != nil {
			return err
		}
	}
	_, err = dec.Token()
	return err
}

func cdxComponentRecord(component cdxComponent, source string) model.ComponentRecord {
	name := model.NormalizeName(component.Name)
	if name == "" {
		return model.ComponentRecord{}
	}

	record := model.ComponentRecord{
		Name:     name,
		Version:  model.NormalizeVersion(component.Version),
		PURL:     strings.TrimSpace(component.PURL),
		Supplier: cdxSupplier(component.Supplier),
		Source:   source,
	}
	for _, license := range component.Licenses {
		value := strings.TrimSpace(license.License.ID)
		if value == "" {
			value = strings.TrimSpace(license.License.Name)
		}
		if value == "" {
			value = strings.TrimSpace(license.Expression)
		}
		if value != "" {
			record.Licenses = append(record.Licenses, value)
		}
	}
	for _, hash := range component.Hashes {
		alg := strings.TrimSpace(hash.Alg)
		content := strings.TrimSpace(hash.Content)
		if alg != "" && content != "" {
			record.Hashes = append(record.Hashes, model.Hash{Algorithm: alg, Value: content})
		}
	}
	return record
}

func cdxSupplier(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return strings.TrimSpace(text)
	}
	var obj struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(raw, &obj); err == nil {
		return strings.TrimSpace(obj.Name)
	}
	return ""
}
