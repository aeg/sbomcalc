package engine

import "github.com/aeg/sbomcalc/internal/model"

func Changed(oldPath, newPath string) (model.DiffResult, error) {
	result, err := Diff(oldPath, newPath)
	if err != nil {
		return model.DiffResult{}, err
	}
	result.Added = nil
	result.Removed = nil
	result.Unchanged = nil
	return result, nil
}
