/* dir.go - directory operations */
/*
modification history
--------------------
2016/9/15, by Zhang Miao, create
*/
/*
DESCRIPTION
*/

package file_util

import (
	"fmt"
	"os"
	"path/filepath"
)

// dirCreate(): check and create dir if nonexist
func DirCreate(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		/* create directory */
		err = os.MkdirAll(dir, 0744)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
 generate full file path

 Return:
 	(string in format "rootDir/day/prefix_dayTime.suffix")
*/
func FullPathGen(rootDir, day, dayTime, prefix, suffix string) string {
	dirName := filepath.Join(rootDir, day)
	fileName := fmt.Sprintf("%s_%s.%s", prefix, dayTime, suffix)
	return filepath.Join(dirName, fileName)
}
