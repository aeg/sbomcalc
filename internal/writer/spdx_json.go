package writer

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/eiji/sbomcalc/internal/model"
)

func WriteSPDXJSON(w io.Writer, result model.QueryResult, version string) error {
	doc := map[string]any{
		"spdxVersion":       "SPDX-" + version,
		"dataLicense":       "CC0-1.0",
		"SPDXID":            "SPDXRef-DOCUMENT",
		"name":              "sbomcalc-result",
		"documentNamespace": fmt.Sprintf("https://sbomcalc.local/doc/%d", time.Now().UnixNano()),
		"creationInfo": map[string]any{
			"created":  time.Now().UTC().Format(time.RFC3339),
			"creators": []string{"Tool: sbomcalc"},
		},
		"packages": spdxPackages(result.Components),
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(doc)
}

func spdxPackages(records []model.ComponentRecord) []map[string]any {
	packages := make([]map[string]any, 0, len(records))
	for i, record := range records {
		pkg := map[string]any{
			"name":                  record.Name,
			"SPDXID":                fmt.Sprintf("SPDXRef-Package-%d", i+1),
			"downloadLocation":      "NOASSERTION",
			"filesAnalyzed":         false,
			"licenseConcluded":      "NOASSERTION",
			"licenseDeclared":       "NOASSERTION",
			"copyrightText":         "NOASSERTION",
			"supplier":              supplierOrNoAssertion(record.Supplier),
			"externalRefs":          spdxExternalRefs(record),
			"checksums":             spdxChecksums(record.Hashes),
			"versionInfo":           record.Version,
			"primaryPackagePurpose": "LIBRARY",
		}
		if len(record.Licenses) > 0 {
			pkg["licenseConcluded"] = record.Licenses[0]
		}
		packages = append(packages, pkg)
	}
	return packages
}

func supplierOrNoAssertion(value string) string {
	if value == "" {
		return "NOASSERTION"
	}
	return value
}

func spdxExternalRefs(record model.ComponentRecord) []map[string]string {
	if record.PURL == "" {
		return nil
	}
	return []map[string]string{{
		"referenceCategory": "PACKAGE-MANAGER",
		"referenceType":     "purl",
		"referenceLocator":  record.PURL,
	}}
}

func spdxChecksums(hashes []model.Hash) []map[string]string {
	out := make([]map[string]string, 0, len(hashes))
	for _, hash := range hashes {
		out = append(out, map[string]string{"algorithm": hash.Algorithm, "checksumValue": hash.Value})
	}
	return out
}
