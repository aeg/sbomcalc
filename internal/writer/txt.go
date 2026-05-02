package writer

import (
	"fmt"
	"io"

	"github.com/aeg/sbomcalc/internal/model"
)

func WriteQueryText(w io.Writer, result model.QueryResult) error {
	for _, component := range result.Components {
		if result.Level == model.Level1 {
			if _, err := fmt.Fprintln(w, component.Name); err != nil {
				return err
			}
			continue
		}
		if _, err := fmt.Fprintf(w, "%s@%s\n", component.Name, component.Version); err != nil {
			return err
		}
	}
	return nil
}

func WriteDiffText(w io.Writer, result model.DiffResult, changedOnly bool) error {
	return WriteDiffTable(w, result, changedOnly)
}
