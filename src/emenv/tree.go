package emenv

func NewStack(tokens []Token) Stack {
	return Stack{Tokens: tokens}
}

func (stack *Stack) Parse() (Node, error) {

	head := stack.Tokens[0]
	stack.Tokens = stack.Tokens[1:]

	switch {
	case head.Type == NumberToken:
		return Node{Type: NumberNode, Number: head.Number}, nil
	case head.Type == NilToken:
		return Node{Type: NilNode}, nil
	case head.Type == StringToken:
		return Node{Type: StringNode, String: head.String}, nil
	case head.Type == DotToken:
		return Node{Type: DotNode}, nil
	case head.Type == KeywordToken:
		return Node{Type: KeywordNode, String: head.String}, nil
	case head.Type == SymbolToken:
		return Node{Type: SymbolNode, String: head.String}, nil
	case head.Type == OpenVectorToken:
		node := Node{Type: VectorNode}
		for {
			subhead := stack.Tokens[0]
			if subhead.Type == CloseVectorToken {
				stack.Tokens = stack.Tokens[1:]
				break
			}
			if len(stack.Tokens) == 0 {
				return Node{}, DanglingVectorError
			}
			subnode, err := stack.Parse()

			if err != nil {
				return Node{}, err
			}
			node.Children = append(node.Children, subnode)
		}
		return node, nil
	case head.Type == OpenParToken:
		node := Node{Type: ListNode}
		for {
			subhead := stack.Tokens[0]
			if subhead.Type == CloseParToken {
				stack.Tokens = stack.Tokens[1:]
				break
			}
			if len(stack.Tokens) == 0 {
				return Node{}, DanglingListError
			}
			subnode, err := stack.Parse()
			if err != nil {
				return Node{}, err
			}
			node.Children = append(node.Children, subnode)
		}
		return node, nil
	case head.Type == CloseVectorToken:
		return Node{}, StrayVectorError
	case head.Type == CloseParToken:
		return Node{}, StrayListError

	default:
		return Node{}, UnknownTokenError
	}
	return Node{}, UnreachableError
}

func ParseTree(tokens []Token) (Node, error) {
	stack := NewStack(tokens)

	node, err := stack.Parse()
	if err != nil {
		return Node{}, err
	}

	if len(stack.Tokens) > 1 {
		return Node{}, TrailingTokensError
	}
	if len(stack.Tokens) == 1 && stack.Tokens[0].Type != EOFToken {
		return Node{}, TrailingTokensError
	}
	return node, nil
}

func ParseForm(tokens *[]Token) (Node, error) {

	stack := NewStack(*tokens)
	if len(stack.Tokens) < 1 {
		return Node{Type: EOFNode}, nil
	}
	if len(stack.Tokens) == 1 && stack.Tokens[0].Type == EOFToken {
		return Node{Type: EOFNode}, nil
	}

	node, err := stack.Parse()
	if err != nil {
		return Node{}, err
	}
	*tokens = stack.Tokens
	return node, nil
}
