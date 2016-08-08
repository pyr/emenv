package emenv

import (
	"strings"
)

type Options struct {
	ImplicitYes bool
}

type TokenType int

const (
	StringToken TokenType = iota
	NumberToken
	OpenVectorToken
	CloseVectorToken
	OpenParToken
	CloseParToken
	SymbolToken
	KeywordToken
	DotToken
	QuoteToken
	NilToken
	EOFToken
)

type Token struct {
	Type   TokenType
	Number int
	String string
}

type Tokenizer struct {
	r *strings.Reader
}

type NodeType int

const (
	ListNode NodeType = iota
	VectorNode
	PairNode
	DotNode
	StringNode
	SymbolNode
	KeywordNode
	NumberNode
	NilNode
	EOFNode
)

type Node struct {
	Type     NodeType
	Number   int
	String   string
	Children []Node
}

type Stack struct {
	Tokens []Token
}

type Version struct {
	Members []int
	Literal string
}

type StorageType int

const (
	FileStorage StorageType = iota
	TarStorage
)

type PackageDef struct {
	Name    string
	Repo    string
	Type    PackageType
	Version Version
}

type PackageType int

const (
	RootPackage PackageType = iota
	StandardPackage
	ThemePackage
	ProvidedPackage
	ShadowPackage
	DependencyPackage
)

type Package struct {
	Name         string
	Version      Version
	Desc         string
	Type         StorageType
	URL          string
	Dependencies []PackageDef
}

type Repository struct {
	Name     string
	URL      string
	Version  int
	Packages []Package
}

type Source struct {
	Name string
	URL  string
}

type SourceConfig struct {
	Sources  map[string]Source
	Prefer   []string
	Packages []PackageDef
	Provided []string
}

type InstallDef struct {
	Type      PackageType
	URL       string
	Repo      string
	Name      string
	Version   string
	StoreType StorageType
	Depth     int
	Parent    *InstallNode
}

type InstallNode struct {
	Children []InstallNode
	Def      InstallDef
}

type InstallSet struct {
	Tree     InstallNode
	Packages map[string]InstallDef
}

type Upgrade struct {
	Prev InstallDef
	Next InstallDef
}

type DiffSet struct {
	Keep []InstallDef
	Delete []InstallDef
	Install []InstallDef
	Upgrade []Upgrade
}

type Env struct {
	BaseDir      string
	ArchiveDir   string
	PackageDir   string
	Packages     []PackageDef
	Sources      map[string]Source
	Prefer       []string
	Previous     map[string]InstallDef
	Provided     []string
	Repositories map[string]Repository
	InstallSet   InstallSet
	DiffSet      DiffSet
	Options      Options
}
