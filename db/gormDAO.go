package db

import (
	"cbridgewrapper/config"
	"cbridgewrapper/entity"
	"cbridgewrapper/logger"
	"database/sql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var gormDB *gorm.DB
var MyDAO *GormDAO

type GormDAO struct {
	TransactionDAO TransactionDAO
}

func init() {
	dsn := config.GetConfig("psql.connectstring")
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Logger.Fatalf("Can't connect to db")
	}
	gormDB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		logger.Logger.Fatalf("Can't wrap gorm %v", err)
	}

	gormDB.AutoMigrate(&entity.RelayTransaction{})
	MyDAO = &GormDAO{
		TransactionDAO: TransactionDAO{},
	}
}
