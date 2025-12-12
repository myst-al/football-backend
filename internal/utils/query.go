package utils

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FilterOperator string

const (
	OpEq   FilterOperator = "eq"
	OpNe   FilterOperator = "ne"
	OpGt   FilterOperator = "gt"
	OpLt   FilterOperator = "lt"
	OpGte  FilterOperator = "gte"
	OpLte  FilterOperator = "lte"
	OpLike FilterOperator = "like"
	OpIn   FilterOperator = "in"
)

type QueryParams struct {
	Page    int
	Limit   int
	Sort    string
	Order   string
	Filters map[string]map[FilterOperator]string
}

func NewQueryParams() QueryParams {
	return QueryParams{
		Page:    1,
		Limit:   20,
		Sort:    "id",
		Order:   "ASC",
		Filters: make(map[string]map[FilterOperator]string),
	}
}

func ParseQuery(c *gin.Context) QueryParams {
	return ParseFromValues(c.Request.URL.Query())
}

func ParseFromValues(values url.Values) QueryParams {
	q := NewQueryParams()

	if v := values.Get("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			q.Page = p
		}
	}

	if v := values.Get("limit"); v != "" {
		if l, err := strconv.Atoi(v); err == nil && l > 0 {
			q.Limit = l
		}
	}

	if v := values.Get("sort"); v != "" {
		q.Sort = v
	}

	if v := values.Get("order"); v != "" {
		up := strings.ToUpper(v)
		if up == "ASC" || up == "DESC" {
			q.Order = up
		}
	}

	for key, val := range values {
		if !strings.HasPrefix(key, "filter[") {
			continue
		}
		trim := strings.TrimPrefix(key, "filter[")
		trim = strings.TrimSuffix(trim, "]")
		parts := strings.Split(trim, "][")
		if len(parts) != 2 {
			continue
		}
		field := parts[0]
		op := FilterOperator(parts[1])
		if _, ok := q.Filters[field]; !ok {
			q.Filters[field] = make(map[FilterOperator]string)
		}
		q.Filters[field][op] = val[0]
	}

	return q
}

func ApplyFilters(db *gorm.DB, q QueryParams) *gorm.DB {
	for field, ops := range q.Filters {
		for op, val := range ops {
			switch op {
			case OpEq:
				db = db.Where(fmt.Sprintf("%s = ?", field), val)
			case OpNe:
				db = db.Where(fmt.Sprintf("%s <> ?", field), val)
			case OpGt:
				db = db.Where(fmt.Sprintf("%s > ?", field), val)
			case OpLt:
				db = db.Where(fmt.Sprintf("%s < ?", field), val)
			case OpGte:
				db = db.Where(fmt.Sprintf("%s >= ?", field), val)
			case OpLte:
				db = db.Where(fmt.Sprintf("%s <= ?", field), val)
			case OpLike:
				db = db.Where(fmt.Sprintf("%s LIKE ?", field), "%"+val+"%")
			case OpIn:
				parts := strings.Split(val, ",")
				db = db.Where(fmt.Sprintf("%s IN ?", field), parts)
			default:
			}
		}
	}
	return db
}
