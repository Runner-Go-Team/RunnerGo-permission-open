package public

import (
	"os"
	"reflect"
	"regexp"
)

// GetStringNum 获取字符串字符个数
func GetStringNum(stringData string) int {
	num := 0
	for range stringData {
		num++
	}
	return num
}

// CheckStructIsEmpty 判断结构体是否为空
func CheckStructIsEmpty(obj interface{}) bool {
	// 获取结构体的反射值
	value := reflect.ValueOf(obj)
	// 获取结构体的反射类型
	typ := value.Type()

	// 如果传入的不是结构体类型，则认为不为空
	if typ.Kind() != reflect.Struct {
		return false
	}

	// 遍历结构体的每个字段，判断是否有值
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		// 如果有任何一个字段有值，就认为结构体不为空
		if !field.IsZero() {
			return false
		}
	}

	return true
}

// SliceDiff 两个切片的差集
func SliceDiff[T any](a []T, b []T) []T {
	var c []T
	temp := map[any]struct{}{}

	for _, val := range b {
		if _, ok := temp[val]; !ok {
			temp[val] = struct{}{}
		}
	}
	for _, val := range a {
		if _, ok := temp[val]; !ok {
			c = append(c, val)
		}
	}

	return c
}

// SliceUnique 切片去重通过map键的唯一性去重
func SliceUnique[T any](s []T) []T {
	result := make([]T, 0, len(s))

	m := map[any]struct{}{}
	for _, v := range s {
		if _, ok := m[v]; !ok {
			result = append(result, v)
			m[v] = struct{}{}
		}
	}

	return result
}

// ContainsStringSlice 切片中是否包含改元素
func ContainsStringSlice(s []string, elem string) bool {
	for _, a := range s {
		if a == elem {
			return true
		}
	}
	return false
}

// IsEmailValid checks if the email provided is valid by regex.
func IsEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

// StructToMap struct转map 返回的map键为struct的成员名
func StructToMap(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}

	return data
}

// StructToMapJson struct转map 返回的map键为struct的json键名
func StructToMapJson(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		jsonKey := t.Field(i).Tag.Get("json")
		if jsonKey != "-" {
			data[jsonKey] = v.Field(i).Interface()
		}
	}

	return data
}

// CheckFileIsExist 检查目录是否存在
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
