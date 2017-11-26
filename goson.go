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

func (j *JsonObject) ToFormattedBytes() []byte {
    data, e := json.MarshalIndent(j.dataMap, ``, `  `)
    if e != nil {
        return nil
    }
    return data
}

func (j *JsonObject) ToBytes() []byte {
    if len(j.dataMap) == 0  {
        return []byte(`{}`)
    }
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
        jo := JsonObject{dataMap: value.(map[string]interface{})}
        arrJson = append(arrJson, jo)
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
func (j *JsonObject) GetFloatD(path string, defValue float64) float64 {
    val := j.GetFloat(path)
    if val == nil   {
        return defValue
    } else {
        return *val
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
func (j *JsonObject) GetIntD(path string, defValue int) int {
    val := j.GetInt(path)
    if val == nil   {
        return defValue
    } else {
        return *val
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
func (j *JsonObject) GetBooleanD(path string, defValue bool) bool {
    val := j.GetBoolean(path)
    if val == nil   {
        return defValue
    } else {
        return *val
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
func (j *JsonObject) GetStringD(path string, defValue string) string {
    val := j.GetString(path)
    if val == nil   {
        return ``
    } else {
        return *val
    }
}

func (j *JsonObject) Put(path string, value interface{}) *JsonObject    {
    j.putE(path, value)
    return j
}
func (j *JsonObject) putE(path string, value interface{}) error   {
    ptr, ok := value.(types.Pointer)
    if ok   {
        value = ptr.Elem()
    }

    if arrays, ok := value.([]*JsonObject); ok {
        arrayMap := []map[string]interface{}{}
        for _, jo := range arrays {
            arrayMap = append(arrayMap, jo.dataMap)
        }
        value = arrayMap
    } else if arrays, ok := value.([]JsonObject); ok {
        arrayMap := []map[string]interface{}{}
        for _, jo := range arrays {
            arrayMap = append(arrayMap, jo.dataMap)
        }
        value = arrayMap
    } else if ptrJ, ok := value.(*JsonObject); ok {
        value = ptrJ.dataMap
    } else if j, ok := value.(JsonObject); ok {
        value = j.dataMap
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

    if j == nil {
        return nil
    }
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
        case []interface{}:
            return val
        default:
            jsonMap = nil
        }
    }
    return val
}

func Parse(data []byte) *JsonObject    {
    jo := JsonObject{}
    e := json.Unmarshal(data, &jo.dataMap)
    if e != nil {
        return nil
    } else {
        return &jo
    }
}