package models

func (model *DBModel) GetAvailablePrizes(userId uint) (prizes *[]Prize, err error) {
	sqlStatement := `SELECT * FROM prizes WHERE kid_id=$1;`
	prizes = &[]Prize{}
	err = model.Db.Select(prizes, sqlStatement, userId)
	return
}
