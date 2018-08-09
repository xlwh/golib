/* conf_path.go - process path of config file   */
/*
modification history
--------------------
2014/8/21, by Zhang Miao, create
2014/9/29, by Zhang Miao, move from go-bfe to golang-lib
*/
/*
DESCRIPTION
*/
package conf_path

import (
    "path"
    "strings"
)

/* ConfPathProc - Process path of config file
 *
 * Params:
 *      - confPath: origin path for config file
 *      - confRoot: root path of ALL config
 * 
 * Returns:
 *      the final path of config file
 *      (1) path starts with "/", it's absolute path, return path untouched
 *      (2) else, it's relative path, return path.Join(confRoot, path)
 */
func ConfPathProc(confPath string, confRoot string) string {
	if !strings.HasPrefix(confPath, "/") {
	    // relative path to confRoot
	    confPath = path.Join(confRoot, confPath)
	}    
	
	return confPath
}
