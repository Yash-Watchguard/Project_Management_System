package config

import (
	"fmt"
	"strings"
)

func InsertQuery(tableName string, columns []string)string{
	colNames := strings.Join(columns,", ")
	placeholders:=strings.Repeat("?, ",len(columns))

	placeholders =strings.TrimSuffix(placeholders,", ")

	query :=fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",tableName,colNames,placeholders)
	return query
}

func SelectQuery(tableName string, cols []string, conditions ...string) string {
	columns := strings.Join(cols, ", ")
	query := fmt.Sprintf("SELECT %s FROM %s", columns, tableName)

	if len(conditions) > 0 {
		var conds []string
		for _, c := range conditions {
			conds = append(conds, fmt.Sprintf("%s = ?", c))
		}
		query += " WHERE " + strings.Join(conds, " AND ")
	}

	return query
}

func DeleteQuery(tableName string, conditions []string) string {
	query := fmt.Sprintf("DELETE FROM %s", tableName)

	if len(conditions) > 0 {
		var condParts []string
		for _, c := range conditions {
			condParts = append(condParts, fmt.Sprintf("%s = ?", c))
		}
		query += " WHERE " + strings.Join(condParts, " AND ")
	}

	return query
}

func UpdateQuery(tableName, condition1, condition2 string, columns []string) string {
	setClause := make([]string, len(columns))
	for i, col := range columns {
		setClause[i] = fmt.Sprintf("%s = ?", col)
	}
	setClauseStr := strings.Join(setClause, ", ")

	if condition2 == "" {
		return fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?", tableName, setClauseStr, condition1)
	}

	return fmt.Sprintf("UPDATE %s SET %s WHERE %s = ? AND %s = ?", tableName, setClauseStr, condition1, condition2)
}

