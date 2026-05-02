package engine

import (
	"sort"

	"github.com/eiji/sbomcalc/internal/model"
	"github.com/eiji/sbomcalc/internal/reader"
)

func Diff(oldPath, newPath string) (model.DiffResult, error) {
	oldVersions, err := scanVersionsByName(oldPath)
	if err != nil {
		return model.DiffResult{}, err
	}
	newVersions, err := scanVersionsByName(newPath)
	if err != nil {
		return model.DiffResult{}, err
	}

	result := model.DiffResult{From: oldPath, To: newPath}
	names := sortedNameUnion(oldVersions, newVersions)
	for _, name := range names {
		oldSet, inOld := oldVersions[name]
		newSet, inNew := newVersions[name]
		switch {
		case !inOld && inNew:
			result.Added = append(result.Added, model.VersionedName{Name: name, Versions: sortedStrings(newSet)})
		case inOld && !inNew:
			result.Removed = append(result.Removed, model.VersionedName{Name: name, Versions: sortedStrings(oldSet)})
		case equalStringSet(oldSet, newSet):
			result.Unchanged = append(result.Unchanged, model.VersionedName{Name: name, Versions: sortedStrings(oldSet)})
		default:
			result.Changed = append(result.Changed, model.ChangedName{
				Name:       name,
				OldVersion: sortedStrings(oldSet),
				NewVersion: sortedStrings(newSet),
			})
		}
	}
	return result, nil
}

func scanVersionsByName(path string) (map[string]map[string]struct{}, error) {
	out := map[string]map[string]struct{}{}
	err := reader.ScanFile(path, func(record model.ComponentRecord) error {
		name := model.NormalizeName(record.Name)
		if name == "" {
			return nil
		}
		if _, ok := out[name]; !ok {
			out[name] = map[string]struct{}{}
		}
		out[name][model.NormalizeVersion(record.Version)] = struct{}{}
		return nil
	})
	return out, err
}

func sortedNameUnion(a, b map[string]map[string]struct{}) []string {
	seen := map[string]struct{}{}
	for name := range a {
		seen[name] = struct{}{}
	}
	for name := range b {
		seen[name] = struct{}{}
	}
	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func sortedStrings(set map[string]struct{}) []string {
	values := make([]string, 0, len(set))
	for value := range set {
		values = append(values, value)
	}
	sort.Strings(values)
	return values
}

func equalStringSet(a, b map[string]struct{}) bool {
	if len(a) != len(b) {
		return false
	}
	for value := range a {
		if _, ok := b[value]; !ok {
			return false
		}
	}
	return true
}
