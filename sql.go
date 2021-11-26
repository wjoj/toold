package toold

import (
	"fmt"
	"reflect"
)

//WhereIfsType WhereIfsType
type WhereIfsType string

//
const (
	WhereIfsTypeLike         WhereIfsType = " like "
	WhereIfsTypeEqual                     = "="
	WhereIfsTypeLessEqual                 = "<="
	WhereIfsTypeGreaterEqual              = ">="
	WhereIfsTypeGreater                   = ">"
	WhereIfsTypeLess                      = "<"
)

//WhereSource WhereSource
type WhereSource struct {
	s        interface{}
	whs      map[string]*whIfs
	sqlWhere *Where
}

type whIfs struct {
	ifs    WhereIfsType
	field  []string
	sign   string
	notVal string
}

//And and
func (w *WhereSource) And(field []string, ifs WhereIfsType, key string, notval string) *WhereSource {
	w.whs[key] = &whIfs{
		ifs:    ifs,
		field:  field,
		sign:   "and",
		notVal: notval,
	}
	return w
}

//Or or
func (w *WhereSource) Or(field []string, ifs WhereIfsType, key string, notval string) *WhereSource {
	w.whs[key] = &whIfs{
		ifs:    ifs,
		field:  field,
		sign:   "or",
		notVal: notval,
	}
	return w
}

//Done done
func (w *WhereSource) Done() *Where {
	v := w.s
	immutable := reflect.ValueOf(v)
	if immutable.Kind() == reflect.Ptr {
		immutable = immutable.Elem()
	}
	typeIm := reflect.TypeOf(v)
	if typeIm.Kind() == reflect.Ptr {
		typeIm = typeIm.Elem()
	}
	sqlm := ""
	sqlifs := ""
	for i := 0; i < typeIm.NumField(); i++ {
		field := typeIm.Field(i)
		key := field.Tag.Get("json")
		if len(key) == 0 || key == "-" {
			continue
		}
		name := field.Name
		val := immutable.FieldByName(name)
		obj := w.whs[key]
		if obj != nil {
			if len(obj.field) == 0 || fmt.Sprintf("%v", val) == obj.notVal {
				continue
			}
			strval := ""
			if val.Kind() == reflect.String {
				strval = fmt.Sprintf("'%v'", val)
			} else {
				strval = fmt.Sprintf("%v", val)
			}

			sql := ""
			for _, f := range obj.field {
				if len(sql) == 0 {
					sql = fmt.Sprintf("%v%v%v", f, obj.ifs, strval)
				} else {
					sql += fmt.Sprintf(" or %v%v%v", f, obj.ifs, strval)
				}
			}
			if len(obj.field) != 1 {
				sql = fmt.Sprintf("(%v)", sql)
			}
			if len(sql) == 0 {
				continue
			}
			if len(sqlm) == 0 {
				sqlm = sql
				sqlifs = obj.sign
			} else {
				sqlm += fmt.Sprintf(" %v %v", obj.sign, sql)
			}
		}
	}

	if len(sqlm) == 0 {
		return w.sqlWhere
	}
	if len(*w.sqlWhere) == 0 {
		*w.sqlWhere = Where(sqlm)
	} else {
		*w.sqlWhere += Where(fmt.Sprintf(" %v %v", sqlifs, sqlm))
	}
	return w.sqlWhere
}

//Where Where
type Where string

//AndSQL AndSQL
func (w *Where) AndSQL(sql string) {
	if len(*w) == 0 {
		*w = Where(sql)
	} else {
		*w += Where(fmt.Sprintf(" and %v", sql))
	}
}

//ORSQL ORSQL
func (w *Where) ORSQL(sql string) {
	if len(*w) == 0 {
		*w = Where(sql)
	} else {
		*w += Where(fmt.Sprintf(" or %v", sql))
	}
}

//And and
func (w *Where) And(key, ifs, val interface{}, not interface{}) {
	if val == not {
		return
	}

	sql := fmt.Sprintf("%v%v%v", key, ifs, val)
	if len(*w) == 0 {
		*w = Where(sql)
	} else {
		*w += Where(fmt.Sprintf(" and %v", sql))
	}
}

//Or or
func (w *Where) Or(key, ifs, val interface{}, not interface{}) {
	if val == not {
		return
	}
	sql := fmt.Sprintf("%v%v%v", key, ifs, val)
	if len(*w) == 0 {
		*w = Where(sql)
	} else {
		*w += Where(fmt.Sprintf(" or %v", sql))
	}
}

//AndWhere AndWhere
func (w *Where) AndWhere(where *Where) {
	if len(*where) != 0 {
		if len(*w) != 0 {
			*w += Where(fmt.Sprintf(" and %v", *where))
		} else {
			*w = *where
		}
	}
}

//OrWhere OrWhere
func (w *Where) OrWhere(where *Where) {
	if len(*where) != 0 {
		if len(*w) != 0 {
			*w += Where(fmt.Sprintf(" or %v", *where))
		} else {
			*w = *where
		}
	}
}

//SearchKeysSQLDoneOr SearchKeysSQLDoneOr
func (w *Where) SearchKeysSQLDoneOr(keys []string, ifs string, val interface{}, not interface{}) {
	if val == not {
		return
	}
	wh := ""
	for _, key := range keys {
		if len(key) == 0 {
			continue
		}
		sql := fmt.Sprintf("%v%v%v", key, ifs, val)
		if len(wh) == 0 {
			wh = sql
		} else {
			wh += fmt.Sprintf(" or %v", sql)
		}
	}
	if len(wh) == 0 {
		return
	}
	if len(*w) == 0 {
		*w = Where(wh)
	} else {
		*w += Where(fmt.Sprintf(" and %v", wh))
	}
}

//SearchKeysSQLDoneOr SearchKeysSQLDoneOr
func (w *Where) SearchFirstKeysSQLDoneOr(keys []string, ifs string, val interface{}, not interface{}) {
	if val == not {
		return
	}
	wh := ""
	for _, key := range keys {
		if len(key) == 0 {
			continue
		}
		sql := fmt.Sprintf("%v%v%v", key, ifs, val)
		if len(wh) == 0 {
			wh = sql
		} else {
			wh += fmt.Sprintf(" or %v", sql)
		}
	}
	if len(wh) == 0 {
		return
	}
	if len(*w) == 0 {
		*w = Where(fmt.Sprintf("(%v)", wh))
	} else {
		*w += Where(fmt.Sprintf(" and (%v)", wh))
	}
}

//Source Source
func (w *Where) Source(v interface{}) *WhereSource {
	return &WhereSource{
		s:        v,
		whs:      make(map[string]*whIfs),
		sqlWhere: w,
	}
}

//GetWhere GetWhere
func (w *Where) GetWhere() string {
	if len(*w) == 0 {
		return ""
	}
	return fmt.Sprintf("%s", *w)
}
