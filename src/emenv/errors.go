package emenv

import (
	"errors"
	"fmt"
)

func RepositoryNotFoundError(repo string) error {
	return fmt.Errorf("Repository not found: %s", repo)
}

func PackageNotFoundError(pkg string, repo string) error {
	return fmt.Errorf("Package %s not found in %s", pkg, repo)
}

func NoSuchPackageError(pkg string) error {
	return fmt.Errorf("Package %s not found in any repository", pkg)
}

var UnreachableError = errors.New("Unreachable code path")

var TrailingTokensError = errors.New("Trailing tokens")

var UnknownTokenError = errors.New("Unknown token")

var StrayVectorError = errors.New("Stray vector closing")

var StrayListError = errors.New("Stray list closing")

var DanglingVectorError = errors.New("Dangling vector closing")

var DanglingListError = errors.New("Dangling list closing")

var UnknownDirectiveError = errors.New("Unknown directive")

var BadSyntaxError = errors.New("Bad syntax")
