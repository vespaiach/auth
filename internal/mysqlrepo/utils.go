package mysqlrepo

import (
	"fmt"
	"strings"

	"github.com/vespaiach/auth/internal/comtype"
)

func sqlWhereBuilder(join string, m map[string]interface{}) string {
	lenmap := len(m)
	if lenmap == 0 {
		return ""
	}

	st := make([]string, 0, lenmap)
	for key, val := range m {
		switch key {
		case "from_date":
			st = append(st, "created_at >= :from_date")
			continue
		case "to_date":
			st = append(st, "created_at <= :to_date")
			continue
		}

		switch val.(type) {
		case string:
			st = append(st, fmt.Sprintf("%s LIKE :%s", key, key))
		default:
			st = append(st, fmt.Sprintf("%s = :%s", key, key))
		}
	}

	return "WHERE " + strings.Join(st, join)
}

func sqlSortingBuilder(m map[string]comtype.SortDirection) string {
	lenmap := len(m)
	if lenmap == 0 {
		return "created_at DESC"
	}

	st := make([]string, 0, lenmap)
	for key, val := range m {
		if val == comtype.Ascending {
			st = append(st, fmt.Sprintf("%s ASC", key))
		} else {
			st = append(st, fmt.Sprintf("%s DESC", key))
		}
	}

	return strings.Join(st, ", ")
}

func sqlLikeConditionFilter(m map[string]interface{}) map[string]interface{} {
	lenmap := len(m)
	if lenmap == 0 {
		return m
	}

	for key, val := range m {
		switch val.(type) {
		case string:
			m[key] = "%" + m[key].(string) + "%"
		}
	}

	return m
}
