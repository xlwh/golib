/* txt_load_test.go - test file of txt_load.go

/*
modification history
--------------------
2014/7/14, by Li Bingyi, create
*/
/*
DESCRIPTION
*/

package txt_load

import (
	"bytes"
	"net"
	"testing"
)

import (
	"www.baidu.com/golang-lib/ipdict"
	"www.baidu.com/golang-lib/net_util"
)

// test for normal line
func TestCheckSplit_Case0(t *testing.T) {
	var startIP, endIP net.IP
	var line string
	var err error

	line = "1.1.1.1 2.2.2.2"
	startIP, endIP, err = checkSplit(line, " ")
	if !bytes.Equal(startIP, net_util.ParseIPv4("1.1.1.1")) ||
		!bytes.Equal(endIP, net_util.ParseIPv4("2.2.2.2")) ||
		err != nil {
		t.Error("TestCheckSplit():", err)
	}

	line = "1.1.1.1"
	startIP, endIP, err = checkSplit(line, " ")
	if !bytes.Equal(startIP, net_util.ParseIPv4("1.1.1.1")) ||
		!bytes.Equal(endIP, net_util.ParseIPv4("1.1.1.1")) ||
		err != nil {
		t.Error("TestCheckSplit():", err)
	}

	line = "1.1.1.1  2.2.2.2"
	startIP, endIP, err = checkSplit(line, " ")
	if !bytes.Equal(startIP, net_util.ParseIPv4("1.1.1.1")) ||
		!bytes.Equal(endIP, net_util.ParseIPv4("2.2.2.2")) ||
		err != nil {
		t.Error("TestCheckSplit():", err)
	}

	line = "1.1.1.1  \t\t  2.2.2.2"
	startIP, endIP, err = checkSplit(line, " ")
	if !bytes.Equal(startIP, net_util.ParseIPv4("1.1.1.1")) ||
		!bytes.Equal(endIP, net_util.ParseIPv4("2.2.2.2")) ||
		err != nil {
		t.Error("TestCheckSplit():", err)
	}
}

// test for abnormal line
func TestCheckSplit_Case1(t *testing.T) {
	var line string
	var err error

	line = "1.1.1.1 a"
	_, _, err = checkSplit(line, " ")
	if err == nil {
		t.Errorf("TestCheckSplit(): line %s err", line)
	}

	line = "a 1.1.1.1"
	_, _, err = checkSplit(line, " ")
	if err == nil {
		t.Errorf("TestCheckSplit(): line %s err", line)
	}

	line = "a b"
	_, _, err = checkSplit(line, " ")
	if err == nil {
		t.Errorf("TestCheckSplit(): line %s err", line)
	}

	line = "1.1.1.1 2.2.2.2 a"
	_, _, err = checkSplit(line, " ")
	if err == nil {
		t.Errorf("TestCheckSplit(): line %s err", line)
	}
}

// test for normal line
func TestCheckLine_Case0(t *testing.T) {
	var startIP, endIP net.IP
	var line string
	var err error

	line = "1.1.1.1 2.2.2.2"
	startIP, endIP, err = checkLine(line)
	if !bytes.Equal(startIP, net_util.ParseIPv4("1.1.1.1")) ||
		!bytes.Equal(endIP, net_util.ParseIPv4("2.2.2.2")) ||
		err != nil {
		t.Error("TestCheckLine():", err)
	}

	line = "1.1.1.1"
	startIP, endIP, err = checkLine(line)
	if !bytes.Equal(startIP, net_util.ParseIPv4("1.1.1.1")) ||
		!bytes.Equal(endIP, net_util.ParseIPv4("1.1.1.1")) ||
		err != nil {
		t.Error("TestCheckLine():", err)
	}

	line = "1.1.1.1  2.2.2.2"
	startIP, endIP, err = checkLine(line)
	if !bytes.Equal(startIP, net_util.ParseIPv4("1.1.1.1")) ||
		!bytes.Equal(endIP, net_util.ParseIPv4("2.2.2.2")) ||
		err != nil {
		t.Error("TestCheckLine():", err)
	}

	line = "1.1.1.1  \t\t  2.2.2.2"
	startIP, endIP, err = checkLine(line)
	if !bytes.Equal(startIP, net_util.ParseIPv4("1.1.1.1")) ||
		!bytes.Equal(endIP, net_util.ParseIPv4("2.2.2.2")) ||
		err != nil {
		t.Error("TestCheckLine():", err)
	}

	line = "1.1.1.1\t\t   \t\t2.2.2.2"
	startIP, endIP, err = checkLine(line)
	if !bytes.Equal(startIP, net_util.ParseIPv4("1.1.1.1")) ||
		!bytes.Equal(endIP, net_util.ParseIPv4("2.2.2.2")) ||
		err != nil {
		t.Error("TestCheckLine():", err)
	}
}

