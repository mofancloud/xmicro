package pagination

import "testing"

func TestPageQuery(t *testing.T) {
	pageQuery, err := ParsePageQuery(req)
	if err != nil {
		t.Error("ParsePageQuery err")
	}

	t.Logf("pageQuery: %v", pageQuery)
}
