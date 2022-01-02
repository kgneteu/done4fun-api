package models

import (
	"strconv"
	"strings"
)

func (model *DBModel) GetAvailablePrizes(userId uint) (prizes *[]Prize, err error) {
	sqlStatement := `SELECT * FROM prizes WHERE kid_id=$1;`
	prizes = &[]Prize{}
	err = model.Db.Select(prizes, sqlStatement, userId)
	return
}

func (model *DBModel) GetPrize(prizeId uint) (prize *Prize, err error) {
	sqlStatement := `SELECT * FROM prizes WHERE id=$1;`
	prize = &Prize{}
	err = model.Db.Get(prize, sqlStatement, prizeId)
	return prize, err
}

func (model *DBModel) DeletePrize(prizeId uint) (err error) {
	sqlStatement := `DELETE FROM prizes WHERE id=$1;`
	_, err = model.Db.Exec(sqlStatement, prizeId)
	return err
}

func (model *DBModel) CreatePrize(prize map[string]string) (id uint, err error) {
	var fields []string
	var values []interface{}
	var placeholders []string
	i := 1
	for k, v := range prize {
		fields = append(fields, k)
		values = append(values, v)
		placeholders = append(placeholders, "$"+strconv.Itoa(i))
		i++
	}
	fString := strings.Join(fields, ", ")
	pString := strings.Join(placeholders, ", ")
	sqlStatement := `
		INSERT INTO prizes (` + fString + `)
		VALUES (` + pString + `)
		RETURNING id`

	err = model.Db.QueryRow(sqlStatement, values...).Scan(&id)
	return
}

func (model *DBModel) UpdatePrize(prize map[string]string, prizeId uint) (err error) {
	var fields []string
	var values []interface{}
	values = append(values, prizeId)
	i := 2
	for k, v := range prize {
		fields = append(fields, k+"=$"+strconv.Itoa(i))
		values = append(values, v)
		i++
	}

	fString := strings.Join(fields, ", ")
	sqlStatement := `UPDATE prizes SET ` + fString + ", updated_at=NOW() WHERE id=$1"
	_, err = model.Db.Exec(sqlStatement, values...)
	return
}
