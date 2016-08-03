package emenv

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func FileExists(path string) bool {

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (env *Env) FetchRepository(src Source) error {

	fmt.Printf("fetching repository %s from %s\n", src.Name, src.URL)
	contents := fmt.Sprintf("%s/archive-contents", src.URL)

	resp, err := http.Get(contents)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", env.ArchiveDir, src.Name), body, 0644)

	if err != nil {
		return err
	}
	return nil
}

func (env *Env) LoadRepository(src Source) error {

	path := fmt.Sprintf("%s/%s", env.ArchiveDir, src.Name)

	if FileExists(path) == false {
		if err := env.FetchRepository(src); err != nil {
			return err
		}
	}

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
	fmt.Printf("loaded repository %s from %s\n", src.Name, path)
	repo, err := RepositoryFromAST(src.Name, src.URL, tree)
	if err != nil {
		return err
	}
	env.Repositories[repo.Name] = repo
	return nil
}

func (env *Env) LoadRepositories() error {
	for _, src := range env.Sources {
		if err := env.LoadRepository(src); err != nil {
			return err
		}
	}
	return nil
}
