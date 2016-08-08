package emenv

func (env *Env) FindPackageIn(rname string, pname string) (Package, error) {

	repo, ok := env.Repositories[rname]
	if ok == false {
		return Package{}, RepositoryNotFoundError(rname)
	}

	for _, p := range repo.Packages {
		if p.Name == pname {
			return p, nil
		}
	}
	return Package{}, PackageNotFoundError(pname, rname)
}

func (env *Env) TentativelyShadow(parent *InstallNode, pname string, depth int) bool {
	previous, ok := env.InstallSet.Packages[pname]
	if ok == false {
		return false
	}
	if previous.Depth > depth {
		for i, c := range previous.Parent.Children {
			if pname == c.Def.Name {
				previous.Parent.Children[i] = previous.Parent.Children[len(previous.Parent.Children)-1]
				previous.Parent.Children = previous.Parent.Children[:len(previous.Parent.Children)-1]

				idef := InstallDef{Name: pname, Type: ShadowPackage}
				previous.Parent.Children = append(previous.Parent.Children, InstallNode{Def: idef})
				return false
			}
		}
		return false
	}
	return true
}

func (env *Env) AddPkgToInstallSet(parent *InstallNode, repo string, ptype PackageType, pkg Package, depth int) error {

	if env.TentativelyShadow(parent, pkg.Name, depth) == true {
		idef := InstallDef{
			Name: pkg.Name,
			Type: ShadowPackage,
		}
		parent.Children = append(parent.Children, InstallNode{Def: idef})
		return nil
	}

	idef := InstallDef{
		Name:      pkg.Name,
		Repo:      repo,
		StoreType: pkg.Type,
		Type:      ptype,
		Depth:     depth,
		URL:       pkg.URL,
		Parent:    parent,
		Version:   pkg.Version.Literal,
	}
	inode := InstallNode{Def: idef, Children: make([]InstallNode, 0)}
	for _, dep := range pkg.Dependencies {
		if err := env.AddToInstallSet(&inode, dep, depth+1); err != nil {
			return err
		}
	}
	parent.Children = append(parent.Children, inode)
	env.InstallSet.Packages[pkg.Name] = idef
	return nil
}

func (env *Env) AddProvidedToInstallSet(parent *InstallNode, pdef PackageDef, depth int) {

	if env.TentativelyShadow(parent, pdef.Name, depth) == true {
		idef := InstallDef{
			Name: pdef.Name,
			Type: ShadowPackage,
		}
		parent.Children = append(parent.Children, InstallNode{Def: idef})
		return
	}
	idef := InstallDef{
		Name:  pdef.Name,
		Type:  ProvidedPackage,
		Depth: 0,
	}
	env.InstallSet.Packages[pdef.Name] = idef
	parent.Children = append(parent.Children, InstallNode{Def: idef})
}

func (env *Env) AddToInstallSet(parent *InstallNode, pdef PackageDef, depth int) error {

	for _, pv := range env.Provided {
		if pv == pdef.Name {
			env.AddProvidedToInstallSet(parent, pdef, depth)
			return nil
		}
	}

	if len(pdef.Repo) > 0 {
		pkg, err := env.FindPackageIn(pdef.Repo, pdef.Name)
		if err != nil {
			return err
		}
		if err = env.AddPkgToInstallSet(parent, pdef.Repo, pdef.Type, pkg, depth); err != nil {
			return err
		}
		return nil
	}

	for _, r := range env.Prefer {
		pkg, err := env.FindPackageIn(r, pdef.Name)
		if err == nil {
			if err = env.AddPkgToInstallSet(parent, r, pdef.Type, pkg, depth); err != nil {
				return err
			}
			return nil
		}
	}

	return NoSuchPackageError(pdef.Name)
}

func (env *Env) ResolveInstallSet() error {

	for _, p := range env.Packages {
		if err := env.AddToInstallSet(&env.InstallSet.Tree, p, 0); err != nil {
			return err
		}
	}
	return nil
}

func NewInstallSet() InstallSet {
	return InstallSet{
		Tree: InstallNode{
			Def:      InstallDef{Type: RootPackage},
			Children: make([]InstallNode, 0),
		},
		Packages: make(map[string]InstallDef),
	}
}
