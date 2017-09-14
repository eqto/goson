package marshaller

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
    parsed      map[string]interface{}
}


func (j *JsonObject) Parse(data []byte)    {
    json.Unmarshal(data, &j.parsed)
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
            jo := JsonObject{parsed: mapValue}
            arrJson = append(arrJson, jo)
        }
    }
    return arrJson
}
func (j *JsonObject) GetJsonObject(path string) *JsonObject    {
    obj := j.get(path)

    v, ok := obj.(map[string]interface{})
    if ok   {
        jo := JsonObject{ parsed: v }
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
func (j *JsonObject) GetString(path string) (string, error) {
    obj := j.get(path)

    switch obj.(type) {
    case string:
        str, _ := obj.(string)
        return str, nil
    case float64:
        float, _ := obj.(float64)
        str := strconv.FormatFloat(float, 'f', -1, 64)
        return str, nil
    default:
        return ``, errors.New(`unable to get ` + path + `, is not string`)
    }

}

func (j *JsonObject) get(path string) interface{} {
    splittedPath := strings.Split(path, `.`)

    var jsonMap interface{}
    jsonMap = j.parsed
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
