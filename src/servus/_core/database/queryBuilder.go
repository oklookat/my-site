package database

import (
	"fmt"
	"strings"
)

type queryBuilded struct {
	fields string
	dollars string
	values []interface{}
}

func queryBuilder(userData map[string]string, columns []string) queryBuilded{
	var fields []string
	var dollars []string
	var values []interface{}
	i := 1
	for _, k := range columns {
		if _, ok := userData[k]; ok {
			if len(userData[k]) < 1{
				continue
			}
			values = append(values, userData[k])
			dollars = append(dollars, fmt.Sprintf("$%v", i))
			fields = append(fields, fmt.Sprintf("%s", k))
			i++
		}
	}
	return queryBuilded{fields: strings.Join(fields,", "), dollars: strings.Join(dollars,", "), values: values}
}
