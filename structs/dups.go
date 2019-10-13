package structs

type FileTmplObj struct {
	Path      string
	PathParts []string
	Hash      string
}

type DupsTmplObj struct {
	Files   map[string][]FileTmplObj
	PageNum int
}
