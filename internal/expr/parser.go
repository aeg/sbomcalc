package expr

import "fmt"

type parser struct {
	tokens []token
	pos    int
}

func Parse(input string) (Node, error) {
	tokens, err := lex(input)
	if err != nil {
		return nil, err
	}
	p := parser{tokens: tokens}
	node, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.peek().kind != tokenEOF {
		return nil, fmt.Errorf("unexpected token: %s", p.peek().value)
	}
	return node, nil
}

func (p *parser) parseExpr() (Node, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	for p.peek().kind == tokenOperator {
		op := p.next().op
		right, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		left = BinaryNode{Op: op, Left: left, Right: right}
	}
	return left, nil
}

func (p *parser) parsePrimary() (Node, error) {
	switch p.peek().kind {
	case tokenFile:
		return FileNode{Path: p.next().value}, nil
	case tokenLParen:
		p.next()
		node, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if p.peek().kind != tokenRParen {
			return nil, fmt.Errorf("unmatched parenthesis")
		}
		p.next()
		return node, nil
	case tokenOperator:
		return nil, fmt.Errorf("unexpected operator: %s", p.peek().value)
	case tokenRParen:
		return nil, fmt.Errorf("unmatched parenthesis")
	default:
		return nil, fmt.Errorf("unexpected end of expression")
	}
}

func (p *parser) peek() token {
	return p.tokens[p.pos]
}

func (p *parser) next() token {
	tok := p.tokens[p.pos]
	p.pos++
	return tok
}
