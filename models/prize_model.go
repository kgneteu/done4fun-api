package models

import (
	"strconv"
	"strings"
)

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

func (model *DBModel) CreatePrize(prize map[string]interface{}) (id uint, err error) {
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

func (model *DBModel) UpdatePrize(prize map[string]interface{}, prizeId uint) (err error) {
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

type PrizeList struct {
	Total  int
	Prizes *[]Prize
}

func (model *DBModel) GetAvailablePrizes(userId uint, page int, limit int, order string) (prizeList PrizeList, err error) {
	if page < 1 {
		page = 1
	}
	if order == "" {
		order = "'id'"
	}
	if limit < 0 {
		limit = 0
	}
	offset := (page - 1) * limit
	prizeList.Prizes = &[]Prize{}

	if limit > 0 {
		sqlStatement := `SELECT * FROM prizes WHERE kid_id=$1 AND published=true ORDER BY $4 OFFSET $2 LIMIT $3`
		err = model.Db.Select(prizeList.Prizes, sqlStatement, userId, offset, limit, order)
		if err = model.Db.Get(&prizeList.Total, `SELECT COUNT(*) as total FROM prizes WHERE kid_id=$1 and published=true`, userId); err != nil {
			return
		}
	} else {
		sqlStatement := `SELECT * FROM prizes WHERE kid_id=$1 AND published=true ORDER BY $3 OFFSET $2`
		if err = model.Db.Select(prizeList.Prizes, sqlStatement, userId, offset, order); err != nil {
			prizeList.Total = len(*prizeList.Prizes)
		}
	}

	return
}
