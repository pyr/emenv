package emenv

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

func (env *Env) Sync() error {

	for _, src := range(env.Sources) {
		fmt.Printf("syncing repository %s at %s\n", src.Name, src.URL)
		if err := env.FetchRepository(src); err != nil {
			return err
		}
	}
	return nil
}

func (env *Env) FetchInstallSet() error {


	for _, idef := range(env.InstallSet.Packages) {
		if idef.Type == ProvidedPackage {
			continue
		}
		if err := env.FetchPackage(idef); err != nil {
			return err
		}
	}
	return nil
}

func (env *Env) WritePackageList() error {
	loadbuf := bytes.NewBufferString(";; autoload-file for Emenv\n")
	pkgbuf  := bytes.NewBufferString(";; package list file for Emenv\n(\n")
	for _, idef := range env.InstallSet.Packages {
		if idef.Type == ProvidedPackage {
			loadbuf.WriteString(fmt.Sprintf(";; %s is provided\n", idef.Name))
			continue
		}
		if idef.Type == ThemePackage {
			loadbuf.WriteString(fmt.Sprintf("(add-to-list 'custom-theme-load-path \"%s%s-%s\")\n",
				env.PackageDir,
				idef.Name,
				idef.Version))
		}
		loadbuf.WriteString(fmt.Sprintf("(add-to-list 'load-path \"%s%s-%s\")\n",
			env.PackageDir,
			idef.Name,
			idef.Version))
		pkgbuf.WriteString(fmt.Sprintf("(%s \"%s\" %s)\n", idef.Name, idef.Version, idef.Repo))
	}
	pkgbuf.WriteString(")\n")
	ioutil.WriteFile(fmt.Sprintf("%s/load.el", env.BaseDir), loadbuf.Bytes(), 0644)
	ioutil.WriteFile(fmt.Sprintf("%s/plist.el", env.BaseDir), pkgbuf.Bytes(), 0644)
	return nil
}

func (env *Env) DeletePackage(p InstallDef) error {
	fmt.Printf("Deleting: %s\n", p.Name)
	return os.RemoveAll(fmt.Sprintf("%s/%s-%s", env.PackageDir, p.Name, p.Version))
}

func (env *Env) InstallPackage(p InstallDef) error {
	return env.FetchPackage(p)
}

func (env *Env) UpgradePackage(up Upgrade) error {
	if err := env.DeletePackage(up.Prev); err != nil {
		return err
	}
	return env.InstallPackage(up.Next)
}

func (env *Env) ApplyDiffSet() error {
	for _, p := range(env.DiffSet.Delete) {
		if err:= env.DeletePackage(p); err != nil {
			return err
		}
	}
	for _, u := range(env.DiffSet.Upgrade) {
		if err := env.UpgradePackage(u); err != nil {
			return err
		}
	}
	for _, p := range(env.DiffSet.Install) {
		if err := env.InstallPackage(p); err != nil {
			return err
		}
	}
	return nil
}

func (env *Env) NoOpDiffSet() bool {
	return (len(env.DiffSet.Install) == 0 &&
		len(env.DiffSet.Upgrade) == 0 &&
		len(env.DiffSet.Delete) == 0)

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

	if err := env.LoadPreviousInstallSet(); err == nil {
		if env.NoOpDiffSet() {
			fmt.Println("nothing to do, bye.")
			return nil
		}
		if !(env.Options.ImplicitYes || Confirm()) {
			return nil
		}

		if err = env.ApplyDiffSet(); err != nil {
			return err
		}
	} else {

		if !(env.Options.ImplicitYes || Confirm()) {
			return nil
		}

		if err := env.FetchInstallSet(); err != nil {
			return err
		}
	}
	return env.WritePackageList()
}
