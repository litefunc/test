package internal

import (
	"cloud/lib/logger"
	model "test/model/gorm/portal"

	"github.com/jinzhu/gorm"
)

func Portal(db *gorm.DB) {
	tx := db.Begin()
	defer tx.Rollback()

	var accs model.Accounts
	logger.Error(tx.Find(&accs).Error)
	logger.Debug(accs)
	logger.Debug(len(accs))

	var acc, acc1 model.Account

	bio := model.Bio{Email: "test1", Company: 1}
	acc = model.Account{Bio: bio}

	logger.Error(tx.Create(&acc).Error)
	logger.Error(tx.Find(&accs).Error)
	logger.Debug(len(accs))

	acc.Name = "n1"
	logger.Error(tx.Save(acc).Error)
	logger.Error(tx.Where(&acc).Find(&acc1).Error)
	logger.Debug(acc1)

	logger.Error(tx.Delete(acc).Error)
	logger.Error(tx.Where(&acc).Find(&acc1).Error)
	logger.Debug(acc1)

	logger.Error(tx.Find(&accs).Error)
	logger.Debug(len(accs))

}
