package reflect

import (
	"fmt"
	"reflect"
	"strings"
)

var StructInfoMap = make(map[reflect.Type]*StructInfo)

//结构体信息
type StructInfo struct {
	FieldsMap map[string]*StructField //字段字典集合
	Name      string                  //类型名
}

//结构体字段信息
type StructField struct {
	Name           string //字段名
	FieldType      reflect.Type
	Value          reflect.Value //字段值
	TableFieldName string        //表属性名
}

//获得结构体的信息
func GetStructInfo(target interface{}) (*StructInfo, error) {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("not ptr param")
	}
	t := v.Elem().Type()
	//判断target的类型
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("not struct param")
	}
	return GetReflectInfo(t, v.Elem())
}

// 获得结构体的反射的信息
func GetReflectInfo(t reflect.Type, v reflect.Value) (*StructInfo, error) {
	var structInfo *StructInfo

	fieldsMap := make(map[string]*StructField)
	// 从map里取结构体信息, 如果map没有则新建一个然后存map
	if value, ok := StructInfoMap[t]; ok {
		structInfo = value
		//更新缓存的结构体字段的值,这一部分肯定不能使用缓存,
		// 因为每一个变量的value是不同的
		for key, _ := range structInfo.FieldsMap {
			structField := structInfo.FieldsMap[key]
			if structField.Value.CanSet() && structField.Value.IsValid() {
				//更新字段的value属性
				structField.Value.Set(v.FieldByName(structField.Name))
			} else {
				return nil, fmt.Errorf("StructField [%s] is can not set or is not valid", structInfo.FieldsMap[key].Name)
			}
		}
	} else {
		fmt.Printf("will reflect\n")

		// 遍历所有属性
		for index := 0; index < t.NumField(); index++ {
			structField := t.Field(index)
			structFieldValue := v.Field(index)
			// 获取field标签的值 作为数据库字段名
			tableField := strings.TrimSpace(structField.Tag.Get("bson"))
			structFieldType := structField.Type

			// 如果字段
			if len(tableField) != 0 {
				// 构造一个新的StructField
				sf := &StructField{
					Name:           structField.Name,
					TableFieldName: tableField,
					FieldType:      structFieldType,
					Value:          structFieldValue,
				}
				// 将新的StructField放入Map
				fieldsMap[tableField] = sf
			}
		}

		//构造一个新的StructInfo
		structInfo = &StructInfo{
			Name:      t.Name(),
			FieldsMap: fieldsMap,
		}
		//将新的StructInfo放入Map当缓存用
		StructInfoMap[t] = structInfo
	}
	return structInfo, nil
}
