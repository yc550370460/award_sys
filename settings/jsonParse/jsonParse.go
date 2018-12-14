package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//定义配置文件解析后的结构
type Config struct {
	Static  string
	Template string
}

func JsonLoad(path string) Config{
	JsonParse := NewJsonStruct()
	v := Config{}

	JsonParse.Load(path, &v)
	return v
}

type JsonStruct struct {
}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}

func (jst *JsonStruct) Load(filename string, v interface{}) {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	fmt.Println(filename)
	if err != nil {
		return
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		return
	}
}