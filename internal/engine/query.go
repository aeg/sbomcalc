package engine

import (
	"github.com/eiji/sbomcalc/internal/expr"
	"github.com/eiji/sbomcalc/internal/model"
	"github.com/eiji/sbomcalc/internal/reader"
)

func Query(expression string, level model.Level) (model.QueryResult, error) {
	ast, err := expr.Parse(expression)
	if err != nil {
		return model.QueryResult{}, err
	}

	keySets := map[string]model.KeySet{}
	for _, path := range expr.Files(ast) {
		set, err := reader.ScanKeySet(path, level)
		if err != nil {
			return model.QueryResult{}, err
		}
		keySets[path] = set
	}

	resultKeys, err := expr.Eval(ast, func(path string) (model.KeySet, error) {
		return keySets[path], nil
	})
	if err != nil {
		return model.QueryResult{}, err
	}

	records, err := collectResultRecords(expr.Files(ast), resultKeys, level)
	if err != nil {
		return model.QueryResult{}, err
	}
	model.SortRecords(records, level)
	return model.QueryResult{Level: level, Components: records}, nil
}

func collectResultRecords(paths []string, keys model.KeySet, level model.Level) ([]model.ComponentRecord, error) {
	collected := make(map[model.ComponentKey]model.ComponentRecord, len(keys))
	for _, path := range paths {
		if err := reader.ScanFile(path, func(record model.ComponentRecord) error {
			key := model.KeyFor(record, level)
			if !keys.Has(key) {
				return nil
			}
			if _, ok := collected[key]; ok {
				return nil
			}
			if level == model.Level1 {
				record = model.ComponentRecord{Name: key.Name, Source: record.Source}
			}
			collected[key] = record
			return nil
		}); err != nil {
			return nil, err
		}
	}

	records := make([]model.ComponentRecord, 0, len(collected))
	for _, record := range collected {
		records = append(records, record)
	}
	return records, nil
}
