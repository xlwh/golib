/* check.go - reload_trigger adaptor interface */
/*
modification history
--------------------
2014/12/22, by caodong, create
*/
/*
DESCRIPTION
*/
package conf_check

type checker interface {
    // LoadAndCheck - check config
    //
    // params:
    //      - filename  : config name
    // returns:
    //      - version   : version number when err is nil; else empty string
    //      - error     : return nil means OK
    LoadAndCheck(filename string) (string, error)
}

type Config struct {
    filename string
    check    checker
}

// create Config
func NewConfig(f string, c checker) *Config {
    config := new(Config)
    config.filename = f
    config.check = c
    return config
}

// exposed to the external interface to check config
func (c *Config) Check() (string, error) {
    return c.check.LoadAndCheck(c.filename)
}
