package goson

import (
    "encoding/json"
    "strings"
    "errors"
    "strconv"
    "go/types"
)

/**
 * Created by tuxer on 9/6/17.
 */

type JsonObject struct {
    dataMap map[string]interface{}
}

func (j *JsonObject) ToBytes() []byte {
    data, e := json.Marshal(j.dataMap)
    if e != nil {
        return nil
    }
    return data
}

func (j *JsonObject) ToString() *string {
    data := j.ToBytes()
    str := ``
    if data != nil  {
        str = string(data)
    } else {
        return nil
    }
    return &str
}

func (j *JsonObject) GetDataMap() map[string]interface{}   {
    return j.dataMap
}

func (j *JsonObject) GetJsonArray(path string) []JsonObject    {
    obj := j.get(path)

    values, ok := obj.([]interface{})

    if !ok  {
        return nil
    }
    var arrJson []JsonObject
    for _, value := range values   {
        mapValue, ok := value.(map[string]interface{})
        if ok   {
            jo := JsonObject{dataMap: mapValue}
            arrJson = append(arrJson, jo)
        }
    }
    return arrJson
}
func (j *JsonObject) GetJsonObject(path string) *JsonObject    {
    obj := j.get(path)

    v, ok := obj.(map[string]interface{})
    if ok   {
        jo := JsonObject{ dataMap: v }
        return &jo
    }
    return nil
}

func (j *JsonObject) GetFloat(path string) *float64 {
    obj := j.get(path)

    switch obj.(type) {
    case float64:
        float, _ := obj.(float64)
        return &float
    case string:
        str, _ := obj.(string)
        val, e := strconv.ParseFloat(str, 64)
        if e != nil {
            return nil
        }
        return &val
    default:
        return nil
    }
}

func (j *JsonObject) GetInt(path string) *int {
    obj := j.get(path)

    switch obj.(type) {
    case float64:
        float, _ := obj.(float64)
        val := int(float)
        return &val
    case string:
        str, _ := obj.(string)
        val, e := strconv.Atoi(str)
        if e != nil {
            return nil
        }
        return &val
    default:
        return nil
    }
}
func (j *JsonObject) GetBoolean(path string) *bool {
    obj := j.get(path)
    b, ok := obj.(bool)
    if ok   {
        return &b
    } else {
        return nil
    }
}

func (j *JsonObject) GetString(path string) *string {
    obj := j.get(path)

    switch obj.(type) {
    case string:
        str, _ := obj.(string)
        return &str
    case float64:
        float, _ := obj.(float64)
        str := strconv.FormatFloat(float, 'f', -1, 64)
        return &str
    default:
        return nil
    }
}

func (j *JsonObject) Put(path string, value interface{}) *JsonObject    {
    j.PutE(path, value)
    return j
}
func (j *JsonObject) PutE(path string, value interface{}) error   {
    ptr, ok := value.(types.Pointer)
    if ok   {
        value = ptr.Elem()
    }

    arrays, ok := value.([]JsonObject)
    if ok   {
        arrayMap := []map[string]interface{}{}
        for _, jo := range arrays {
            arrayMap = append(arrayMap, jo.dataMap)
        }
        value = arrayMap
    }
    _, ok = value.(*JsonObject)
    if ok   {
        value = value.(*JsonObject).dataMap
    }
    _, ok = value.(JsonObject)
    if ok   {
        value = value.(JsonObject).dataMap
    }

    if j.dataMap == nil  {
        j.dataMap = make(map[string]interface{})
    }

    rootMap := j.dataMap
    currentMap := rootMap

    splittedPath := strings.Split(path, `.`)
    for index, pathItem := range splittedPath   {
        if index < len(splittedPath) - 1    {
            _, ok := currentMap[pathItem]
            if !ok   {
                currentMap[pathItem] = make(map[string]interface{})
            }
            currentMap, ok = currentMap[pathItem].(map[string]interface{})
            if !ok  {
                return errors.New(pathItem + `is not a json object`)
            }
        } else {
            currentMap[pathItem] = value
        }
    }
    j.dataMap = rootMap
    return nil
}

func (j *JsonObject) get(path string) interface{} {
    splittedPath := strings.Split(path, `.`)

    var jsonMap interface{}
    jsonMap = j.dataMap
    var val interface{}
    for _, pathItem := range splittedPath   {
        if jsonMap == nil   {
            return nil
        }
        val = jsonMap.(map[string]interface{})[pathItem]

        switch val.(type) {
        case map[string]interface{}:
            jsonMap = val
        default:
            jsonMap = nil
        }
    }
    return val
}

func Parse(data []byte) *JsonObject    {
    jo := JsonObject{}
    json.Unmarshal(data, &jo.dataMap)
    return &jo
}