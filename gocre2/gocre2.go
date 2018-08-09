/* gocre2.go - go wrapper for calling lib of cre2 */
/*
modification history
--------------------
2014/3/13, by Taochunhua, create
*/
/*
DESCRIPTION

Setup:
    1, Gcc4 is required. Setup the PATH and LD_LIBRARY_PATH for gcc4
    2, Setup GOPATH
    3, Replace "#cgo LDFLAGS: -L/your_dev_machine_gopath" for real path
    4, cd /path/of/gocre2_go_file
    5, tar zxvf lib.tar.gz
    6, Execute command : go build
    7, Execute command : go install
    
Usage:
    import "www.baidu.com/golang-lib/gocre2"
    
    // a regex pattern
    pattern := "^\\w+@\\w+\\.\\w{2,4}$"
    
    // create a regex object
    re := gocre2.NewRegex(pattern)
    
    // match a string
    match := gocre2.Match(re, "anny@baidu.com")
    
    if match {
       // do some thing
    }
    
    // release object
    gocre2.DeleteRegex(re)
*/

package gocre2 

/*
#cgo CFLAGS: -I./include
#cgo LDFLAGS: -L/your_dev_machine_gopath/gocre2/lib/cre2 -lcre2
#cgo LDFLAGS: -L/your_dev_machine_gopath/gocre2/lib/re2 -lre2
#cgo LDFLAGS: -lstdc++ -lm
#include <cre2.h>
*/
import "C"
import "unsafe"
import "errors"

/* Create a new regex object */
func NewRegex(pattern string) (unsafe.Pointer, error) {
    opt := C.cre2_opt_new()
    if opt == nil {
        return nil, errors.New("cre2_opt_new() failed.")
    }
    
    re := C.cre2_new(C.CString(pattern), C.int(len(pattern)), opt)
    if re == nil {
        return nil, errors.New("cre2_new() failed.")
    }
    
    C.cre2_opt_delete(opt)
    
    return re, nil
}

/* Delete a regex object */
func DeleteRegex(re unsafe.Pointer) {
    C.cre2_delete(re)
}

/* Match test function */
func Match(re unsafe.Pointer, target unsafe.Pointer, length int) bool {
    res := C.cre2_match(re,
                        (*C.char)(target), //pointer of target string
                        C.int(length),     //len of string
                        0,                 //match from beginning(position = 0)
                        C.int(length),     //match length
                        C.CRE2_UNANCHORED, //default
                        nil,               //no return string
                        C.int(0))          //no return string
    if res == 1 {
        return true
    }
    return false
}
