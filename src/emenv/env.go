package emenv

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

func (env *Env) Sync() error {

	for _, src := range env.Sources {
		fmt.Printf("syncing repository %s at %s\n", src.Name, src.URL)
		if err := env.FetchRepository(src); err != nil {
			return err
		}
	}
	return nil
}

func (env *Env) FetchInstallSet() error {

	buffer := bytes.NewBufferString(";; autoload-file for Emenv\n")
	for _, idef := range env.InstallSet.Packages {
		if idef.Type == ProvidedPackage {
			buffer.WriteString(fmt.Sprintf(";; %s is provided\n", idef.Name))
			continue
		}
		if err := env.FetchPackage(idef); err != nil {
			return err
		}
		if idef.Type == ThemePackage {
			buffer.WriteString(fmt.Sprintf("(add-to-list 'custom-theme-load-path \"%s%s-%s\")\n",
				env.PackageDir,
				idef.Name,
				idef.Version))
		} else {
			buffer.WriteString(fmt.Sprintf("(add-to-list 'load-path \"%s%s-%s\")\n",
				env.PackageDir,
				idef.Name,
				idef.Version))
		}
	}
	ioutil.WriteFile(fmt.Sprintf("%s/load.el", env.BaseDir), buffer.Bytes(), 0644)
	return nil
}

func (env *Env) Install() error {

	err := env.LoadRepositories()
	if err != nil {
		return err
	}

	if err := env.ResolveInstallSet(); err != nil {
		return err
	}

	DumpNode(env.InstallSet.Tree, 0)

	if err := env.FetchInstallSet(); err != nil {
		return err
	}
	return nil
}
