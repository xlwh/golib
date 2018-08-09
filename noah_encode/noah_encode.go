/* noah_encode.go - encode struct to noah format(i.e., lines of key:value)  */
/*
modification history
--------------------
2014/9/28, by Zhang Miao, create
*/
/*
DESCRIPTION
This should only be used when struct has no member of struct, slice or map
*/
package noah_encode

import (
    "bytes"
    "fmt"
    "reflect"
)

// get attribute types of m 
func attributes(m interface{}) (map[string]reflect.Type, error) {
    typ := reflect.TypeOf(m)
    // if a pointer to a struct is passed, get the type of the dereferenced object
    if typ.Kind() == reflect.Ptr{
        typ = typ.Elem()
    }
    
    // Only structs are supported. So return err if the passed object isn't a struct
    if typ.Kind() != reflect.Struct {        
        return nil, fmt.Errorf("%v type can't have attributes inspected\n", typ.Kind())
    }

    // create an attribute data structure as a map of types keyed by a string.
    attrs := make(map[string]reflect.Type)
    
    // loop through the struct's fields and set the map
    for i := 0; i < typ.NumField(); i++ {
        p := typ.Field(i)
        if !p.Anonymous {
            attrs[p.Name] = p.Type
        }
    }
    
    return attrs, nil
}

// output noah string (lines of key:value) for data
func Encode(data interface{}) ([]byte, error) {
    return EncodeData(data, "", false)
}

/* output noah string (lines of key:value) for data
 *
 * Params:
 *     - data: data to encode
 *     - prefix: prefix of noah key
 *     - ingoreUnknown: ingore unsupported type
 *
 * Return:
 *     - noah string for data
 */
func EncodeData(data interface{}, prefix string, ingoreUnknown bool) ([]byte, error) {
	var buf bytes.Buffer
	var str string

    // set prefix 
    if prefix != "" {
        prefix = prefix + "_"
    }

    // get attributes of data
    Attrs, err := attributes(data)
    if err != nil {
        return nil, err
    }

    // iterate through the attributes of data
    value := reflect.ValueOf(data)
    for name, mtype := range Attrs {
        switch mtype.Name() {
            case "string":
                str = fmt.Sprintf("%s%s:%s\n", prefix, name, value.FieldByName(name).String())
            case "bool":
                str = fmt.Sprintf("%s%s:%t\n", prefix, name, value.FieldByName(name).Bool())
            case "int", "int8", "int16", "int32", "int64":
                str = fmt.Sprintf("%s%s:%d\n", prefix, name, value.FieldByName(name).Int())
            case "uint", "uint8", "uint16", "uint32", "uint64":
                str = fmt.Sprintf("%s%s:%d\n", prefix, name, value.FieldByName(name).Uint())
            case "float32", "float64":
                str = fmt.Sprintf("%s%s:%f\n", prefix, name, value.FieldByName(name).Float())                
            default:
                if !ingoreUnknown {
                    return nil, fmt.Errorf("unsupported type:%s", mtype.Name())
                }
                continue
        }
        buf.WriteString(str)
    }

    return buf.Bytes(), nil
}
