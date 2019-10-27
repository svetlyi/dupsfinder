package dups

import (
	"fmt"
	"github.com/svetlyi/dupsfinder/app"
	"github.com/svetlyi/dupsfinder/config"
	"github.com/svetlyi/dupsfinder/structs"
)

// Get returns duplicates for path.
// Pages start with 0
func Get(path string, page int, app *app.App) (structs.DupsTmplObj, error) {
	selectDupsStmt := GetSelectDupsByDir(app.DB)

	rows, err := selectDupsStmt.Query(path, config.DupsPerPage, config.DupsPerPage*page)
	if nil != err {
		app.Fatal(fmt.Sprintf("An error while querying: %v", err))
		return structs.DupsTmplObj{}, err
	}
	defer rows.Close()

	var searchDupsObj structs.DupsTmplObj
	searchDupsObj.Files = make(map[string][]structs.FileTmplObj)

	for rows.Next() {
		var fileInfo structs.FileInfo
		err := rows.Scan(&fileInfo.Path, &fileInfo.Hash)
		if err != nil {
			app.Fatal(fmt.Sprintf("An error while scanning: %v", err))
			return structs.DupsTmplObj{}, err
		}

		searchDupsObj.Files[fileInfo.Hash] = append(
			searchDupsObj.Files[fileInfo.Hash],
			structs.FileTmplObj{
				Path:      fileInfo.Path,
				PathParts: fileInfo.SplitPath(),
				Hash:      fileInfo.Hash,
			},
		)
	}
	err = rows.Err()
	if err != nil {
		app.Fatal(fmt.Sprintf("An unexpected error: %v", err))
		return structs.DupsTmplObj{}, err
	}

	return searchDupsObj, nil
}
