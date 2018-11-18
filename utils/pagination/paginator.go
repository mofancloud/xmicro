package pagination

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

func ParsePageQueryFromRequest(req *http.Request) (*PageQuery, error) {
	defer req.Body.Close()
	in, _ := ioutil.ReadAll(req.Body) //获取post的数据

	var pageQuery PageQuery
	err := json.Unmarshal(in, &pageQuery)
	if err != nil {
		return nil, err
	}
	return &pageQuery, nil
}

func ParsePageQueryFromReader(reader io.Reader) (*PageQuery, error) {
	in, err := ioutil.ReadAll(reader) //获取post的数据
	if err != nil {
		return nil, err
	}

	var pageQuery PageQuery
	err = json.Unmarshal(in, &pageQuery)
	if err != nil {
		return nil, err
	}
	return &pageQuery, nil
}
