/* file_test.go - unit test file for file.go*/
/*
modification history
--------------------
2015/4/9, by zhangjiyang01, create
*/
/*
DESCRIPTION
*/

package file_util

import (
    "io/ioutil"
    "os"
    "path/filepath"
    "testing"
)

func TestAtomicCopy(t *testing.T) {
    srcFilePath, _ := filepath.Abs("testdata/src_file")
    dstFilePath, _ := filepath.Abs("testdata/dst_file")
    testString := "testString"

    // set tmpDir env
    if err := os.Setenv("TMPDIR", "testdata"); err != nil {
        t.Error(err)
    }
    srcFile, err := os.Create(srcFilePath)
    if err != nil {
        t.Error(err)
    }
    if _, err := srcFile.WriteString(testString); err != nil {
        t.Error(err)
    }

    if err := AtomicCopy(srcFilePath, dstFilePath); err != nil {
        t.Error(err)
    }
    dstFile, err := os.Open(dstFilePath)
    if err != nil {
        t.Error(err)
    }

    buf, err := ioutil.ReadAll(dstFile)
    if err != nil {
        t.Error(err)
    }

    if string(buf) != testString {
        t.Error("copy failed")
    }
}
