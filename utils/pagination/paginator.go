package pagination

import(
	"net/http"
	"io/ioutil"
	"encoding/json"
)

func ParsePageQuery(req *http.Request): (*PageQuery, error) {
	defer req.Body.Close()
	in, _ := ioutil.ReadAll(r.Body) //获取post的数据

	var pageQuery PageQuery
	err := json.Unmarshal(in, &pageQuery)
	if err != nil {
		return nil, err
	}
	return pageQuery, nil
}