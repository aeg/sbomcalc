package writer

import (
	"encoding/json"
	"io"

	"github.com/eiji/sbomcalc/internal/model"
)

func WriteCycloneDXJSON(w io.Writer, result model.QueryResult, version string) error {
	doc := map[string]any{
		"bomFormat":   "CycloneDX",
		"specVersion": version,
		"version":     1,
		"components":  cdxComponents(result.Components),
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(doc)
}

func cdxComponents(records []model.ComponentRecord) []map[string]any {
	components := make([]map[string]any, 0, len(records))
	for _, record := range records {
		component := map[string]any{
			"type":    "library",
			"name":    record.Name,
			"version": record.Version,
		}
		if record.PURL != "" {
			component["purl"] = record.PURL
		}
		if record.Supplier != "" {
			component["supplier"] = map[string]string{"name": record.Supplier}
		}
		if len(record.Licenses) > 0 {
			component["licenses"] = cdxLicenses(record.Licenses)
		}
		if len(record.Hashes) > 0 {
			component["hashes"] = cdxHashes(record.Hashes)
		}
		components = append(components, component)
	}
	return components
}

func cdxLicenses(licenses []string) []map[string]any {
	out := make([]map[string]any, 0, len(licenses))
	for _, license := range licenses {
		out = append(out, map[string]any{"license": map[string]string{"id": license}})
	}
	return out
}

func cdxHashes(hashes []model.Hash) []map[string]string {
	out := make([]map[string]string, 0, len(hashes))
	for _, hash := range hashes {
		out = append(out, map[string]string{"alg": hash.Algorithm, "content": hash.Value})
	}
	return out
}
