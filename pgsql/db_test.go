package pgsql

import (
	"cloud/lib/logger"
	"database/sql"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	_ "github.com/lib/pq"
)

func TestBasicCRUD(t *testing.T) {

	dbConfig := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5431, "test", "abcd", "test")
	db, err := Connect(dbConfig)
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()

	var mdas TbAas
	db.Select(&mdas).SQL()
	if err := db.Select(&mdas).Run(); err != nil {
		t.Error(err)
		return
	}
	// len(mdas) should be 0
	if err := ModelsEqual(db, mdas); err != nil {
		t.Error(err)
		return
	}

	defer db.Truncate(mdas).Run()

	if err := db.Insert(mda1).Returning("id").Run().Scan(&mda1.ID); err != nil {
		t.Error(err)
		return
	}
	if err := db.Insert(mda2).Returning("id").Run().Scan(&mda2.ID); err != nil {
		t.Error(err)
		return
	}

	md := TbAa{ID: mda1.ID}
	if err := db.SelectByPk(&md).Run(); err != nil {
		t.Error(err)
		return
	}
	if err := ModelEqual(mda1, md); err != nil {
		t.Error(err)
		return
	}

	if err := ModelsEqual(db, append(mdas, mda1, mda2)); err != nil {
		t.Error(err)
		return
	}

	mda2.EmbedAa = mda2.EmbedAa * 10
	mda2.EmbedAb = mda2.EmbedAb + mda2.EmbedAb
	mda2.Note = mda2.Note + mda2.Note
	if _, err := db.UpdateByPk(mda2).Run(); err != nil {
		t.Error(err)
		return
	}

	if err := ModelsEqual(db, append(mdas, mda1, mda2)); err != nil {
		t.Error(err)
		return
	}

	var n int
	if err := db.Count(mdas).Run().Scan(&n); err != nil {
		t.Error(err)
		return
	}
	if n != 2 {
		t.Error(n)
	}

	n = 0
	if err := db.CountDistinct(mdas, "note").Run().Scan(&n); err != nil {
		t.Error(err)
		return
	}
	if n != 2 {
		t.Error(n)
	}

	if _, err := db.Delete(mda2).Run(); err != ErrDeleteWithoutCondition {
		t.Error(err)
		return
	}

	if _, err := db.DeleteByPk(mda2).Run(); err != nil {
		t.Error(err)
		return
	}
	if err := ModelsEqual(db, append(mdas, mda1)); err != nil {
		t.Error(err)
		return
	}

	if _, err := db.Truncate(mdas).Run(); err != nil {
		t.Error(err)
		return
	}

	db.Select(&mdas).SQL()
	if err := db.Select(&mdas).Run(); err != nil {
		t.Error(err)
		return
	}
	// len(mdas) should be 0
	if err := ModelsEqual(db, mdas); err != nil {
		t.Error(err)
		return
	}

}

