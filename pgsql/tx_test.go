package pgsql

import (
	"errors"
	"fmt"
	"testing"

	"github.com/lib/pq"
)

func TestTxBasicCRUD(t *testing.T) {

	dbConfig := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5433, "test", "abcd", "test")
	db, err := Connect(dbConfig)
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer tx.Rollback()
	test := NewTester(t)

	var mdas TbAas
	test.Run(tx.Select(&mdas))

	// len(mdas) should be 0
	if err := ModelsEqual(tx, mdas); err != nil {
		t.Error(err)
		return
	}

	mda1 := NewTbAa(1, "b1", "n1")
	mda2 := NewTbAa(2, "b2", "n2")
	test.Scan(tx.Insert(mda1).Returning("id").Run(), &mda1.ID)
	test.Scan(tx.Insert(mda2).Returning("id").Run(), &mda2.ID)

	md := TbAa{ID: mda1.ID}
	test.Run(tx.SelectByPk(&md))
	if err := ModelEqual(mda1, md); err != nil {
		t.Error(err)
		return
	}

	md = TbAa{ID: mda1.ID}
	test.Run(tx.SelectByPk(&md, "embed_aa", "embed_ab"))
	if err := ModelEqual(TbAa{ID: mda1.ID, Embed: Embed{EmbedAa: 1, EmbedAb: "b1"}, Note: ""}, md); err != nil {
		t.Error(err)
		return
	}

	// column order dosen't matter
	md = TbAa{}
	test.Run(tx.Select(&md, "note", "embed_ab", "embed_aa").Where("id=?", mda1.ID))
	if err := ModelEqual(TbAa{Embed: Embed{EmbedAa: 1, EmbedAb: "b1"}, Note: "n1"}, md); err != nil {
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
	test.Run(tx.UpdateByPk(mda2))
	if err := ModelsEqual(tx, append(mdas, mda1, mda2)); err != nil {
		t.Error(err)
		return
	}

	var n int
	test.Scan(tx.Count(mdas).Run(), &n)
	if n != 2 {
		t.Error(n)
		return
	}

	n = 0
	test.Scan(tx.CountDistinct(mdas, "note").Run(), &n)
	if n != 2 {
		t.Error(n)
		return
	}

	if err := tx.Delete(mda2).Run(); err != ErrDeleteWithoutCondition {
		t.Error(err)
		return
	}

	test.Run(tx.DeleteByPk(mda2))
	if err := ModelsEqual(tx, append(mdas, mda1)); err != nil {
		t.Error(err)
		return
	}

	test.Run(tx.Truncate(mdas))
	test.Run(tx.Select(&mdas))

	// len(mdas) should be 0
	if err := ModelsEqual(tx, mdas); err != nil {
		t.Error(err)
		return
	}

}

func TestTxWhere(t *testing.T) {

	dbConfig := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5433, "test", "abcd", "test")
	db, err := Connect(dbConfig)
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer tx.Rollback()
	test := NewTester(t)

	var mdas TbAas
	test.Run(tx.Select(&mdas))

	// len(mdas) should be 0
	if err := ModelsEqual(tx, mdas); err != nil {
		t.Error(err)
		return
	}

	defer tx.Truncate(mdas).Run()
	mda1 := NewTbAa(1, "b1", "n1")
	mda2 := NewTbAa(2, "b2", "n2")
	test.Scan(tx.Insert(mda1).Returning("id").Run(), &mda1.ID)
	test.Scan(tx.Insert(mda2).Returning("id").Run(), &mda2.ID)
	if err := ModelsEqual(tx, append(mdas, mda1, mda2)); err != nil {
		t.Error(err)
		return
	}

	md := TbAa{}
	test.Run(tx.Select(&md).Where("note=?", mda1.Note))
	if err := ModelEqual(mda1, md); err != nil {
		t.Error(err)
		return
	}

	mda2.EmbedAa = mda2.EmbedAa * 10
	mda2.EmbedAb = mda2.EmbedAb + mda2.EmbedAb
	test.Run(tx.Update(TbAa{}).Set("embed_aa=?, embed_ab=?", mda2.EmbedAa, mda2.EmbedAb).Where("note=?", mda2.Note))
	if err := ModelsEqual(tx, append(mdas, mda1, mda2)); err != nil {
		t.Error(err)
		return
	}

	var n int
	test.Scan(tx.Count(TbAa{}).Where("note=?", mda2.Note).Run(), &n)
	if n != 1 {
		t.Error(n)
		return
	}

	n = 0
	test.Scan(tx.CountDistinct(TbAa{}, "note").Where("note=?", mda2.Note).Run(), &n)
	if n != 1 {
		t.Error(n)
		return
	}

	test.Run(tx.Delete(TbAa{}).Where("note=?", mda2.Note))
	if err := ModelsEqual(tx, append(mdas, mda1)); err != nil {
		t.Error(err)
		return
	}

}

