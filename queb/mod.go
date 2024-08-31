package queb

import (
	"fmt"
	"reflect"
	"strings"
)

type queryMod interface {
	toRawSql() string
}

type rawMod struct {
	sql string
}

func (m rawMod) toRawSql() string {
	return strings.TrimSpace(m.sql)
}

func Raw(sql string) rawMod {
	return rawMod{
		sql: sql,
	}
}

type whereMod struct {
	logicalOp  string
	clause     string
	field      interface{}
	isRequired bool
}

func (m whereMod) toRawSql() (sql string) {
	hasValue := hasValue(m.field)
	if m.isRequired || hasValue {
		sql = strings.TrimSpace(m.clause)
	}
	return
}

func Where(clause string, field interface{}, isReqs ...bool) whereMod {
	isRequired := isRequired(isReqs)
	return whereMod{
		logicalOp:  "WHERE",
		clause:     clause,
		field:      field,
		isRequired: isRequired,
	}
}

func AndWhere(clause string, field interface{}, isReqs ...bool) whereMod {
	isRequired := isRequired(isReqs)
	return whereMod{
		logicalOp:  "AND",
		clause:     clause,
		field:      field,
		isRequired: isRequired,
	}
}

func OrWhere(clause string, field interface{}, isReqs ...bool) whereMod {
	isRequired := isRequired(isReqs)
	return whereMod{
		logicalOp:  "OR",
		clause:     clause,
		field:      field,
		isRequired: isRequired,
	}
}

type bracketMod struct {
	logicalOp string
	mods      []queryMod
}

func (m bracketMod) toRawSql() (sql string) {
	for _, mod := range m.mods {
		if len(mod.toRawSql()) <= 0 {
			continue
		}

		switch mod := mod.(type) {
		case rawMod:
			sql = fmt.Sprintf("%s %s", sql, mod.toRawSql())
		case whereMod:
			if len(sql) <= 0 {
				sql = fmt.Sprintf("%s %s", sql, mod.toRawSql())
			} else {
				sql = fmt.Sprintf("%s %s %s", sql, mod.logicalOp, mod.toRawSql())
			}
		}
		sql = strings.TrimSpace(sql)
	}

	if len(sql) > 0 {
		sql = fmt.Sprintf("%s (%s)", m.logicalOp, strings.TrimSpace(sql))
	}

	return
}

func WhereBracket(mods ...queryMod) bracketMod {
	return bracketMod{
		logicalOp: "WHERE",
		mods:      mods,
	}
}

func AndBracket(mods ...queryMod) bracketMod {
	return bracketMod{
		logicalOp: "AND",
		mods:      mods,
	}
}

func OrBracket(mods ...queryMod) bracketMod {
	return bracketMod{
		logicalOp: "OR",
		mods:      mods,
	}
}

func isRequired(isReqs []bool) (isRequired bool) {
	isRequired = false
	if len(isReqs) > 0 {
		isRequired = isReqs[0]
	}
	return
}

func hasValue(value interface{}) (hasValue bool) {
	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Interface, reflect.Chan:
		hasValue = !v.IsNil()
	case reflect.String:
		hasValue = v.Len() != 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		hasValue = v.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		hasValue = v.Uint() != 0
	case reflect.Float32, reflect.Float64:
		hasValue = v.Float() != 0.0
	case reflect.Bool:
		hasValue = true
	default:
		hasValue = false
	}
	return
}
