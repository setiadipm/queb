package queb

import (
	"fmt"
	"strings"
)

func Build(mods ...queryMod) (sql string) {
	err := validateMods(mods)
	if err != nil {
		panic(err)
	}

	sql = generateSql(mods)

	return
}

func validateMods(mods []queryMod) (err error) {
	whereCount := 0
	hasNonWhereCount := false

	for _, mod := range mods {
		switch mod := mod.(type) {
		case whereMod:
			if mod.logicalOp == "WHERE" {
				whereCount++
			}
			if mod.logicalOp != "WHERE" {
				hasNonWhereCount = true
			}

			if mod.logicalOp == "WHERE" {
				err = validateWhereMod(whereCount, hasNonWhereCount, mod)
				if err != nil {
					return
				}
			}
		case bracketMod:
			if mod.logicalOp == "WHERE" {
				whereCount++
			}

			if mod.logicalOp == "WHERE" {
				err = validateWhereBracketMod(whereCount, hasNonWhereCount, mod)
				if err != nil {
					return
				}
			}
			err = validateMods(mod.mods)
			if err != nil {
				return
			}
		}
	}
	return
}

func validateWhereMod(whereCount int, hasNonWhereCount bool, mod whereMod) (err error) {
	if whereCount > 1 {
		err = fmt.Errorf("second Where mod is found on clause \"%s\"", mod.clause)
	}
	if hasNonWhereCount {
		err = fmt.Errorf("invalid Where mod position found on clause \"%s\"", mod.clause)
	}
	return
}

func validateWhereBracketMod(whereCount int, hasNonWhereCount bool, mod bracketMod) (err error) {
	if whereCount > 1 {
		err = fmt.Errorf("second Where bracket mod is found for mod \"%s\"", mod)
	}
	if hasNonWhereCount {
		err = fmt.Errorf("invalid Where bracket mod position found for mod \"%s\"", mod)
	}
	return
}

func generateSql(mods []queryMod) (sql string) {
	for _, mod := range mods {
		if len(mod.toRawSql()) <= 0 {
			continue
		}

		switch mod := mod.(type) {
		case rawMod:
			sql = fmt.Sprintf("%s %s", sql, mod.toRawSql())
		case whereMod:
			switch mod.logicalOp {
			case "WHERE":
				if !strings.Contains(sql, "WHERE") {
					sql = fmt.Sprintf("%s WHERE %s", sql, mod.toRawSql())
				}
			case "AND":
				if !strings.Contains(sql, "WHERE") {
					sql = fmt.Sprintf("%s WHERE %s", sql, mod.toRawSql())
				} else {
					sql = fmt.Sprintf("%s AND %s", sql, mod.toRawSql())
				}
			case "OR":
				if !strings.Contains(sql, "WHERE") {
					sql = fmt.Sprintf("%s WHERE %s", sql, mod.toRawSql())
				} else {
					sql = fmt.Sprintf("%s OR %s", sql, mod.toRawSql())
				}
			}
		case bracketMod:
			sql = fmt.Sprintf("%s %s", sql, mod.toRawSql())
		}
		sql = strings.TrimSpace(sql)
	}

	return
}
