package internal

import (
	"cloud/lib/logger"
	"fmt"
	model "test/model/xorm/portal"

	"github.com/go-xorm/xorm"
	"github.com/google/go-cmp/cmp"
)

func Portal(db *xorm.Engine) {
	tx := db.NewSession()

	tx.Begin()
	defer tx.Rollback()

	var accs model.Accounts

	if err := tx.Find(&accs); err != nil {
		logger.Error(err)
	}
	logger.Debug(accs)

	// insert
	bio1 := model.Bio{Email: "test1", Name: "name1", Company: 1}
	acc1 := model.Account{Bio: bio1}
	bio2 := model.Bio{Email: "test2", Name: "name2", Company: 2}
	acc2 := model.Account{Bio: bio2}
	if _, err := tx.Insert(&acc1, &acc2); err != nil {
		logger.Error(err)
	}

	if err := AccountsEqual(tx, append(accs, acc1, acc2)); err != nil {
		logger.Error(err)
	}

	// select by id
	acc11 := model.Account{}
	if _, err := tx.Where("id = ?", acc1.ID).Get(&acc11); err != nil {
		logger.Error(err)
	}
	if err := AccountEqual(acc1, acc11); err != nil {
		logger.Error(err)
	}

	// select by name
	acc11 = model.Account{}
	if _, err := tx.Where("name = ?", acc1.Name).Get(&acc11); err != nil {
		logger.Error(err)
	}
	if err := AccountEqual(acc1, acc11); err != nil {
		logger.Error(err)
	}

	// update by id
	acc2.Email = "test22"
	acc2.Company = 3
	if _, err := tx.Where("id = ?", acc2.ID).Update(acc2); err != nil {
		logger.Error(err)
	}

	// update by name
	acc2.Email = "test2"
	acc2.Company = 2
	if _, err := tx.Where("name = ?", acc2.Name).Update(acc2); err != nil {
		logger.Error(err)
	}

	if err := AccountsEqual(tx, append(accs, acc1, acc2)); err != nil {
		logger.Error(err)
	}

	// delete
	if _, err := tx.Delete(acc1); err != nil {
		logger.Error(err)
	}

	if err := AccountsEqual(tx, append(accs, acc2)); err != nil {
		logger.Error(err)
	}

}

func AccountsEqual(tx *xorm.Session, want model.Accounts) error {
	var got model.Accounts
	if err := tx.Find(&got); err != nil {
		logger.Error(err)
		return err
	}
	if !cmp.Equal(want, got) {
		return fmt.Errorf("want:%+v,\n got: %+v", want, got)
	}
	return nil

}

func AccountEqual(want, got model.Account) error {

	if !cmp.Equal(want, got) {
		return fmt.Errorf("want:%+v,\n got: %+v", want, got)
	}
	return nil

}
