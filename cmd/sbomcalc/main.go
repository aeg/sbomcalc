package main

import (
	"fmt"
	"os"

	"github.com/aeg/sbomcalc/internal/cli"
	"github.com/aeg/sbomcalc/internal/engine"
	"github.com/aeg/sbomcalc/internal/writer"
	"github.com/spf13/cobra"
)

func main() {
	if err := newRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:           "sbomcalc",
		Short:         "Calculate component sets from SBOM files",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(newQueryCommand(), newDiffCommand(false), newDiffCommand(true))
	return root
}

func newQueryCommand() *cobra.Command {
	var l1 bool
	var l2 bool
	var outputValues []string
	cmd := &cobra.Command{
		Use:   "query [-l1|-l2] EXPR [-o FORMAT[=FILE] ...]",
		Short: "Evaluate a component set expression",
		Long: `Evaluate a component set expression.

Output formats:
  table, txt, json
  spdx-json, spdx-json@2.2, spdx-json@2.3
  cyclonedx-json, cyclonedx-json@1.5, cyclonedx-json@1.6, cyclonedx-json@1.7

Notes:
  - FORMAT=FILE writes to a file.
  - FORMAT without FILE writes to stdout.
  - SBOM formats are only supported with -l2.`,
		Example: `  sbomcalc query -l1 "a.json and b.json"
  sbomcalc query -l2 "new.json minus old.json" -o json
  sbomcalc query -l2 "new.json minus old.json" -o table -o cyclonedx-json@1.7=added.cdx.json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			level, err := cli.ParseLevel(l1, l2)
			if err != nil {
				return err
			}
			specs, err := cli.ParseOutputSpecs(outputValues)
			if err != nil {
				return err
			}
			if err := cli.ValidateQueryOutputs(level, specs); err != nil {
				return err
			}
			result, err := engine.Query(args[0], level)
			if err != nil {
				return err
			}
			return writer.WriteQuery(result, specs)
		},
	}
	cmd.Flags().BoolVar(&l1, "l1", false, "use package name as component key")
	cmd.Flags().BoolVar(&l2, "l2", false, "use package name and version as component key")
	cmd.Flags().StringArrayVarP(&outputValues, "output", "o", nil, "output format: FORMAT or FORMAT=FILE")
	return cmd
}

func newDiffCommand(changedOnly bool) *cobra.Command {
	var outputValues []string
	use := "diff old.json new.json [-o FORMAT[=FILE] ...]"
	short := "Compare two SBOM files"
	run := engine.Diff
	if changedOnly {
		use = "changed old.json new.json [-o FORMAT[=FILE] ...]"
		short = "Show changed components between two SBOM files"
		run = engine.Changed
	}
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long: short + `.

Output formats:
  table, txt, json

Notes:
  - FORMAT=FILE writes to a file.
  - FORMAT without FILE writes to stdout.
  - SBOM formats are not supported for diff or changed.`,
		Example: `  sbomcalc diff old.json new.json
  sbomcalc diff old.json new.json -o json=result.json
  sbomcalc changed old.json new.json -o txt`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			specs, err := cli.ParseOutputSpecs(outputValues)
			if err != nil {
				return err
			}
			if err := cli.ValidateDiffOutputs(specs); err != nil {
				return err
			}
			result, err := run(args[0], args[1])
			if err != nil {
				return err
			}
			return writer.WriteDiff(result, specs, changedOnly)
		},
	}
	cmd.Flags().StringArrayVarP(&outputValues, "output", "o", nil, "output format: FORMAT or FORMAT=FILE")
	return cmd
}
