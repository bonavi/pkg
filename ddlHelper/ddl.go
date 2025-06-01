package ddlHelper

import (
	"fmt"
	"strconv"
)

func BuildJoin(joinTableWithAlias, joinColumnWithPrefix, columnWithPrefix string) string {
	return joinTableWithAlias + " ON " + joinColumnWithPrefix + " = " + columnWithPrefix
}

func WithCustomAlias(table, alias string) string {
	return table + " " + alias
}

func WithCustomPrefix(column, prefix string) string {
	return prefix + "." + column

}

func Distinct(column string) string {
	return "DISTINCT " + column
}

func As(column, newName string) string {
	return column + " AS " + newName
}

func Cast(column, castType string) string {
	return "CAST(" + column + " AS " + castType + ")"
}

func Divide(column1, column2 string) string {
	return column1 + " / " + column2
}

func Multiply(column string, value string) string {
	return column + " * " + value
}

func Coalesce(column, defaultValue string) string {
	return "COALESCE(" + column + ", " + defaultValue + ")"
}

func Max(column string) string {
	return "MAX(" + column + ")"
}

func Min(column string) string {
	return "MIN(" + column + ")"
}

func Sum(column string) string {
	return "SUM(" + column + ")"
}

func Avg(column string) string {
	return "AVG(" + column + ")"
}

func Plus(column string, value int) string {
	return column + " + " + strconv.Itoa(value)
}

func Minus(column string, value int) string {
	return column + " - " + strconv.Itoa(value)
}

func Count(column string) string {
	return "COUNT(" + column + ")"
}

func Lower(column string) string {
	return "LOWER(" + column + ")"
}

func Desc(column string) string {
	return column + " DESC"
}

func Asc(column string) string {
	return column + " ASC"
}

func PartContains(column string, value any) (string, any) {
	return fmt.Sprintf("%s @> ?", column), value
}

func AllSubstring(value string) string {
	return fmt.Sprintf("%%%s%%", value)
}

const SelectAll = "*"
