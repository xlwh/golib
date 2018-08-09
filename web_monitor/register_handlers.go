/* register_handlers.go - register web handlers */
/*
modification history
--------------------
2014/9/23, by Zhang Miao, create
*/
/*
DESCRIPTION
*/
package web_monitor

import (
    "errors"
    "fmt"
)

// RegisterHandlers - register handlers in handler-table to WebHandlers
//
// Params:
//      - wh    : WebHandlers
//      - hType : hanlder type, WEB_HANDLE_MONITOR or WEB_HANDLE_RELOAD
//      - ht    : handler table
//
// Returns:
//      error
func RegisterHandlers(wh *WebHandlers, hType int, ht map[string]interface{}) error {
    // check WebHandlers
    if wh == nil {
        return errors.New("nil WebHandlers")
    }
    
    // check hType
    var typeStr string
    switch hType {
        case WEB_HANDLE_MONITOR:
            typeStr = "MONITOR"
        case WEB_HANDLE_RELOAD:
            typeStr = "RELOAD"
        default:
            return fmt.Errorf("invalid handler type:%d", hType)
    }
    
    // register handlers
    for name, handler := range ht {
        err := wh.RegisterHandler(hType, name, handler) 
        if err != nil {
            return fmt.Errorf("register:%s:%s:%s", typeStr, name, err.Error())
        }
    }
    
    return nil
}