package structs

import (
	"os"
	"strings"
)

type FileInfo struct {
	Path string
	Hash string
}

func (fi *FileInfo) SplitPath() []string {
	return strings.Split(fi.Path, string(os.PathSeparator))
}
