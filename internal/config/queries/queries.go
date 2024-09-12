package queries

import "strings"

const (
	BaseSelect      = "SELECT {columns} FROM {table}"
	BaseSelectWhere = "SELECT {columns} FROM {table} WHERE {conditions}"
	BaseInsert      = "INSERT INTO {table} ({columns}) VALUES ({values})"
	BaseUpdate      = "UPDATE {table} SET {assignments} WHERE {conditions}"
)

func QueryBuilder(query string, params map[string]string) string {
	for placeholder, value := range params {
		query = strings.ReplaceAll(query, "{"+placeholder+"}", value)
	}
	return query
}
