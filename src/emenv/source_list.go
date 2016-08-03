package emenv

import (
	"fmt"
	"io/ioutil"
	"os"
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

func LoadEnv(path string) (*Env, error) {

	body, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tokens, err := ParseTokens(body)
	if err != nil {
		return nil, err
	}

	sources := make(map[string]Source)

	sources["melpa-stable"] = Source{Name: "melpa-stable", URL: "https://stable.melpa.org/packages"}
	sources["melpa"] = Source{Name: "melpa", URL: "https://melpa.org/packages"}
	sources["gnu"] = Source{Name: "gnu", URL: "https://elpa.gnu.org/packages"}
	sources["org"] = Source{Name: "org", URL: "http://orgmode.org/elpa"}
	sources["sunrise"] = Source{Name: "sunrise", URL: "http://joseito.republika.pl/sunrise-commander"}

	prefer := make([]string, 5)
	prefer[0] = "melpa-stable"
	prefer[1] = "org"
	prefer[2] = "gnu"
	prefer[3] = "melpa"
	prefer[4] = "sunrise"

	provided := make([]string, 4)
	provided[0] = "emacs"
	provided[1] = "cl-lib"
	provided[2] = "eieio"
	provided[3] = "json"

	packages := make([]PackageDef, 0)

	env := Env{
		Sources:    sources,
		Prefer:     prefer,
		Provided:   provided,
		Packages:   packages,
		InstallSet: NewInstallSet(),
	}

	for {
		tree, err := ParseForm(&tokens)
		if err != nil {
			return nil, err
		}

		if tree.Type == EOFNode {
			break
		}

		if tree.Type != ListNode {
			return nil, BadSyntaxError
		}
		err = env.AddToConfig(tree.Children)
		if err != nil {
			return nil, err
		}
	}

	bdir := os.ExpandEnv("${PWD}/.emenv")
	adir := fmt.Sprintf("%s/archives/", bdir)
	pdir := fmt.Sprintf("%s/packages/", bdir)
	err = os.MkdirAll(adir, 0755)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(pdir, 0755)
	if err != nil {
		return nil, err
	}
	env.ArchiveDir = adir
	env.BaseDir = bdir
	env.PackageDir = pdir
	env.Repositories = make(map[string]Repository)
	return &env, nil
}
