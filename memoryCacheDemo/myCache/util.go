package myCache

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const (
	B = 1 << (10 * iota)
	KB
	MB
	GB
	TB
	PB
)

func parseSize(size string) (int64, string) {
	re, _ := regexp.Compile("[0-9]+")
	unit := re.ReplaceAllString(size, "")
	nums, err := strconv.ParseInt(strings.Replace(size, unit, "", 1), 10, 64)
	if err != nil {
		log.Println("parse size err, invalue size input:", err)
		//默认缓存大小
		nums = 100
		unit = "MB"
	}
	unit = strings.ToUpper(unit)
	var memSize int64
	switch unit {
	case "B":
		memSize = nums * B
	case "KB":
		memSize = nums * KB
	case "MB":
		memSize = nums * MB
	case "GB":
		memSize = nums * GB
	case "TB":
		memSize = nums * TB
	case "PB":
		memSize = nums * PB
	default:
		log.Println("invalid unit:", unit)
		nums = 100
		unit = "MB"
		memSize = nums * MB
	}
	size = strconv.Itoa(int(nums)) + unit
	return memSize, size
}

func CalSize(data interface{}) int64 {
	/*通过reflect计算interface指向值的大小
	val := reflect.ValueOf(data)
	size := int64(0)

	switch val.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		size = int64(val.Len()) * int64(val.Type().Elem().Size())
	case reflect.Struct:
		numFields := val.NumField()
		for i := 0; i < numFields; i++ {
			//递归调用计算大小
			//fieldSize := int64(val.Field(i).Type().Size())
			fieldSize := int64(val.Field(i).Type().Size())
			size += fieldSize
		}
	default:
		size = int64(val.Type().Size())
	}
	return size
	*/

	//通过序列化计算interface指向值的大小
	bytes, _ := json.Marshal(data)
	size := len(bytes)
	fmt.Println("the size of value interface is:", size)
	return int64(size)
}
