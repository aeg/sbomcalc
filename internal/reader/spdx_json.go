package reader

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/aeg/sbomcalc/internal/model"
)

type spdxPackage struct {
	SPDXID           string            `json:"SPDXID"`
	Name             string            `json:"name"`
	VersionInfo      string            `json:"versionInfo"`
	Supplier         string            `json:"supplier"`
	LicenseConcluded string            `json:"licenseConcluded"`
	ExternalRefs     []spdxExternalRef `json:"externalRefs"`
	Checksums        []spdxChecksum    `json:"checksums"`
}

type spdxExternalRef struct {
	ReferenceType    string `json:"referenceType"`
	ReferenceLocator string `json:"referenceLocator"`
}

type spdxChecksum struct {
	Algorithm     string `json:"algorithm"`
	ChecksumValue string `json:"checksumValue"`
}

func scanSPDXJSON(r io.Reader, source string, fn ComponentFunc) error {
	dec := json.NewDecoder(r)
	if err := expectObject(dec); err != nil {
		return err
	}

	foundPackages := false
	for dec.More() {
		token, err := dec.Token()
		if err != nil {
			return err
		}
		key, ok := token.(string)
		if !ok {
			return fmt.Errorf("invalid SPDX JSON object")
		}
		if key != "packages" {
			if err := skipValue(dec); err != nil {
				return err
			}
			continue
		}
		foundPackages = true
		if err := scanSPDXPackages(dec, source, fn); err != nil {
			return err
		}
	}
	if !foundPackages {
		return fmt.Errorf("SPDX JSON packages must be an array")
	}
	return nil
}

func scanSPDXPackages(dec *json.Decoder, source string, fn ComponentFunc) error {
	token, err := dec.Token()
	if err != nil {
		return err
	}
	delim, ok := token.(json.Delim)
	if !ok || delim != '[' {
		return fmt.Errorf("SPDX JSON packages must be an array")
	}

	for dec.More() {
		var pkg spdxPackage
		if err := dec.Decode(&pkg); err != nil {
			return err
		}
		record := spdxPackageRecord(pkg, source)
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

func spdxPackageRecord(pkg spdxPackage, source string) model.ComponentRecord {
	if pkg.SPDXID == "SPDXRef-DOCUMENT" {
		return model.ComponentRecord{}
	}
	name := model.NormalizeName(pkg.Name)
	if name == "" {
		return model.ComponentRecord{}
	}

	record := model.ComponentRecord{
		Name:     name,
		Version:  model.NormalizeVersion(pkg.VersionInfo),
		PURL:     firstSPDXPURL(pkg.ExternalRefs),
		Supplier: normalizeUnknown(pkg.Supplier),
		Source:   source,
	}
	if license := normalizeLicense(pkg.LicenseConcluded); license != "" {
		record.Licenses = []string{license}
	}
	for _, checksum := range pkg.Checksums {
		algorithm := strings.TrimSpace(checksum.Algorithm)
		value := strings.TrimSpace(checksum.ChecksumValue)
		if algorithm != "" && value != "" {
			record.Hashes = append(record.Hashes, model.Hash{Algorithm: algorithm, Value: value})
		}
	}
	return record
}

func firstSPDXPURL(refs []spdxExternalRef) string {
	for _, ref := range refs {
		if strings.TrimSpace(ref.ReferenceType) == "purl" {
			return strings.TrimSpace(ref.ReferenceLocator)
		}
	}
	return ""
}

func normalizeUnknown(value string) string {
	value = strings.TrimSpace(value)
	if value == "NOASSERTION" {
		return ""
	}
	return value
}

func normalizeLicense(value string) string {
	value = strings.TrimSpace(value)
	switch value {
	case "", "NOASSERTION", "NONE":
		return ""
	default:
		return value
	}
}
