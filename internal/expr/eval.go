package expr

import (
	"fmt"

	"github.com/aeg/sbomcalc/internal/model"
)

type KeySetProvider func(path string) (model.KeySet, error)

func Files(node Node) []string {
	seen := map[string]struct{}{}
	var files []string
	var walk func(Node)
	walk = func(n Node) {
		switch v := n.(type) {
		case FileNode:
			if _, ok := seen[v.Path]; !ok {
				seen[v.Path] = struct{}{}
				files = append(files, v.Path)
			}
		case BinaryNode:
			walk(v.Left)
			walk(v.Right)
		}
	}
	walk(node)
	return files
}

func Eval(node Node, provider KeySetProvider) (model.KeySet, error) {
	switch v := node.(type) {
	case FileNode:
		return provider(v.Path)
	case BinaryNode:
		left, err := Eval(v.Left, provider)
		if err != nil {
			return nil, err
		}
		right, err := Eval(v.Right, provider)
		if err != nil {
			return nil, err
		}
		switch v.Op {
		case OpAnd:
			return left.And(right), nil
		case OpOr:
			return left.Or(right), nil
		case OpMinus:
			return left.Minus(right), nil
		case OpXor:
			return left.Xor(right), nil
		default:
			return nil, fmt.Errorf("unknown operator")
		}
	default:
		return nil, fmt.Errorf("unknown expression node")
	}
}
