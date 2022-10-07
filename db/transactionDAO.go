package db

import (
	"cbridgewrapper/constant"
	"cbridgewrapper/entity"
)

type TransactionDAO struct {
}

func (t *TransactionDAO) Find(txId string) (*entity.RelayTransaction, error) {
	var trans entity.RelayTransaction
	tx := gormDB.First(&trans, entity.RelayTransaction{TransferId: txId})
	return &trans, tx.Error
}
func (t *TransactionDAO) FindPendingTxs() ([]entity.RelayTransaction, error) {
	var data []entity.RelayTransaction
	tx := gormDB.Where("status = ?", constant.PENDING).Find(&data)
	return data, tx.Error
}
func (t *TransactionDAO) Create(transaction entity.RelayTransaction) error {
	tx := gormDB.Create(&transaction)
	return tx.Error
}

func (t *TransactionDAO) Update(transaction entity.RelayTransaction) error {
	tx := gormDB.Model(&entity.RelayTransaction{TransferId: transaction.TransferId}).Updates(transaction)
	return tx.Error
}
