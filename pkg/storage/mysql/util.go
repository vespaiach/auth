package mysql

import "github.com/vespaiach/auth/pkg/storage"

func getOrderDirection(i storage.Direction) string {
	switch i {
	case storage.Ascendant:
		return "ASC"
	case storage.Descendant:
		return "DESC"
	default:
		return ""
	}
}