func TestWhere(t *testing.T) {

	dbConfig := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5431, "test", "abcd", "test")
	db, err := Connect(dbConfig)
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()

	var mdas TbAas
	db.Select(&mdas).SQL()
	if err := db.Select(&mdas).Run(); err != nil {
		t.Error(err)
		return
	}
	// len(mdas) should be 0
	if err := ModelsEqual(db, mdas); err != nil {
		t.Error(err)
		return
	}

	defer db.Truncate(mdas).Run()

	if err := db.Insert(mda1).Returning("id").Run().Scan(&mda1.ID); err != nil {
		t.Error(err)
		return
	}
	if err := db.Insert(mda2).Returning("id").Run().Scan(&mda2.ID); err != nil {
		t.Error(err)
		return
	}

	if err := ModelsEqual(db, append(mdas, mda1, mda2)); err != nil {
		t.Error(err)
		return
	}

	md := TbAa{}
	if err := db.Select(&md).Where("note=?", mda1.Note).Run(); err != nil {
		t.Error(err)
		return
	}
	if err := ModelEqual(mda1, md); err != nil {
		t.Error(err)
		return
	}

	mda2.EmbedAa = mda2.EmbedAa * 10
	mda2.EmbedAb = mda2.EmbedAb + mda2.EmbedAb
	if _, err := db.Update(TbAa{}).Set("embed_aa=?, embed_ab=?", mda2.EmbedAa, mda2.EmbedAb).Where("note=?", mda2.Note).Run(); err != nil {
		t.Error(err)
		return
	}

	if err := ModelsEqual(db, append(mdas, mda1, mda2)); err != nil {
		t.Error(err)
		return
	}

	var n int
	if err := db.Count(TbAa{}).Where("note=?", mda2.Note).Run().Scan(&n); err != nil {
		t.Error(err)
		return
	}
	if n != 1 {
		t.Error(n)
	}

	n = 0
	if err := db.CountDistinct(TbAa{}, "note").Where("note=?", mda2.Note).Run().Scan(&n); err != nil {
		t.Error(err)
		return
	}
	if n != 1 {
		t.Error(n)
	}

	if _, err := db.Delete(TbAa{}).Where("note=?", mda2.Note).Run(); err != nil {
		t.Error(err)
		return
	}

	if err := ModelsEqual(db, append(mdas, mda1)); err != nil {
		t.Error(err)
		return
	}

}

type database interface {
	Select(md interface{}) Query
}

func ModelsEqual(db database, want TbAas) error {
	var got TbAas
	if err := db.Select(&got).Run(); err != nil {
		logger.Error(err)
		return err
	}

	if len(want) != len(got) {
		return fmt.Errorf("\nwant:%+v,\ngot :%+v", want, got)
	}

	if !cmp.Equal(want, got) {
		return fmt.Errorf("\nwant:%+v,\ngot :%+v", want, got)
	}
	return nil

}

func ModelEqual(want, got TbAa) error {

	if !cmp.Equal(want, got) {
		return fmt.Errorf("\nwant:%+v,\ngot :%+v", want, got)
	}
	return nil

}

func Benchmark01pgsqlDB(b *testing.B) {
	db := setupBenchData()
	defer db.Close()
	defer db.Truncate(mdas).Run()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.Select(&TbAas{}).Run()
	}
}

func Benchmark02sqlxDB(b *testing.B) {
	db := setupBenchData()
	defer db.Close()
	defer db.Truncate(mdas).Run()
	sqlxDB := db.DB

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sqlxDB.Select(&TbAas{}, "SELECT id, embed_aa, embed_ab, note FROM tb_aa")
	}
}

func Benchmark03sqlxDB(b *testing.B) {
	db := setupBenchData()
	defer db.Close()
	defer db.Truncate(mdas).Run()
	q := db.Select(&TbAas{}).SQL()
	sqlxDB := db.DB

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sqlxDB.Select(&TbAas{}, q)
	}
}

func Benchmark04sqlDB(b *testing.B) {
	db := setupBenchData()
	defer db.Close()
	defer db.Truncate(mdas).Run()
	sqlDB := db.DB.DB

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sqlSelect(sqlDB)
	}
}

func setupBenchData() *DB {
	dbConfig := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5431, "test", "abcd", "test")
	db, err := Connect(dbConfig)
	if err != nil {
		logger.Panic(err)
	}
	for i := 0; i < 1000; i++ {
		if _, err := db.Insert(mda1).Run(); err != nil {
			logger.Panic(err)
		}
	}
	return db
}

func sqlSelect(db *sql.DB) {
	rows, err := db.Query("SELECT id, embed_aa, embed_ab, note FROM tb_aa")
	if err != nil {
		panic(err)
	}
	var mds TbAas
	for rows.Next() {
		var md TbAa
		if err := rows.Scan(&md.ID, &md.EmbedAa, &md.EmbedAb, &md.Note); err != nil {
			panic(err)
		}
		mds = append(mds, md)
	}
}
