package data

type FilterType string

const (
	FilterType_EQ       FilterType = "EQ"       //相等
	FilterType_NE       FilterType = "NE"       //不相等
	FilterType_GT       FilterType = "GT"       //大于
	FilterType_GTE      FilterType = "GTE"      //大于等于
	FilterType_LT       FilterType = "LT"       //小于
	FilterType_LTE      FilterType = "LTE"      //小于等于
	FilterType_IN       FilterType = "IN"       //在什么范围内
	FilterType_NOT_IN   FilterType = "NOT_IN"   //不在什么范围内
	FilterType_LIKE     FilterType = "LIKE"     //like
	FilterType_NOT_LIKE FilterType = "NOT_LIKE" //not like
	FilterType_MATCH    FilterType = "MATCH"    //匹配
	FilterType_AND      FilterType = "AND"      //AND
	FilterType_OR       FilterType = "OR"       //OR
	FilterType_NOR      FilterType = "NOR"      //NOR
)

type TimeType int // 数据库的时间类型
const (
	TIME      TimeType = 1 // 时间类型 time.Time
	TIMESTAMP TimeType = 2 // 时间戳 int64
)

type SortType string

const (
	SortType_DEFAULT SortType = "DEFAULT"
	SortType_ASC     SortType = "ASC" // 升序
	SortType_DSC     SortType = "DSC" // 降序
)

type SortSpec struct {
	Property   string   `json:"property"`   // 属性名
	Type       SortType `json:"type"`       // 排序类型
	IgnoreCase bool     `json:"ignoreCase"` // 忽略大小写
}

/**
 * 新的 PageQuery
 */
type PageQuery struct {
	Filters  map[string]interface{} `json:"filters"`
	PageNo   int64                  `json:"pageNo"`
	PageSize int32                  `json:"pageSize"`
	Sort     []*SortSpec            `json:"sort"`
}
