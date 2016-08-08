package emenv

import (
	"fmt"
	"io/ioutil"
)

func (env *Env) AddPackageToConfig(list []Node, ptype PackageType) error {

	if list[0].Type != SymbolNode && list[0].Type != StringNode {
		return BadSyntaxError
	}
	pdef := PackageDef{Name: list[0].String, Type: ptype}

	if len(list) > 1 {
		for _, elem := range list[1:] {
			if elem.Type != ListNode || elem.Children[0].Type != SymbolNode {
				return BadSyntaxError
			}
			arg := elem.Children[0].String
			switch {
			case arg == "repo":
				if len(elem.Children) != 2 && elem.Children[1].Type != SymbolNode {
					return BadSyntaxError
				}
				pdef.Repo = elem.Children[1].String
			default:
				return UnknownDirectiveError
			}
		}
	}
	env.Packages = append(env.Packages, pdef)
	return nil
}

func (env *Env) AddSourceToConfig(list []Node) error {

	if len(list) != 2 || list[0].Type != SymbolNode || list[1].Type != StringNode {
		return BadSyntaxError
	}
	sdef := Source{Name: list[0].String, URL: list[1].String}
	env.Sources[list[0].String] = sdef
	return nil
}

func (env *Env) SetPreferenceOrder(list []Node) error {

	prefer := make([]string, 0)

	for _, elem := range list {
		if elem.Type != SymbolNode && elem.Type != StringNode {
			return BadSyntaxError
		}
		prefer = append(prefer, elem.String)
	}
	env.Prefer = prefer
	return nil
}

func (env *Env) SetProvidedDependencies(list []Node) error {

	provided := make([]string, 0)
	for _, elem := range list {
		if elem.Type != SymbolNode && elem.Type != StringNode {
			return BadSyntaxError
		}
		provided = append(provided, elem.String)
	}
	env.Provided = provided
	return nil
}

func (env *Env) AddToConfig(list []Node) error {

	if list[0].Type != SymbolNode {
		return BadSyntaxError
	}

	switch {
	case list[0].String == "package":
		return env.AddPackageToConfig(list[1:], StandardPackage)
	case list[0].String == "theme":
		return env.AddPackageToConfig(list[1:], ThemePackage)
	case list[0].String == "source":
		return env.AddSourceToConfig(list[1:])
	case list[0].String == "prefer":
		return env.SetPreferenceOrder(list[1:])
	case list[0].String == "provided":
		return env.SetProvidedDependencies(list[1:])
	default:
		return UnknownDirectiveError
	}
	return UnreachableError
}

func FindInSet(pkgs map[string]InstallDef, id InstallDef) (*InstallDef, bool) {

	p, ok := pkgs[id.Name]
	if !ok {
		return nil, false
	}
	return &p, (p.Version == id.Version && p.Repo == id.Repo)
}

func (env *Env) LoadPreviousInstallSet() error {

	path := fmt.Sprintf("%s/plist.el", env.BaseDir)

	body, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	tokens, err := ParseTokens(body)
	if err != nil {
		return err
	}

	tree, err := ParseTree(tokens)
	if err != nil {
		return err
	}

	if tree.Type != ListNode {
		return BadSyntaxError
	}
	for _, node := range tree.Children {
		if node.Type != ListNode || len(node.Children) != 3 {
			return BadSyntaxError
		}
		if (node.Children[0].Type != SymbolNode ||
			node.Children[1].Type != StringNode ||
			node.Children[2].Type != SymbolNode) {
			return BadSyntaxError
		}
		env.Previous[node.Children[0].String] = InstallDef{
			Name: node.Children[0].String,
			Version: node.Children[1].String,
			Repo: node.Children[2].String,
		}
	}

	// Now that we have a previous installed set, compute differences

	for _, prev := range env.Previous {
		id, equal := FindInSet(env.InstallSet.Packages, prev)
		switch {
		case equal:
			env.DiffSet.Keep = append(env.DiffSet.Keep, prev)
		case id != nil:
			env.DiffSet.Upgrade = append(env.DiffSet.Upgrade, Upgrade{Prev: prev, Next: *id})
		default:
			env.DiffSet.Delete = append(env.DiffSet.Delete, prev)
		}
	}

	for _, p := range env.InstallSet.Packages {
		if (p.Type != StandardPackage && p.Type != ThemePackage && p.Type != DependencyPackage) {
			continue
		}
		id, _ := FindInSet(env.Previous, p)
		if id == nil {
			env.DiffSet.Install = append(env.DiffSet.Install, p)
		}
	}
	DumpDiffSet(env.DiffSet)
	return nil
}
