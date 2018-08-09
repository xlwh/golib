/* file.go - file copy implementation */
/*
modification history
--------------------
2014/9/15, by Weiwei02, create
*/
/*
DESCRIPTION
*/
package file_util

import (
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "path"
    "path/filepath"
)

// atomicCopy copy src to dst in a atomic op
// step1: copy src to a temp file
// step2: rename to dst
func AtomicCopy(src, dst string) error {
    // open src file
    srcFile, err := os.Open(src)
    if err != nil {
        return fmt.Errorf("open src error %s", err)
    }
    defer srcFile.Close()

    srcFileStat, err := srcFile.Stat()
    if err != nil {
        return fmt.Errorf("stat src error %s", err)
    }

    if !srcFileStat.Mode().IsRegular() {
        return fmt.Errorf("%s is not a regular file", src)
    }

    // open temp file
    tempFile, err := ioutil.TempFile(os.TempDir(), filepath.Base(src))
    if err != nil {
        return err
    }

    // copy to temp file
    _, err = io.Copy(tempFile, srcFile)
    if err != nil {
        tempFile.Close()
        return err
    }
    tempFile.Close()

    // mkdir all dir
    dirPath := path.Dir(dst)
    if err = os.MkdirAll(dirPath, 0755); err != nil {
        return err
    }

    // atomic rename file
    return os.Rename(tempFile.Name(), dst)
}

// regular file copy from src to dst
// return file length, error
func CopyFile(src, dst string) (int64, error) {
    srcFile, err := os.Open(src)
    if err != nil {
        return 0, fmt.Errorf("open src error %s", err)
    }
    defer srcFile.Close()

    srcFileStat, err := srcFile.Stat()
    if err != nil {
        return 0, fmt.Errorf("stat src error %s", err)
    }

    if !srcFileStat.Mode().IsRegular() {
        return 0, fmt.Errorf("%s is not a regular file", src)
    }

    // mkdir all dir
    dirPath := path.Dir(dst)
    if err = os.MkdirAll(dirPath, 0755); err != nil {
        return 0, err
    }

    dstFile, err := os.Create(dst)
    if err != nil {
        return 0, fmt.Errorf("create dst error %s", dst)
    }
    defer dstFile.Close()

    return io.Copy(dstFile, srcFile)
}

// backup file, atomic rename
func BackupFile(path string, bakPath string) error {
    copyPath := fmt.Sprintf("%s.%d.bak", path, os.Getpid())
    if _, err := CopyFile(path, copyPath); err != nil {
        return err
    }

    if err := os.Rename(copyPath, bakPath); err != nil {
        return err
    }

    return nil
}
