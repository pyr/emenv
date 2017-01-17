package emenv

import (
	"os"
	"fmt"
	"io/ioutil"
)

func LoadEnv(path string, opts Options) (*Env, error) {

	body, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tokens, err := ParseTokens(body)
	if err != nil {
		return nil, err
	}

	sources := make(map[string]Source)

	sources["melpa-stable"] = Source{Name: "melpa-stable", URL: "http://stable.melpa.org/packages"}
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
	previous := make(map[string]InstallDef)

	env := Env{
		Sources:     sources,
		Prefer:      prefer,
		Provided:    provided,
		Previous:    previous,
		Packages:    packages,
		InstallSet:  NewInstallSet(),
		Options:     opts,
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

	for sname, _ := range(sources) {
		found := false
		for _, p := range(env.Prefer) {
			if p == sname {
				found = true
			}

		}
		if !found {
			env.Prefer = append(env.Prefer, sname)
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
