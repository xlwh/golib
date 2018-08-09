/* pdb.go - get noah_id to pdb_id table from rms interface */
/*
modification history
--------------------
2016/7/22, by Zhang Jiyang, create
*/
/*
DESCRIPTION

Format of response is:
{
    "status": 0,
    "data": {
        "noah": {
            "0": 0,
            "100024421": 114787220,
            ...
        },
        "rms": {
            "100024421": 114787220,
            ...
        }
}
For now, we only need to get noahID to PdbID
*/

package baidu_rms

import (
	"encoding/json"
	"fmt"
)

import (
	"www.baidu.com/golang-lib/http_util"
)

var noahIDToPdbIDAddr = "http://rms.baidu.com/?r=pdbInterface/GetAllRelationsForTcoByLevel&version=V2&level=4"

// getNoahID2PdbIDTable get noahID to PdbID table by given url
// Params:
//  - url: rms interface url
// Returns:
//  - (noah id => Pdb ID , err) 
func getNoahID2PdbIDTable(url string) (map[string]int, error) {
	// request rms
	resp, err := http_util.Read(url, TIME_OUT, nil)
	if err != nil {
		return nil, fmt.Errorf("http_util.Read(): %s", err.Error())
	}

	// decode response
	result := struct {
		Status int
		Data   struct {
			Noah map[string]int //noahID(string) -> pdbID(int)
		}
	}{}
	if err = json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("json.Unmarshal(): %s", err.Error())
	}

	// check status
	if result.Status != 0 {
		return nil, fmt.Errorf("resp.Status not 0")
	}

	return result.Data.Noah, nil
}

// GetNoahID2PdbIDTable get noah id to pdb id map
func GetNoahID2PdbIDTable() (map[string]int, error) {
	return getNoahID2PdbIDTable(noahIDToPdbIDAddr)
}
