package goson

import (
    "encoding/json"
    "strings"
    "errors"
    "strconv"
)

/**
 * Created by tuxer on 9/6/17.
 */

type JsonObject struct {
    dataMap map[string]interface{}
}

func (j *JsonObject) ToString() *string {
    //buffer := ``
    //for key, value := range j.dataMap   {
    //    switch value.(type) {
    //    case string:
    //
    //    }
    //}
    data, e := json.Marshal(j.dataMap)
    if e != nil {
        return nil
    }
    str := string(data)
    return &str
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

func (j *JsonObject) GetInt(path string) (int, error) {
    obj := j.get(path)

    switch obj.(type) {
    case float64:
        float, _ := obj.(float64)
        return int(float), nil
    case string:
        str, _ := obj.(string)
        i, e := strconv.Atoi(str)
        if e != nil {
            return 0, e
        }
        return i, nil
    default:
        return 0, errors.New(`unable to get ` + path + `, is not int`)
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

func (j *JsonObject) Put(path string, value interface{}) error   {
    splittedPath := strings.Split(path, `.`)

    if j.dataMap == nil  {
        j.dataMap = make(map[string]interface{})
    }

    rootMap := j.dataMap
    currentMap := rootMap

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