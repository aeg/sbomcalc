package engine

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/aeg/sbomcalc/internal/model"
)

func TestQueryL1Intersection(t *testing.T) {
	oldPath := filepath.Join("..", "..", "testdata", "old.spdx.json")
	newPath := filepath.Join("..", "..", "testdata", "new.cdx.json")
	result, err := Query(oldPath+" and "+newPath, model.Level1)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, component := range result.Components {
		names = append(names, component.Name)
	}
	want := []string{"curl", "openssl"}
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("names = %#v, want %#v", names, want)
	}
}

func TestQueryComplexL1Intersection(t *testing.T) {
	oldPath := filepath.Join("..", "..", "testdata", "complex-old.spdx.json")
	newPath := filepath.Join("..", "..", "testdata", "complex-new.cdx.json")
	result, err := Query(oldPath+" and "+newPath, model.Level1)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, component := range result.Components {
		names = append(names, component.Name)
	}
	want := []string{"curl", "empty-version", "openssl", "zlib"}
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("names = %#v, want %#v", names, want)
	}
}

func TestQueryComplexL2Intersection(t *testing.T) {
	oldPath := filepath.Join("..", "..", "testdata", "complex-old.spdx.json")
	newPath := filepath.Join("..", "..", "testdata", "complex-new.cdx.json")
	result, err := Query(oldPath+" and "+newPath, model.Level2)
	if err != nil {
		t.Fatal(err)
	}
	got := componentPairs(result.Components)
	want := []string{"curl@8.0.0", "openssl@3.0.0", "zlib@1.2.11"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("components = %#v, want %#v", got, want)
	}
}

func TestQueryComplexL2Minus(t *testing.T) {
	oldPath := filepath.Join("..", "..", "testdata", "complex-old.spdx.json")
	newPath := filepath.Join("..", "..", "testdata", "complex-new.cdx.json")
	result, err := Query(newPath+" minus "+oldPath, model.Level2)
	if err != nil {
		t.Fatal(err)
	}
	got := componentPairs(result.Components)
	want := []string{"curl@8.1.0", "empty-version@1.0.0", "nginx@1.24.0", "openssl@3.2.0"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("components = %#v, want %#v", got, want)
	}
}

func TestQueryComplexL2Xor(t *testing.T) {
	oldPath := filepath.Join("..", "..", "testdata", "complex-old.spdx.json")
	newPath := filepath.Join("..", "..", "testdata", "complex-new.cdx.json")
	result, err := Query(oldPath+" xor "+newPath, model.Level2)
	if err != nil {
		t.Fatal(err)
	}
	got := componentPairs(result.Components)
	want := []string{
		"curl@7.81.0",
		"curl@8.1.0",
		"empty-version@",
		"empty-version@1.0.0",
		"log4j@2.14.1",
		"nginx@1.24.0",
		"openssl@1.1.1",
		"openssl@3.2.0",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("components = %#v, want %#v", got, want)
	}
}

func TestDiff(t *testing.T) {
	oldPath := filepath.Join("..", "..", "testdata", "old.spdx.json")
	newPath := filepath.Join("..", "..", "testdata", "new.cdx.json")
	result, err := Diff(oldPath, newPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Added) != 1 || result.Added[0].Name != "nginx" {
		t.Fatalf("added = %#v", result.Added)
	}
	if len(result.Removed) != 1 || result.Removed[0].Name != "log4j" {
		t.Fatalf("removed = %#v", result.Removed)
	}
	if len(result.Changed) != 1 || result.Changed[0].Name != "openssl" {
		t.Fatalf("changed = %#v", result.Changed)
	}
	if len(result.Unchanged) != 1 || result.Unchanged[0].Name != "curl" {
		t.Fatalf("unchanged = %#v", result.Unchanged)
	}
}

func TestDiffComplexVersionSets(t *testing.T) {
	oldPath := filepath.Join("..", "..", "testdata", "complex-old.spdx.json")
	newPath := filepath.Join("..", "..", "testdata", "complex-new.cdx.json")
	result, err := Diff(oldPath, newPath)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := result.Added, []model.VersionedName{{Name: "nginx", Versions: []string{"1.24.0"}}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("added = %#v, want %#v", got, want)
	}
	if got, want := result.Removed, []model.VersionedName{{Name: "log4j", Versions: []string{"2.14.1"}}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("removed = %#v, want %#v", got, want)
	}
	if got, want := result.Unchanged, []model.VersionedName{{Name: "zlib", Versions: []string{"1.2.11"}}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("unchanged = %#v, want %#v", got, want)
	}

	wantChanged := []model.ChangedName{
		{Name: "curl", OldVersion: []string{"7.81.0", "8.0.0"}, NewVersion: []string{"8.0.0", "8.1.0"}},
		{Name: "empty-version", OldVersion: []string{""}, NewVersion: []string{"1.0.0"}},
		{Name: "openssl", OldVersion: []string{"1.1.1", "3.0.0"}, NewVersion: []string{"3.0.0", "3.2.0"}},
	}
	if !reflect.DeepEqual(result.Changed, wantChanged) {
		t.Fatalf("changed = %#v, want %#v", result.Changed, wantChanged)
	}
}

func componentPairs(components []model.ComponentRecord) []string {
	pairs := make([]string, 0, len(components))
	for _, component := range components {
		pairs = append(pairs, component.Name+"@"+component.Version)
	}
	return pairs
}
