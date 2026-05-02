package writer

import (
	"fmt"
	"io"
	"os"

	"github.com/eiji/sbomcalc/internal/cli"
	"github.com/eiji/sbomcalc/internal/model"
)

func WriteQuery(result model.QueryResult, specs []cli.OutputSpec) error {
	for _, spec := range specs {
		if err := writeTo(spec, func(w io.Writer) error {
			switch spec.Kind {
			case cli.OutputTable:
				return WriteQueryTable(w, result)
			case cli.OutputText:
				return WriteQueryText(w, result)
			case cli.OutputJSON:
				return WriteQueryJSON(w, result)
			case cli.OutputSPDXJSON:
				return WriteSPDXJSON(w, result, spec.Version)
			case cli.OutputCycloneDXJSON:
				return WriteCycloneDXJSON(w, result, spec.Version)
			default:
				return fmt.Errorf("unsupported output")
			}
		}); err != nil {
			return err
		}
	}
	return nil
}

func WriteDiff(result model.DiffResult, specs []cli.OutputSpec, changedOnly bool) error {
	for _, spec := range specs {
		if err := writeTo(spec, func(w io.Writer) error {
			switch spec.Kind {
			case cli.OutputTable:
				return WriteDiffTable(w, result, changedOnly)
			case cli.OutputText:
				return WriteDiffText(w, result, changedOnly)
			case cli.OutputJSON:
				return WriteDiffJSON(w, result, changedOnly)
			default:
				return fmt.Errorf("unsupported output")
			}
		}); err != nil {
			return err
		}
	}
	return nil
}

func writeTo(spec cli.OutputSpec, fn func(io.Writer) error) error {
	if spec.File == "" {
		return fn(os.Stdout)
	}
	file, err := os.Create(spec.File)
	if err != nil {
		return err
	}
	defer file.Close()
	return fn(file)
}
