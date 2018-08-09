/* error.go - error code for dict-ip-location-file */
/*
modification history
--------------------
2016/6/6, by Jiang Hui, create
*/
/*
DESCRIPTION


*/

package loc_load

import (
	"errors"
)

var (
	// file version not change, needn't load the file
	ErrNoNeedUpdate = errors.New("Version no change no need update")
	// line num of file larger than maxline configured
	ErrMaxLineExceed = errors.New("Max line exceed")
	// wrong meta info
	ErrWrongMetaInfo = errors.New("Wrong meta info")
)
