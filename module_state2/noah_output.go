/* noah_output.go - for noah format output	*/
/*
modification history
--------------------
2015/7/17, by Zhang Miao, create
*/
/*
DESCRIPTION

*/
package module_state2

import (
	"fmt"
	"strings"
	"unicode"
)

/*
generate key for noah output

Params:
	- key: the original key
	- noahKeyPrefix: e.g., "mod_gtc"
	- programName: e.g., "go-bfereader"
	- withProgramName: whether program name should be included in the result

Returns:
	final key, e.g., "bfe_reader_ERR_PB_SEEK", "go-bfereader.bfe_reader_ERR_PB_SEEK"
*/
func noahKeyGen(key string, noahKeyPrefix string, programName string, withProgramName bool) string {
	if programName != "" && withProgramName {
		if noahKeyPrefix == "" {
			return fmt.Sprintf("%s.%s", programName, key)
		}

		return fmt.Sprintf("%s.%s_%s", programName, noahKeyPrefix, key)
	} else {
		if noahKeyPrefix == "" {
			return key
		}

		return fmt.Sprintf("%s_%s", noahKeyPrefix, key)
	}
}

/*
argus collect plugin only support letter, num, ".", "_", - use "_" instead of unsupport character. 
more info here http://devops.baidu.com/new/argus/acquisition.md

Params:
    - originKey: original key

Returns:
    - noahKey: key for arugs
*/
func escapeNoahKey(originKey string) string {
	noahKey := originKey
	for _, str := range originKey {
		switch {
		case unicode.IsLetter(str):
		case unicode.IsNumber(str):
		case byte(str) == '-':
		case byte(str) == '_':
		case byte(str) == '.':
		default:
			noahKey = strings.Replace(noahKey, string(str), "_", -1)
		}
	}
	return noahKey
}

/*
generate and escape key for noah output

Params:
	- key: the original key
	- noahKeyPrefix: e.g., "mod_gtc"
	- programName: e.g., "go-bfereader"
	- withProgramName: whether program name should be included in the result

Returns:
	final key, e.g., "bfe_reader_ERR_PB_SEEK", "go-bfereader.bfe_reader_ERR_PB_SEEK"
*/
func NoahKeyGen(key string, noahKeyPrefix string, programName string, withProgramName bool) string {
    noahKey := noahKeyGen(key, noahKeyPrefix, programName, withProgramName)
    return escapeNoahKey(noahKey)
}
