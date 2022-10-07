package service

import (
	"cbridgewrapper/db"
	"cbridgewrapper/logger"
	"github.com/go-co-op/gocron"
	"time"
)

var Scheduler *gocron.Scheduler

func init() {

	Scheduler = gocron.NewScheduler(time.UTC)

	Scheduler.Every(1).Minute().Do(func() { fetchService() })

	Scheduler.StartAsync()
}

func fetchService() {
	logger.Logger.Infof("Fetching pending request from moralis")
	trans, err := db.FindAllPendingRelayTransaction()
	if err == nil {
		for _, trans := range trans {
			err := db.MyDAO.TransactionDAO.Create(trans)
			if err == nil {
				AddTransactionToQueue(trans)
			}
		}
	}
}
