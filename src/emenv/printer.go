package emenv

import (
	"fmt"
)

func DumpToken(token Token) {
	switch {
	case token.Type == EOFToken:
		fmt.Println("eof: <EOF>")
		break
	case token.Type == NilToken:
		fmt.Println("nil: <nil>")
		break
	case token.Type == NumberToken:
		fmt.Printf("num: %d\n", token.Number)
		break
	case token.Type == KeywordToken:
		fmt.Printf("kwd: :%s\n", token.String)
		break
	case token.Type == SymbolToken:
		fmt.Printf("sym: %s\n", token.String)
		break
	case token.Type == StringToken:
		fmt.Printf("str: \"%s\"\n", token.String)
		break
	case token.Type == OpenParToken:
		fmt.Println("par: (")
		break
	case token.Type == CloseParToken:
		fmt.Println("par: )")
		break
	case token.Type == OpenVectorToken:
		fmt.Println("vec: [")
		break
	case token.Type == CloseVectorToken:
		fmt.Println("vec: ]")
		break
	case token.Type == DotToken:
		fmt.Println("dot: .")
		break
	default:
		fmt.Println("err:")
		break
	}
}

func DumpTokens(tokens []Token) {
	for _, token := range tokens {
		DumpToken(token)
	}
}

func DumpTree(node Node) {
	switch {
	case node.Type == StringNode:
		fmt.Printf("str: \"%s\"\n", node.String)
		break
	case node.Type == KeywordNode:
		fmt.Printf("kwd: :%s\n", node.String)
		break
	case node.Type == SymbolNode:
		fmt.Printf("sym: %s\n", node.String)
		break
	case node.Type == NumberNode:
		fmt.Printf("num: %d\n", node.Number)
		break
	case node.Type == NilNode:
		fmt.Printf("nil: nil\n")
		break
	case node.Type == DotNode:
		fmt.Printf("dot: .\n")
		break
	case node.Type == PairNode:
		fmt.Printf("par: 0\n")
		DumpTree(node.Children[0])
		fmt.Printf("par: 1\n")
		DumpTree(node.Children[1])
		break
	case node.Type == ListNode:
		fmt.Printf("lst: (\n")
		for _, child := range node.Children {
			DumpTree(child)
		}
		fmt.Printf("lst: )\n")
		break
	case node.Type == VectorNode:
		fmt.Printf("vec: [\n")
		for _, child := range node.Children {
			DumpTree(child)
		}
		fmt.Printf("vec: ]\n")
		break
	}
}

func DumpNode(node InstallNode, depth int) {

	if node.Def.Type == ShadowPackage {
		return
	}
	for i := 0; i < depth; i++ {
		fmt.Printf("  ")
	}
	switch {
	case node.Def.Type == RootPackage:
		fmt.Printf("root\n")
	case node.Def.Type == StandardPackage:
		fmt.Printf("package %s %s from %s\n",
			node.Def.Name, node.Def.Version, node.Def.Repo)
	case node.Def.Type == DependencyPackage:
		fmt.Printf("dependency %s %s from %s\n",
			node.Def.Name, node.Def.Version, node.Def.Repo)
	case node.Def.Type == ThemePackage:
		fmt.Printf("theme %s %s from %s\n",
			node.Def.Name, node.Def.Version, node.Def.Repo)
	case node.Def.Type == ShadowPackage:
		fmt.Printf("shadowed %s\n", node.Def.Name)
	case node.Def.Type == ProvidedPackage:
		fmt.Printf("provided %s\n", node.Def.Name)
	default:
		fmt.Printf("WHAT?\n")
	}
	for _, n := range node.Children {
		DumpNode(n, depth+1)
	}
}

func DumpDefs(defs map[string]InstallDef) {

	for _, d := range defs {
		switch {
		case d.Type == RootPackage:
			fmt.Printf("root\n")
		case d.Type == StandardPackage:
			fmt.Printf("package %s %s from %s depth %d\n",
				d.Name, d.Version, d.Repo, d.Depth)
		case d.Type == DependencyPackage:
			fmt.Printf("dependency %s %s from %s depth %d\n",
				d.Name, d.Version, d.Repo, d.Depth)
		case d.Type == ShadowPackage:
			fmt.Printf("shadowed %s\n", d.Name)
		case d.Type == ThemePackage:
			fmt.Printf("theme %s %s from %s depth %d\n",
				d.Name, d.Version, d.Repo, d.Depth)
		case d.Type == ProvidedPackage:
			fmt.Printf("provided %s\n", d.Name)
		default:
			fmt.Printf("WHAT?\n")
		}
	}
}
