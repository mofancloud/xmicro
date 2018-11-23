package mongodb

import (
	"errors"
	"fmt"
	"reflect"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/mofancloud/xmicro/data"
	xreflect "github.com/mofancloud/xmicro/reflect"
	"github.com/mofancloud/xmicro/utils"
)

func All(c *mgo.Collection, m Model) *mgo.Query {
	return Where(c, nil)
}

func Find(c *mgo.Collection, m Model) *mgo.Query {
	return Where(c, m.Unique()).Limit(1)
}

func Where(c *mgo.Collection, q interface{}) *mgo.Query {
	return c.Find(q)
}

func Update(c *mgo.Collection, m Model) (*mgo.ChangeInfo, error) {
	return Find(c, m).Apply(mgo.Change{
		ReturnNew: true,
		Update: bson.M{
			"$set": m,
		},
	}, m)
}

func Insert(c *mgo.Collection, m Model) error {
	return c.Insert(m)
}

func Delete(c *mgo.Collection, m Model) error {
	return c.Remove(m.Unique())
}

func Page(c *mgo.Collection, pageQuery *data.PageQuery, m Model, list interface{}) error {
	filters, pageNo, pageSize, sorts := ParsePageQuery(m, pageQuery)

	offset := int((int32(pageNo) - 1) * pageSize)
	limit := int(pageSize)

	err := c.Find(bson.M(filters)).Skip(offset).Limit(limit).Sort(sorts...).All(&list)
	return err
}

func ensureIndexes(db *mgo.Database, m Indexed) {
	coll := db.C(m.Collection())
	for _, i := range m.Indexes() {
		coll.EnsureIndex(i)
	}
}

type DBFunc func(*mgo.Collection) error

func Execute(mongoSession *mgo.Session, databaseName string, collectionName string, fn DBFunc) error {
	session := mongoSession.Clone()
	defer session.Close()

	db := session.DB(databaseName)

	collection := db.C(collectionName)
	if collection == nil {
		err := fmt.Errorf("Collection %s does not exist", collectionName)
		return err
	}
	err := fn(collection)
	if err != nil {
		return err
	}

	return nil
}

/**
 * 根据 model的类型来验证 data的字段是否合法，并转为 bson
 */
func I2bson(m Model, data map[string]interface{}) (bson.M, error) {
	s, err := xreflect.GetStructInfo(m)

	if err != nil {
		return nil, err
	}

	r := bson.M{}

	for k, v := range data {
		f, ok := s.FieldsMap[k]

		// 如果field存在
		if !ok {
			return nil, errors.New(fmt.Sprintf("field %s not exist", k))
		}
		// TODO 字段类型匹配
		t1 := reflect.TypeOf(v)
		t2 := f.FieldType
		if t1 != t2 {
			return nil, errors.New(fmt.Sprintf("field %s type not matched, %v => %v", k, t1, t2))
		}

		r[k] = v
	}

	return r, nil
}

func BuildCriteria(m Model, filters map[string]interface{}) bson.M {
	mStruct, _ := xreflect.GetStructInfo(m)

	criteria := bson.M{}

	if filters == nil || len(filters) == 0 {
		return criteria
	}

	for k, v := range filters {
		filterType := data.FilterType(k)

		switch filterType {
		case data.FilterType_AND:
			{
				//fmt.Printf("filter value %v", v)
				//subFilters := v.([]map[string]interface{})
				subFilters := v.([]interface{})

				subCriterias := []bson.M{}
				for _, item := range subFilters {
					subFilter := item.(map[string]interface{})
					subCriterias = append(subCriterias, BuildCriteria(m, subFilter))
				}

				criteria["$and"] = subCriterias // 同一级别只允许一个
			}
		case data.FilterType_OR:
			{
				//fmt.Printf("filter value %v", v)
				//subFilters := v.([]map[string]interface{})
				subFilters := v.([]interface{})

				subCriterias := []bson.M{}
				for _, item := range subFilters {
					subFilter := item.(map[string]interface{})
					subCriterias = append(subCriterias, BuildCriteria(m, subFilter))
				}

				criteria["$or"] = subCriterias // 同一级别只允许一个
			}
		case data.FilterType_NOR:
			{
				subFilters := v.([]map[string]interface{})

				subCriterias := []bson.M{}
				for _, subFilter := range subFilters {
					subFilter = subFilter
					subCriterias = append(subCriterias, BuildCriteria(m, subFilter))
				}

				criteria["$nor"] = subCriterias // 同一级别只允许一个
			}
		default:
			{
				vMap := v.(map[string]interface{})

				subCriteria := bson.M{}

				for vKey, vValue := range vMap {
					filterType = data.FilterType(vKey)
					// time类型转换
					if mStruct != nil {
						fieldInfo, ok := mStruct.FieldsMap[k]
						if ok {
							// 如果字段是时间类型
							if fieldInfo.FieldType == xreflect.TimeType {
								// 并且传入的值是int64, 则把int64转成 time格式先
								fieldType := fieldInfo.FieldType

								if fieldType == xreflect.TimeType {
									// 并且传入的值是int64, 则把int64转成 time格式先
									if v, ok := vValue.(int64); ok {
										vValue = utils.Unix(v, 0)

									} else if v, ok := vValue.(int); ok {
										vValue = utils.Unix(int64(v), 0)
									} else if v, ok := vValue.(float64); ok {
										vValue = utils.Unix(int64(v), 0)
									} else if v, ok := vValue.(int32); ok {
										vValue = utils.Unix(int64(v), 0)
									}
								}
							}
						}
					}

					switch filterType {
					case data.FilterType_EQ:
						{
							subCriteria["$eq"] = vValue
						}
					case data.FilterType_GT:
						{
							subCriteria["$gt"] = vValue
						}
					case data.FilterType_GTE:
						{
							subCriteria["$gte"] = vValue
						}
					case data.FilterType_LT:
						{
							subCriteria["$lt"] = vValue
						}
					case data.FilterType_LTE:
						{
							subCriteria["$lte"] = vValue
						}
					case data.FilterType_LIKE:
						{
							subCriteria["$regex"] = vValue
						}
					case data.FilterType_NE:
						{
							subCriteria["$ne"] = vValue
						}
					case data.FilterType_NOT_IN:
						{
							subCriteria["$nin"] = vValue
						}
					case data.FilterType_IN:
						{
							subCriteria["$in"] = vValue
						}
					}
				}
				criteria[k] = subCriteria // 这里是否有必要
			}
		}
	}
	return criteria
}

func ParsePageQuery(m Model, pageQuery *data.PageQuery) (criteria bson.M, pageNo int64, pageSize int32, sorts []string) {
	// 构造 filterMap排重, 每个 property 都构造一个 子filter数组, 应对 <, > 各种情况
	criteria = BuildCriteria(m, pageQuery.Filters)

	pageNo = pageQuery.PageNo
	if pageQuery.PageNo < 1 {
		pageNo = 1
	}

	pageSize = pageQuery.PageSize
	if pageSize < 1 {
		pageSize = 20
	} else if pageSize > 1000 {
		pageSize = 1000
	}

	sorts = []string{}

	if pageQuery.Sort != nil {
		for _, s := range pageQuery.Sort {
			var s1 string
			switch s.Type {
			case data.SortType_ASC:
				{
					s1 = s.Property
				}
			case data.SortType_DSC:
				{
					s1 = fmt.Sprintf("-%s", s.Property)
				}
			default:
				{
					s1 = s.Property
				}
			}

			sorts = append(sorts, s1)
		}
	}

	return
}
