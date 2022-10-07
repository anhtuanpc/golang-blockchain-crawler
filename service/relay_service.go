package service

import (
	"cbridgewrapper/client"
	"cbridgewrapper/constant"
	"cbridgewrapper/db"
	"cbridgewrapper/entity"
	"cbridgewrapper/logger"
	"cbridgewrapper/utils"
	"go.mongodb.org/mongo-driver/bson"
)

var queue chan entity.RelayTransaction

func init() {
	queue = make(chan entity.RelayTransaction, 10000)
	txs, err := db.MyDAO.TransactionDAO.FindPendingTxs()
	if err == nil {
		for _, tx := range txs {
			queue <- tx
		}
	}
}

func AddTransactionToQueue(trans entity.RelayTransaction) {
	queue <- trans
}

func ProcessTransaction() {
	defer ProcessTransaction()
	for {
		logger.Logger.Infof("Waiting for pending request")
		tx := <-queue
		raw, _ := utils.Marshal(tx)
		logger.Logger.Infof("processing request", string(raw))
		tx.Status = constant.PROCESSING
		err := db.MyDAO.TransactionDAO.Update(tx)
		if err != nil {
			logger.Logger.WithField("transferId", tx.TransferId).Errorf("Update transaction status to PROCESSING fail")
			continue
		}

		err = db.UpdateRelayStatus(bson.M{"transferId": tx.TransferId, "status": constant.PENDING},
			bson.D{
				{"$set", bson.D{{"status", constant.PROCESSING}}},
			})
		if err != nil {
			logger.Logger.WithField("transferId", tx.TransferId).WithField("status", constant.FAIL).Errorf("Sync status to moralis fail")
			continue
		}

		logger.Logger.WithField("transferId", tx.TransferId).Infof("processing swap")
		transaction, err := client.GetEthClient(tx.DestinationChainId).Swap(tx)
		if err != nil {
			logger.Logger.WithField("transferId", tx.TransferId).Errorf("swap fail with err %s", err.Error())
			tx.Status = constant.FAIL
			tx.Result = err.Error()
			err := db.MyDAO.TransactionDAO.Update(tx)
			if err != nil {
				logger.Logger.WithField("transferId", tx.TransferId).Errorf("Update transaction status to FAIL fail")
			} else {
				err = db.UpdateRelayStatus(bson.M{"transferId": tx.TransferId, "status": constant.PENDING},
					bson.D{
						{"$set", bson.D{{"status", constant.FAIL}}},
					})
				if err != nil {
					logger.Logger.WithField("transferId", tx.TransferId).WithField("status", constant.FAIL).Errorf("Sync status to moralis fail")
				}
			}
		} else {
			logger.Logger.WithField("transferId", tx.TransferId).WithField("transactionId", transaction.Hash().Hex()).Infof("Relay success")
			tx.Status = constant.RELAYED
			err := db.MyDAO.TransactionDAO.Update(tx)
			if err != nil {
				logger.Logger.WithField("transferId", tx.TransferId).WithField("transactionId", transaction.Hash().Hex()).Errorf("Update transaction status to RELAYED fail")
			} else {
				err = db.UpdateRelayStatus(bson.M{"transferId": tx.TransferId, "status": constant.PENDING},
					bson.D{
						{"$set", bson.D{{"status", constant.RELAYED}}},
					})
				if err != nil {
					logger.Logger.WithField("transferId", tx.TransferId).WithField("status", constant.RELAYED).Errorf("Sync status to moralis fail")
				}
			}
		}
	}
}
