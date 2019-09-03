package humsize

import "fmt"

const (
	B = 1 << (iota * 10)
	KB
	MB
	GB
	TB
)

func GetSize(size int64) string {
	var sizeStr string
	switch true {
	case size > TB:
		sizeStr = fmt.Sprintf("%.2f TB", float64(size)/TB)
	case size > GB:
		sizeStr = fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size > MB:
		sizeStr = fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size > KB:
		sizeStr = fmt.Sprintf("%.2f KB", float64(size)/KB)
	case size > B:
		sizeStr = fmt.Sprintf("%.2f B", float64(size)/B)
	}
	return sizeStr
}