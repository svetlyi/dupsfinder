package dups

import (
	"fmt"
	"github.com/svetlyi/dupsfinder/humsize"
	"time"
)

type Stats struct {
	ElapsedTime time.Duration
	FilesSize   int64
	FilesAmount int
}

func (s Stats) String() string {
	return fmt.Sprintf(
		"%d files of %s size were compared in %s",
		s.FilesAmount,
		humsize.GetSize(s.FilesSize),
		s.ElapsedTime.String(),
	)
}

type FileInfo struct {
	Path string
	Hash string
}
