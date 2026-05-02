package writer

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/eiji/sbomcalc/internal/model"
)

func WriteQueryTable(w io.Writer, result model.QueryResult) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	if result.Level == model.Level1 {
		if _, err := fmt.Fprintln(tw, "NAME"); err != nil {
			return err
		}
		for _, component := range result.Components {
			if _, err := fmt.Fprintln(tw, component.Name); err != nil {
				return err
			}
		}
		return tw.Flush()
	}

	if _, err := fmt.Fprintln(tw, "NAME\tVERSION\tPURL"); err != nil {
		return err
	}
	for _, component := range result.Components {
		if _, err := fmt.Fprintf(tw, "%s\t%s\t%s\n", component.Name, component.Version, component.PURL); err != nil {
			return err
		}
	}
	return tw.Flush()
}

func WriteDiffTable(w io.Writer, result model.DiffResult, changedOnly bool) error {
	if changedOnly {
		return writeChangedSection(w, result.Changed, false)
	}
	if err := writeVersionedSection(w, "ADDED", result.Added); err != nil {
		return err
	}
	if err := writeVersionedSection(w, "REMOVED", result.Removed); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "CHANGED"); err != nil {
		return err
	}
	if err := writeChangedSection(w, result.Changed, true); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	return writeVersionedSection(w, "UNCHANGED", result.Unchanged)
}

func writeVersionedSection(w io.Writer, title string, values []model.VersionedName) error {
	if _, err := fmt.Fprintln(w, title); err != nil {
		return err
	}
	for _, value := range values {
		for _, version := range value.Versions {
			if _, err := fmt.Fprintf(w, "  %s@%s\n", value.Name, version); err != nil {
				return err
			}
		}
	}
	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	return nil
}

func writeChangedSection(w io.Writer, values []model.ChangedName, indent bool) error {
	for _, value := range values {
		prefix := ""
		childPrefix := "  "
		if indent {
			prefix = "  "
			childPrefix = "    "
		}
		if _, err := fmt.Fprintf(w, "%s%s\n", prefix, value.Name); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "%sold: %s\n", childPrefix, strings.Join(value.OldVersion, ", ")); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "%snew: %s\n", childPrefix, strings.Join(value.NewVersion, ", ")); err != nil {
			return err
		}
	}
	return nil
}
