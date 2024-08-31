package main

import (
	"log"

	"github.com/setiadipm/queb/queb"
)

type FilterProductDto struct {
	Prod string `json:"prod" db:"prod"`
	Name string `json:"name" db:"name"`
	Merk string `json:"merk" db:"merk"`
}

func main() {
	dto := FilterProductDto{
		Prod: "",
		Name: "Metro White",
		Merk: "Valian",
	}

	rawSql := queb.Build(
		queb.Raw("SELECT * FROM product"),
		queb.Where("prod = :prod", dto.Prod),
		queb.AndWhere("name = :name", dto.Name),
		queb.OrWhere("merk = :merk", dto.Merk),
		queb.AndBracket(
			queb.Where("name = :name", dto.Name),
			queb.AndWhere("name = :name", dto.Name),
			queb.OrWhere("merk = :merk", dto.Merk),
		),
		queb.AndBracket(
			queb.AndWhere("name = :name", dto.Name),
			queb.OrWhere("merk = :merk", dto.Merk),
		),
		queb.OrBracket(
			queb.AndWhere("name = :name", dto.Name),
			queb.OrWhere("merk = :merk", dto.Merk),
		),
		queb.Raw("ORDER BY prod DESC"),
	)

	log.Println(rawSql)
}
