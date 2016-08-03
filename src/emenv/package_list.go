package emenv

import (
	"fmt"
	"strings"
)

func VersionFromAST(node Node) (Version, error) {

	members := make([]int, 0)
	strs := make([]string, 0)
	if len(node.Children) < 1 {
		return Version{}, BadSyntaxError
	}
	for _, child := range node.Children {
		if child.Type != NumberNode {
			return Version{}, BadSyntaxError
		}
		members = append(members, child.Number)
		strs = append(strs, fmt.Sprintf("%d", child.Number))
	}
	return Version{Members: members, Literal: strings.Join(strs, ".")}, nil
}

func DependencyFromAST(node Node) (PackageDef, error) {

	if node.Type != ListNode || len(node.Children) != 2 {
		return PackageDef{}, BadSyntaxError
	}
	if node.Children[0].Type != SymbolNode ||
		node.Children[1].Type != ListNode {
		return PackageDef{}, BadSyntaxError
	}
	version, err := VersionFromAST(node.Children[1])
	if err != nil {
		return PackageDef{}, err
	}
	return PackageDef{Name: node.Children[0].String, Type: DependencyPackage, Repo: "", Version: version}, nil
}

func DependenciesFromAST(node Node) ([]PackageDef, error) {

	deps := make([]PackageDef, 0)
	if node.Type == NilNode || len(node.Children) == 0 {
		return deps, nil
	}

	for _, child := range node.Children {
		dep, err := DependencyFromAST(child)
		if err != nil {
			return deps, err
		}
		deps = append(deps, dep)
	}
	return deps, nil
}

func PackageFromAST(url string, node Node) (Package, error) {
	if len(node.Children) < 3 {
		return Package{}, BadSyntaxError
	}

	if node.Children[0].Type != SymbolNode ||
		node.Children[1].Type != DotNode ||
		node.Children[2].Type != VectorNode {
		return Package{}, BadSyntaxError
	}

	details := node.Children[2].Children
	if len(details) < 4 {
		return Package{}, BadSyntaxError
	}

	if details[0].Type != ListNode ||
		(details[1].Type != ListNode &&
			details[1].Type != NilNode) ||
		details[2].Type != StringNode ||
		details[3].Type != SymbolNode {

		return Package{}, BadSyntaxError

	}

	version, err := VersionFromAST(details[0])
	if err != nil {
		return Package{}, err
	}

	deps, err := DependenciesFromAST(details[1])
	if err != nil {
		return Package{}, err
	}

	storage := FileStorage
	suffix := "el"
	switch {
	case details[3].String == "single":
		storage = FileStorage
		suffix = "el"
		break
	case details[3].String == "tar":
		storage = TarStorage
		suffix = "tar"
		break
	default:
		return Package{}, BadSyntaxError
	}

	pkg := Package{Name: node.Children[0].String,
		Type:         storage,
		Version:      version,
		Dependencies: deps,
		URL: fmt.Sprintf("%s/%s-%s.%s",
			url, node.Children[0].String,
			version.Literal, suffix),
		Desc: details[2].String}
	return pkg, nil
}

func RepositoryFromAST(name string, url string, node Node) (Repository, error) {

	if node.Type != ListNode {
		return Repository{}, BadSyntaxError
	}

	if len(node.Children) < 2 {
		return Repository{}, BadSyntaxError
	}

	if node.Children[0].Type != NumberNode {
		return Repository{}, BadSyntaxError
	}

	packages := make([]Package, 0)
	for _, child := range node.Children[1:] {
		if child.Type != ListNode {
			return Repository{}, BadSyntaxError
		}
		pkg, err := PackageFromAST(url, child)
		if err != nil {
			return Repository{}, err
		}
		packages = append(packages, pkg)
	}

	return Repository{Name: name, URL: url, Version: node.Children[0].Number, Packages: packages}, nil
}
