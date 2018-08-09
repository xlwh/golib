/* ip_loc_table.go - test file of ipdict */
/*
modification history
--------------------
2016/7/13, by Jiang Hui, create
*/
/*
DESCRIPTION
*/

package loc_load

import (
	"testing"
)

//test NewIpLocDictFile
func TestNewIpLocDictFile(t *testing.T) {
	//test max
	_, err := NewIpLocDictFile("testdata/iplocation.txt", 1000001)
	if err == nil {
		t.Errorf("should return error but return nil")
		return
	}

	//test min
	_, err = NewIpLocDictFile("testdata/iplocation.txt", 0)
	if err == nil {
		t.Errorf("should return error but return nil")
		return
	}
}

// test ErrMaxLineExceed
func TestIpLocLoadMaxLine(t *testing.T) {
	locLoadFile, err := NewIpLocDictFile("testdata/iplocation.txt", 2)
	if err != nil {
		t.Errorf("should return nil but return error[%s]", err.Error())
		return
	}

	//load ip location file
	_, err = locLoadFile.CheckAndLoad("")
	if err != ErrMaxLineExceed {
		t.Errorf("should return ErrMaxLineExceed but return error[%s]", err.Error())
		return
	}
}

//test version is new ,no need update
func TestNoUpdate(t *testing.T) {
	locLoadFile, err := NewIpLocDictFile("testdata/iplocation.txt", 1000000)
	if err != nil {
		t.Errorf("should return nil but return error[%s]", err.Error())
		return
	}

	//load ip location file
	_, err = locLoadFile.CheckAndLoad("")
	if err != nil {
		t.Errorf("should return nil but return error[%s]", err.Error())
		return
	}

	version := locLoadFile.version
	_, err = locLoadFile.CheckAndLoad(version)
	if err != ErrNoNeedUpdate {
		t.Errorf("should return ErrNoNeedUpdate but return error[%s]", err.Error())
		return
	}

	version = version + "1"
	_, err = locLoadFile.CheckAndLoad(version)
	if err != nil {
		t.Errorf("should return nil but return error[%s]", err.Error())
		return
	}
}
