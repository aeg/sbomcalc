package cli

import (
	"fmt"

	"github.com/aeg/sbomcalc/internal/model"
)

func ParseLevel(l1 bool, l2 bool) (model.Level, error) {
	if l1 && l2 {
		return model.Level2, fmt.Errorf("-l1 and -l2 cannot be used together")
	}
	if l1 {
		return model.Level1, nil
	}
	return model.Level2, nil
}

func ValidateQueryOutputs(level model.Level, specs []OutputSpec) error {
	for _, spec := range specs {
		if level == model.Level1 && spec.IsSBOM() {
			return fmt.Errorf("SBOM output is not supported for query -l1")
		}
	}
	return nil
}

func ValidateDiffOutputs(specs []OutputSpec) error {
	for _, spec := range specs {
		if spec.IsSBOM() {
			return fmt.Errorf("SBOM output is not supported for diff or changed")
		}
	}
	return nil
}
