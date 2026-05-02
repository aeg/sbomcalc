package cli

import (
	"fmt"
	"strings"
)

type OutputKind int

const (
	OutputTable OutputKind = iota
	OutputText
	OutputJSON
	OutputSPDXJSON
	OutputCycloneDXJSON
)

type OutputSpec struct {
	Kind    OutputKind
	Format  string
	Version string
	File    string
}

func ParseOutputSpecs(values []string) ([]OutputSpec, error) {
	if len(values) == 0 {
		return []OutputSpec{{Kind: OutputTable, Format: "table"}}, nil
	}

	specs := make([]OutputSpec, 0, len(values))
	stdoutCount := 0
	files := map[string]struct{}{}
	for _, value := range values {
		spec, err := ParseOutputSpec(value)
		if err != nil {
			return nil, err
		}
		if spec.File == "" {
			stdoutCount++
		} else {
			if _, ok := files[spec.File]; ok {
				return nil, fmt.Errorf("same output file specified multiple times: %s", spec.File)
			}
			files[spec.File] = struct{}{}
		}
		specs = append(specs, spec)
	}
	if stdoutCount > 1 {
		return nil, fmt.Errorf("multiple stdout outputs are not allowed")
	}
	return specs, nil
}

func ParseOutputSpec(value string) (OutputSpec, error) {
	if strings.TrimSpace(value) == "" {
		return OutputSpec{}, fmt.Errorf("empty output format")
	}

	formatPart, file, _ := strings.Cut(value, "=")
	name, version, hasVersion := strings.Cut(formatPart, "@")
	spec := OutputSpec{Format: name, File: file}
	if hasVersion {
		spec.Version = version
	}

	switch name {
	case "table":
		spec.Kind = OutputTable
		if hasVersion {
			return OutputSpec{}, fmt.Errorf("table does not support version: %s", value)
		}
	case "txt":
		spec.Kind = OutputText
		if hasVersion {
			return OutputSpec{}, fmt.Errorf("txt does not support version: %s", value)
		}
	case "json":
		spec.Kind = OutputJSON
		if hasVersion {
			return OutputSpec{}, fmt.Errorf("json does not support version: %s", value)
		}
	case "spdx-json":
		spec.Kind = OutputSPDXJSON
		if spec.Version == "" {
			spec.Version = "2.3"
		}
		if spec.Version != "2.2" && spec.Version != "2.3" {
			return OutputSpec{}, fmt.Errorf("unsupported SPDX JSON version: %s", spec.Version)
		}
	case "cyclonedx-json":
		spec.Kind = OutputCycloneDXJSON
		if spec.Version == "" {
			spec.Version = "1.7"
		}
		if spec.Version != "1.5" && spec.Version != "1.6" && spec.Version != "1.7" {
			return OutputSpec{}, fmt.Errorf("unsupported CycloneDX JSON version: %s", spec.Version)
		}
	default:
		return OutputSpec{}, fmt.Errorf("unsupported output format: %s", name)
	}
	return spec, nil
}

func (s OutputSpec) IsSBOM() bool {
	return s.Kind == OutputSPDXJSON || s.Kind == OutputCycloneDXJSON
}
