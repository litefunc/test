package pgsql

import (
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

func TestTxBasicCRUD(t *testing.T) {

	dbConfig := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5431, "test", "abcd", "test")
	db, _ := Connect(dbConfig)
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer tx.Rollback()

	var mdas TbAas
	tx.Select(&mdas).SQL()
	if err := tx.Select(&mdas).Run(); err != nil {
		t.Error(err)
		return
	}
	// len(mdas) should be 0
	if err := ModelsEqual(tx, mdas); err != nil {
		t.Error(err)
		return
	}

	if err := tx.Insert(mda1).Returning("id").Run().Scan(&mda1.ID); err != nil {
		t.Error(err)
		return
	}
	if err := tx.Insert(mda2).Returning("id").Run().Scan(&mda2.ID); err != nil {
		t.Error(err)
		return
	}

	md := TbAa{ID: mda1.ID}
	if err := tx.SelectByPk(&md).Run(); err != nil {
		t.Error(err)
		return
	}
	if err := ModelEqual(mda1, md); err != nil {
		t.Error(err)
		return
	}

	if err := ModelsEqual(tx, append(mdas, mda1, mda2)); err != nil {
		t.Error(err)
		return
	}

	mda2.EmbedAa = mda2.EmbedAa * 10
	mda2.EmbedAb = mda2.EmbedAb + mda2.EmbedAb
	mda2.Note = mda2.Note + mda2.Note
	if _, err := tx.UpdateByPk(mda2).Run(); err != nil {
		t.Error(err)
		return
	}

	if err := ModelsEqual(tx, append(mdas, mda1, mda2)); err != nil {
		t.Error(err)
		return
	}

	var n int
	if err := tx.Count(mdas).Run().Scan(&n); err != nil {
		t.Error(err)
		return
	}
	if n != 2 {
		t.Error(n)
	}

	n = 0
	if err := tx.CountDistinct(mdas, "note").Run().Scan(&n); err != nil {
		t.Error(err)
		return
	}
	if n != 2 {
		t.Error(n)
	}

	if _, err := tx.Delete(mda2).Run(); err != ErrDeleteWithoutCondition {
		t.Error(err)
		return
	}

	if _, err := tx.DeleteByPk(mda2).Run(); err != nil {
		t.Error(err)
		return
	}
	if err := ModelsEqual(tx, append(mdas, mda1)); err != nil {
		t.Error(err)
		return
	}

	if _, err := tx.Truncate(mdas).Run(); err != nil {
		t.Error(err)
		return
	}

	tx.Select(&mdas).SQL()
	if err := tx.Select(&mdas).Run(); err != nil {
		t.Error(err)
		return
	}
	// len(mdas) should be 0
	if err := ModelsEqual(tx, mdas); err != nil {
		t.Error(err)
		return
	}

}

func TestTxWhere(t *testing.T) {

	dbConfig := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5431, "test", "abcd", "test")
	db, _ := Connect(dbConfig)
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		t.Error(err)
		return
	}

	var mdas TbAas
	tx.Select(&mdas).SQL()
	if err := tx.Select(&mdas).Run(); err != nil {
		t.Error(err)
		return
	}
	// len(mdas) should be 0
	if err := ModelsEqual(tx, mdas); err != nil {
		t.Error(err)
		return
	}

	defer tx.Truncate(mdas).Run()

	if err := tx.Insert(mda1).Returning("id").Run().Scan(&mda1.ID); err != nil {
		t.Error(err)
		return
	}
	if err := tx.Insert(mda2).Returning("id").Run().Scan(&mda2.ID); err != nil {
		t.Error(err)
		return
	}

	if err := ModelsEqual(tx, append(mdas, mda1, mda2)); err != nil {
		t.Error(err)
		return
	}

	md := TbAa{}
	if err := tx.Select(&md).Where("note=?", mda1.Note).Run(); err != nil {
		t.Error(err)
		return
	}
	if err := ModelEqual(mda1, md); err != nil {
		t.Error(err)
		return
	}

	mda2.EmbedAa = mda2.EmbedAa * 10
	mda2.EmbedAb = mda2.EmbedAb + mda2.EmbedAb
	if _, err := tx.Update(TbAa{}).Set("embed_aa=?, embed_ab=?", mda2.EmbedAa, mda2.EmbedAb).Where("note=?", mda2.Note).Run(); err != nil {
		t.Error(err)
		return
	}

	if err := ModelsEqual(tx, append(mdas, mda1, mda2)); err != nil {
		t.Error(err)
		return
	}

	var n int
	if err := tx.Count(TbAa{}).Where("note=?", mda2.Note).Run().Scan(&n); err != nil {
		t.Error(err)
		return
	}
	if n != 1 {
		t.Error(n)
	}

	n = 0
	if err := tx.CountDistinct(TbAa{}, "note").Where("note=?", mda2.Note).Run().Scan(&n); err != nil {
		t.Error(err)
		return
	}
	if n != 1 {
		t.Error(n)
	}

	if _, err := tx.Delete(TbAa{}).Where("note=?", mda2.Note).Run(); err != nil {
		t.Error(err)
		return
	}

	if err := ModelsEqual(tx, append(mdas, mda1)); err != nil {
		t.Error(err)
		return
	}

}
