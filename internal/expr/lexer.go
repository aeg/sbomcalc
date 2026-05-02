package expr

import (
	"fmt"
	"strings"
	"unicode"
)

type tokenKind int

const (
	tokenEOF tokenKind = iota
	tokenFile
	tokenOperator
	tokenLParen
	tokenRParen
)

type token struct {
	kind  tokenKind
	value string
	op    Operator
}

func lex(input string) ([]token, error) {
	if strings.TrimSpace(input) == "" {
		return nil, fmt.Errorf("query expression is empty")
	}

	var tokens []token
	for i := 0; i < len(input); {
		r := rune(input[i])
		if unicode.IsSpace(r) {
			i++
			continue
		}
		switch input[i] {
		case '(':
			tokens = append(tokens, token{kind: tokenLParen, value: "("})
			i++
			continue
		case ')':
			tokens = append(tokens, token{kind: tokenRParen, value: ")"})
			i++
			continue
		case '"', '\'':
			return nil, fmt.Errorf("quoted file paths are not supported")
		}

		start := i
		for i < len(input) && !unicode.IsSpace(rune(input[i])) && input[i] != '(' && input[i] != ')' {
			i++
		}
		value := input[start:i]
		switch value {
		case "and":
			tokens = append(tokens, token{kind: tokenOperator, value: value, op: OpAnd})
		case "or":
			tokens = append(tokens, token{kind: tokenOperator, value: value, op: OpOr})
		case "minus":
			tokens = append(tokens, token{kind: tokenOperator, value: value, op: OpMinus})
		case "xor":
			tokens = append(tokens, token{kind: tokenOperator, value: value, op: OpXor})
		default:
			if strings.ContainsAny(value, "&|^") {
				return nil, fmt.Errorf("symbol operators are not supported: %s", value)
			}
			tokens = append(tokens, token{kind: tokenFile, value: value})
		}
	}
	tokens = append(tokens, token{kind: tokenEOF})
	return tokens, nil
}