func TestTxWhereIn(t *testing.T) {

	dbConfig := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5433, "test", "abcd", "test")
	db, err := Connect(dbConfig)
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer tx.Rollback()
	test := NewTester(t)

	var mdas TbAas
	test.Run(tx.Select(&mdas))

	// len(mdas) should be 0
	if err := ModelsEqual(tx, mdas); err != nil {
		t.Error(err)
		return
	}

	defer tx.Truncate(mdas).Run()
	mda1 := NewTbAa(1, "b1", "n1")
	mda2 := NewTbAa(2, "b2", "n2")
	test.Scan(tx.Insert(mda1).Returning("id").Run(), &mda1.ID)
	test.Scan(tx.Insert(mda2).Returning("id").Run(), &mda2.ID)
	if err := ModelsEqual(tx, append(mdas, mda1, mda2)); err != nil {
		t.Error(err)
		return
	}

	mds := TbAas{}
	test.Run(tx.Select(&mds).Where("note IN (?)", pq.Array([]string{mda1.Note, mda2.Note})))
	if err := ModelEqual(TbAas{mda1, mda2}, mds); err != nil {
		t.Error(err)
		return
	}
	mds = TbAas{}
	test.Run(tx.Select(&mds).Where("note IN (?)", pq.Array([]string{mda2.Note})))
	if err := ModelEqual(TbAas{mda2}, mds); err != nil {
		t.Error(err)
		return
	}

	mds = TbAas{}
	test.Run(tx.Select(&mds).Where("note").In(pq.Array([]string{mda1.Note, mda2.Note})))
	if err := ModelEqual(TbAas{mda1, mda2}, mds); err != nil {
		t.Error(err)
		return
	}
	mds = TbAas{}
	test.Run(tx.Select(&mds).Where("note").In(pq.Array([]string{mda2.Note})))
	if err := ModelEqual(TbAas{mda2}, mds); err != nil {
		t.Error(err)
		return
	}

}

func TestTxJoin(t *testing.T) {

	dbConfig := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5433, "test", "abcd", "test")
	db, err := Connect(dbConfig)
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer tx.Rollback()
	test := NewTester(t)

	var mdas TbAas
	test.Run(tx.Select(&mdas))

	// len(mdas) should be 0
	if err := ModelsEqual(tx, mdas); err != nil {
		t.Error(err)
		return
	}

	defer tx.Truncate(mdas).Run()
	mda1 := NewTbAa(1, "b1", "n1")
	mda2 := NewTbAa(2, "b2", "n2")
	test.Scan(tx.Insert(mda1).Returning("id").Run(), &mda1.ID)
	test.Scan(tx.Insert(mda2).Returning("id").Run(), &mda2.ID)

	q := `SELECT a.note, a.embed_ab, a.embed_aa, a.id FROM test.tb_aa AS a JOIN test.tb_aa AS b ON a.id=b.id`

	mds := TbAas{}
	if err := tx.Tx.Select(&mds, q); err != nil {
		t.Error(err)
		return
	}
	if err := ModelEqual(append(TbAas{}, mda1, mda2), mds); err != nil {
		t.Error(err)
		return
	}

	mds1 := TbAcs{}
	if err := tx.Tx.Select(&mds1, q); err != nil {
		t.Error(err)
		return
	}

	mdc1 := NewTbAc(mda1.EmbedAa, mda1.EmbedAb, mda1.Note)
	mdc1.A = mda1.ID
	mdc2 := NewTbAc(mda2.EmbedAa, mda2.EmbedAb, mda2.Note)
	mdc2.A = mda2.ID
	if err := ModelEqual(append(TbAcs{}, mdc1, mdc2), mds1); err != nil {
		t.Error(err)
		return
	}

	mds2 := Tbs{}
	if err := tx.Tx.Select(&mds2, q); err == errors.New("missing destination name note in *pgsql.Tbs") {
		t.Error(err)
		return
	}
}
