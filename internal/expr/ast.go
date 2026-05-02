package expr

type Node interface {
	node()
}

type FileNode struct {
	Path string
}

func (FileNode) node() {}

type BinaryNode struct {
	Op          Operator
	Left, Right Node
}

func (BinaryNode) node() {}

type Operator int

const (
	OpAnd Operator = iota
	OpOr
	OpMinus
	OpXor
)
