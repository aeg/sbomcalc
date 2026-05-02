package expr

import (
	"reflect"
	"testing"

	"github.com/eiji/sbomcalc/internal/model"
)

func TestParseAndEvalLeftAssociative(t *testing.T) {
	node, err := Parse("a.json and b.json minus c.json")
	if err != nil {
		t.Fatal(err)
	}
	gotFiles := Files(node)
	wantFiles := []string{"a.json", "b.json", "c.json"}
	if !reflect.DeepEqual(gotFiles, wantFiles) {
		t.Fatalf("files = %#v, want %#v", gotFiles, wantFiles)
	}

	sets := map[string]model.KeySet{
		"a.json": model.NewKeySet(model.ComponentKey{Name: "a"}, model.ComponentKey{Name: "b"}),
		"b.json": model.NewKeySet(model.ComponentKey{Name: "b"}, model.ComponentKey{Name: "c"}),
		"c.json": model.NewKeySet(model.ComponentKey{Name: "b"}),
	}
	got, err := Eval(node, func(path string) (model.KeySet, error) {
		return sets[path], nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("result = %#v, want empty set", got)
	}
}