func TestNewTxtFileLoader(t *testing.T) {
	fileName := "./testdata/ipdict.conf"

	f := NewTxtFileLoader(fileName)

	if f.fileName != fileName {
		t.Errorf("TestNewTxtFileLoader(): fileName %s != %s", f.fileName, fileName)
	}
}

// file not exist case
func TestLoad_Case0(t *testing.T) {
	fileName := "./testdata/no_exist.conf"
	fileLoader := NewTxtFileLoader(fileName)
	_, err := fileLoader.CheckAndLoad("")

	if err == nil {
		t.Errorf("TestCheckAndLoad(): err is not nill")
	}
}

// ip format error case
func TestLoad_Case1(t *testing.T) {
	fileName := "./testdata/ipdict.conf1"
	fileLoader := NewTxtFileLoader(fileName)
	_, err := fileLoader.CheckAndLoad("")

	if err == nil {
		t.Errorf("TestCheckAndLoad(): err is not nill")
	}
}

// startIP > endIP error case
func TestLoad_Case2(t *testing.T) {
	fileName := "./testdata/ipdict.conf1"
	fileLoader := NewTxtFileLoader(fileName)
	_, err := fileLoader.CheckAndLoad("")

	if err == nil {
		t.Errorf("TestCheckAndLoad(): err is not nill")
	}
}

// line format error case
func TestLoad_Case3(t *testing.T) {
	fileName := "./testdata/ipdict.conf3"
	fileLoader := NewTxtFileLoader(fileName)
	_, err := fileLoader.CheckAndLoad("")

	if err == nil {
		t.Errorf("TestCheckAndLoad(): err is not nill")
	}
}

// normal case
func TestLoad_Case4(t *testing.T) {
	fileName := "./testdata/ipdict.conf4"
	fileLoader := NewTxtFileLoader(fileName)
	_, err := fileLoader.CheckAndLoad("")

	if err != nil {
		t.Errorf("TestCheckAndLoad(): err is %s", err.Error())
	}
}

func TestLoad_Case5(t *testing.T) {
	table := ipdict.NewIPTable()

	fileName := "./testdata/innocent.dict"
	fileLoader := NewTxtFileLoader(fileName)
	ipItems, err := fileLoader.CheckAndLoad("")
	if err != nil {
		t.Errorf("TestCheckAndLoad(): err is not nill, %s", err.Error())
	}
	table.Update(ipItems)
	if !table.Search(net_util.ParseIPv4("1.1.1.1")) {
		t.Error("TestCheckAndLoad(); 1.1.1.1 not in table")
	}

	if !table.Search(net_util.ParseIPv4("220.181.165.194")) {
		t.Error("TestCheckAndLoad(); 220.181.165.194 not in table")
	}

}

func TestLoad_Case6(t *testing.T) {
	table := ipdict.NewIPTable()

	fileName := "./testdata/ipdict.conf5"
	fileLoader := NewTxtFileLoader(fileName)
	ipItems, err := fileLoader.CheckAndLoad("")
	if err != nil {
		t.Errorf("TestCheckAndLoad(): err is not nill %s", err.Error())
	}

	if ipItems.Length() != 5 || ipItems.Version != "" {
		t.Errorf("ipItems.Length should be 5 version should be nil")
	}

	table.Update(ipItems)
	if table.Search(net_util.ParseIPv4("1.1.1.1")) {
		t.Error("TestCheckAndLoad(); 1.1.1.1 not in table")
	}

	if !table.Search(net_util.ParseIPv4("220.181.165.194")) {
		t.Error("TestCheckAndLoad(); 220.181.165.194 not in table")
	}

}
