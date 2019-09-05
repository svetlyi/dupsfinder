package dups

import (
	"fmt"
	"github.com/svetlyi/dupsfinder/humsize"
	"time"
)

type Stats struct {
	StartTime   time.Time
	EndTime     time.Time
	FilesSize   int64
	FilesAmount int
}

func (s Stats) String() string {
	var stringStat string

	if s.EndTime.IsZero() {
		stringStat = s.getCurrentStat()
	} else {
		stringStat = s.getFinalStat()
	}

	return stringStat
}

func (s Stats) getCurrentStat() string {
	return fmt.Sprintf(
		"%d files of %s size were compared in %s. The process has not finished yet",
		s.FilesAmount,
		humsize.GetSize(s.FilesSize),
		time.Since(s.StartTime),
	)
}

func (s Stats) getFinalStat() string {
	return fmt.Sprintf(
		"%d files of %s size were compared in %s",
		s.FilesAmount,
		humsize.GetSize(s.FilesSize),
		s.EndTime.Sub(s.StartTime),
	)
}

type FileInfo struct {
	Path string
	Hash string
}
