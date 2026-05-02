package writer

import (
	"encoding/json"
	"io"

	"github.com/aeg/sbomcalc/internal/model"
)

func WriteQueryJSON(w io.Writer, result model.QueryResult) error {
	type queryJSON struct {
		Level      string                  `json:"level"`
		Components []model.ComponentRecord `json:"components"`
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(queryJSON{Level: result.Level.String(), Components: result.Components})
}

func WriteDiffJSON(w io.Writer, result model.DiffResult, changedOnly bool) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if changedOnly {
		return enc.Encode(struct {
			From    string              `json:"from"`
			To      string              `json:"to"`
			Changed []model.ChangedName `json:"changed"`
		}{
			From: result.From, To: result.To, Changed: result.Changed,
		})
	}
	return enc.Encode(struct {
		From      string                `json:"from"`
		To        string                `json:"to"`
		Added     []model.VersionedName `json:"added"`
		Removed   []model.VersionedName `json:"removed"`
		Changed   []model.ChangedName   `json:"changed"`
		Unchanged []model.VersionedName `json:"unchanged"`
	}{
		From:      result.From,
		To:        result.To,
		Added:     result.Added,
		Removed:   result.Removed,
		Changed:   result.Changed,
		Unchanged: result.Unchanged,
	})
}
