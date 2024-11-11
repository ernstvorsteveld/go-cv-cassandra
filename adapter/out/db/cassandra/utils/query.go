package cassandra_util

import (
	"fmt"

	"github.com/ernstvorsteveld/go-cv-cassandra/domain/port/out"
	"github.com/pkg/errors"
)

const (
	EXP_TABLE_NAME = "cv_experiences"
)

var (
	Stmt_insert       string = fmt.Sprintf("INSERT INTO %s (id, name, tags) VALUES (?, ?, ?)", EXP_TABLE_NAME)
	Stmt_select_by_id string = fmt.Sprintf("SELECT id, name, tags FROM %s WHERE id = ?", EXP_TABLE_NAME)
	Stmt_update       string = fmt.Sprintf("UPDATE %s SET name = ?, tags = ? WHERE id = ?", EXP_TABLE_NAME)
	QryErrorNotFound         = errors.Errorf("Not Found")
	Stmt_select_page  string = fmt.Sprintf("SELECT id, name, tags FROM %s WHERE (1==1)", EXP_TABLE_NAME)
)

func GetStatement(params *out.GetParams) string {
	limit := getLimit(params.Limit)
	page := getPage(params.Page)
	tag := getTag(params.Tag)
	name := getName(params.Name)

	var stmt string = Stmt_select_page

	if page != "" {
		stmt = fmt.Sprintf("%s AND page > '%s' ", stmt, page)
	}
	if tag != "" {
		return fmt.Sprintf("%s AND tags CONTAINS '%s' LIMIT %d", stmt, tag, limit)
	}
	if name != "" {
		return fmt.Sprintf("%s AND name CONTAINS '%s' LIMIT %d", stmt, name, limit)
	}

	return fmt.Sprintf("%s LIMIT %d", Stmt_select_page, limit)
}

func getName(name *string) string {
	if name == nil {
		return ""
	}
	if *name == "" {
		return ""
	}
	return *name
}

func getPage(page *string) string {
	if page == nil {
		return ""
	}
	if *page == "" {
		return ""
	}
	return *page
}

func getTag(tag *string) string {
	if tag == nil {
		return ""
	}
	if *tag == "" {
		return ""
	}
	return *tag
}

func getLimit(limit *int32) int32 {
	if limit == nil {
		return 100
	}
	if *limit < 1 {
		return 100
	}
	return *limit
}
