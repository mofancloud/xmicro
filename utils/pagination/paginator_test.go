package pagination

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

func TestPageQuery(t *testing.T) {
	req, err := http.NewRequest("POST", "/users/page", bytes.NewBufferString(`{"pageNo":1, "pageSize":100, "filters":{}}`))
	if err != nil {
		t.Error("NewRequest err")
	}

	pageQuery, err := ParsePageQueryFromRequest(req)
	if err != nil {
		t.Error("ParsePageQueryFromRequest err")
	}

	fmt.Printf("pageQuery: %v", pageQuery)
}
