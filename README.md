# goson
Java-like json parser for Go.

Goson created to ease json parsing and building without hasle using struct.

# Install

```
go get -u github.com/firabliz/goson

# Example

```go
//parsing json
jsonVal := []byte(`{"string_data": "string_value", "numeric_data": 1}`)

jsObj   := goson.Parse(jsonVal)

fmt.println(jsObj.GetString(`key`))         //print: string_value

fmt.println(jsObj.GetInt(`numeric_data`))   //print: 1

```