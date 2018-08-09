/* web_handler.go - web handler framework   */
/*
modification history
--------------------
2014/7/8, by Zhang Miao, create
2014/8/7, by Zhang Miao, copy from go-bfe
2014/9/1, by Sijie YANG, reload handler support args
*/
/*
DESCRIPTION
*/
package web_monitor

import (
	"fmt"
	"net/url"
)

// type of web handler
const (
	WEB_HANDLE_MONITOR = 0 // handler for monitor
	WEB_HANDLE_RELOAD  = 1 // handler for reload
)

var handlerTypeNames = map[int]string{
	0: "monitor",
	1: "reload",
}

type WebHandlerMap map[string]interface{}

type WebHandlers struct {
	Handlers map[int]*WebHandlerMap
}

// create new WebHandlerMap
func NewWebHandlerMap() *WebHandlerMap {
	whm := make(WebHandlerMap)
	return &whm
}

// create new WebHandlers
func NewWebHandlers() *WebHandlers {
	// create bfeCallbacks
	wh := new(WebHandlers)
	wh.Handlers = make(map[int]*WebHandlerMap)

	// handlers for monitor
	wh.Handlers[WEB_HANDLE_MONITOR] = NewWebHandlerMap()
	// handlers for reload
	wh.Handlers[WEB_HANDLE_RELOAD] = NewWebHandlerMap()

	return wh
}

func (wh *WebHandlers) validateHandler(hType int, f interface{}) error {
	var err error
	switch hType {
	case WEB_HANDLE_MONITOR:
		switch f.(type) {
		case func() ([]byte, error):
		case func(map[string][]string) ([]byte, error):
		case func(url.Values) ([]byte, error):
		default:
			err = fmt.Errorf("invalid monitor handler type %T", f)
		}

	case WEB_HANDLE_RELOAD:
		switch f.(type) {
		case func() error:
		case func(map[string][]string) error:
		case func(url.Values) error:
		case func(url.Values) (string, error):
		default:
			err = fmt.Errorf("invalid reload handler type %T", f)
		}

	default:
		err = fmt.Errorf("invalid handler type[%d]", hType)
	}
	return err
}

// add filter to given callback point
func (wh *WebHandlers) RegisterHandler(hType int, command string, f interface{}) error {
	var ok bool
	var hm *WebHandlerMap

	// check format of f
	if err := wh.validateHandler(hType, f); err != nil {
		return err
	}

	// get WebHandlerMap for given hType
	hm, ok = wh.Handlers[hType]
	if !ok {
		return fmt.Errorf("invalid handler type[%d]", hType)
	}

	// handler exist already?
	_, ok = (*hm)[command]
	if ok {
		return fmt.Errorf("handler exist already, type[%s], command[%s]",
			handlerTypeNames[hType], command)
	}

	// add to WebHandlerMap
	(*hm)[command] = f

	return nil
}

// get handler list for given callback point
func (wh *WebHandlers) GetHandler(hType int, command string) (interface{}, error) {
	var ok bool
	var hm *WebHandlerMap
	var h interface{}

	// get WebHandlerMap for given hType
	hm, ok = wh.Handlers[hType]
	if !ok {
		return nil, fmt.Errorf("invalid handler type[%d]", hType)
	}

	// handler exist already?
	h, ok = (*hm)[command]
	if !ok {
		return nil, fmt.Errorf("handler not exist, type[%s], command[%s]",
			handlerTypeNames[hType], command)
	}

	return h, nil
}
