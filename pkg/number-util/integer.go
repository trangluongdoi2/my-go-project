package numberutil

import (
	"fmt"
	"strconv"
)

func ToInt64(v interface{}) (int64, error) {
	switch val := v.(type) {
	case int64:
		return val, nil
	case string:
		return strconv.ParseInt(val, 10, 64)
	default:
		return 0, fmt.Errorf("unsupported type %T", v)
	}
}
