package pgsql

import (
	"cloud/lib/logger"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lib/pq"
)

func TestBasicCRUD(t *testing.T) {

	dbConfig := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5433, "test", "abcd", "test")
	db, err := Connect(dbConfig)
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()
	test := NewTester(t)

	var mdas TbAas
	test.Run(db.Select(&mdas))
	// len(mdas) should be 0
	if err := ModelsEqual(db, mdas); err != nil {
		t.Error(err)
		return
	}

	defer db.Truncate(mdas).Run()
	mda1 := NewTbAa(1, "b1", "n1")
	mda2 := NewTbAa(2, "b2", "n2")
	test.Scan(db.Insert(mda1).Returning("id").Run(), &mda1.ID)
	test.Scan(db.Insert(mda2).Returning("id").Run(), &mda2.ID)

	md := TbAa{ID: mda1.ID}
	test.Run(db.SelectByPk(&md))
	if err := ModelEqual(mda1, md); err != nil {
		t.Error(err)
		return
	}

	md = TbAa{ID: mda1.ID}
	test.Run(db.SelectByPk(&md, "embed_aa", "embed_ab"))
	if err := ModelEqual(TbAa{ID: mda1.ID, Embed: Embed{EmbedAa: 1, EmbedAb: "b1"}, Note: ""}, md); err != nil {
		t.Error(err)
		return
	}

	// column order dosen't matter
	md = TbAa{}
	test.Run(db.Select(&md, "note", "embed_ab", "embed_aa").Where("id=?", mda1.ID))
	if err := ModelEqual(TbAa{Embed: Embed{EmbedAa: 1, EmbedAb: "b1"}, Note: "n1"}, md); err != nil {
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
	test.Run(db.UpdateByPk(mda2))
	if err := ModelsEqual(db, append(mdas, mda1, mda2)); err != nil {
		t.Error(err)
		return
	}

	var n int
	test.Scan(db.Count(mdas).Run(), &n)
	if n != 2 {
		t.Error(n)
	}

	n = 0
	test.Scan(db.CountDistinct(mdas, "note").Run(), &n)
	if n != 2 {
		t.Error(n)
		return
	}

	if err := db.Delete(mda2).Run(); err != ErrDeleteWithoutCondition {
		t.Error(err)
		return
	}

	test.Run(db.DeleteByPk(mda2))
	if err := ModelsEqual(db, append(mdas, mda1)); err != nil {
		t.Error(err)
		return
	}

	test.Run(db.Truncate(mdas))
	test.Run(db.Select(&mdas))

	// len(mdas) should be 0
	if err := ModelsEqual(db, mdas); err != nil {
		t.Error(err)
		return
	}

}

func TestWhere(t *testing.T) {

	dbConfig := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5433, "test", "abcd", "test")
	db, err := Connect(dbConfig)
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()

	test := NewTester(t)

	var mdas TbAas
	test.Run(db.Select(&mdas))

	// len(mdas) should be 0
	if err := ModelsEqual(db, mdas); err != nil {
		t.Error(err)
		return
	}

	defer db.Truncate(mdas).Run()
	mda1 := NewTbAa(1, "b1", "n1")
	mda2 := NewTbAa(2, "b2", "n2")
	test.Scan(db.Insert(mda1).Returning("id").Run(), &mda1.ID)
	test.Scan(db.Insert(mda2).Returning("id").Run(), &mda2.ID)
	if err := ModelsEqual(db, append(mdas, mda1, mda2)); err != nil {
		t.Error(err)
		return
	}

	md := TbAa{}
	test.Run(db.Select(&md).Where("note=?", mda1.Note))
	if err := ModelEqual(mda1, md); err != nil {
		t.Error(err)
		return
	}

	mda2.EmbedAa = mda2.EmbedAa * 10
	mda2.EmbedAb = mda2.EmbedAb + mda2.EmbedAb
	test.Run(db.Update(TbAa{}).Set("embed_aa=?, embed_ab=?", mda2.EmbedAa, mda2.EmbedAb).Where("note=?", mda2.Note))
	if err := ModelsEqual(db, append(mdas, mda1, mda2)); err != nil {
		t.Error(err)
		return
	}

	var n int
	test.Scan(db.Count(TbAa{}).Where("note=?", mda2.Note).Run(), &n)
	if n != 1 {
		t.Error(n)
		return
	}

	n = 0
	test.Scan(db.CountDistinct(TbAa{}, "note").Where("note=?", mda2.Note).Run(), &n)
	if n != 1 {
		t.Error(n)
		return
	}

	test.Run(db.Delete(TbAa{}).Where("note=?", mda2.Note))
	if err := ModelsEqual(db, append(mdas, mda1)); err != nil {
		t.Error(err)
		return
	}

}

func TestWhereIn(t *testing.T) {

	dbConfig := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5433, "test", "abcd", "test")
	db, err := Connect(dbConfig)
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()

	test := NewTester(t)

	var mdas TbAas
	test.Run(db.Select(&mdas))

	// len(mdas) should be 0
	if err := ModelsEqual(db, mdas); err != nil {
		t.Error(err)
		return
	}

	defer db.Truncate(mdas).Run()
	mda1 := NewTbAa(1, "b1", "n1")
	mda2 := NewTbAa(2, "b2", "n2")
	test.Scan(db.Insert(mda1).Returning("id").Run(), &mda1.ID)
	test.Scan(db.Insert(mda2).Returning("id").Run(), &mda2.ID)
	if err := ModelsEqual(db, append(mdas, mda1, mda2)); err != nil {
		t.Error(err)
		return
	}

	mds := TbAas{}
	test.Run(db.Select(&mds).Where("note IN (?)", pq.Array([]string{mda1.Note, mda2.Note})))
	if err := ModelEqual(TbAas{mda1, mda2}, mds); err != nil {
		t.Error(err)
		return
	}
	mds = TbAas{}
	test.Run(db.Select(&mds).Where("note IN (?)", pq.Array([]string{mda2.Note})))
	if err := ModelEqual(TbAas{mda2}, mds); err != nil {
		t.Error(err)
		return
	}

	mds = TbAas{}
	test.Run(db.Select(&mds).Where("note").In(pq.Array([]string{mda1.Note, mda2.Note})))
	if err := ModelEqual(TbAas{mda1, mda2}, mds); err != nil {
		t.Error(err)
		return
	}
	mds = TbAas{}
	test.Run(db.Select(&mds).Where("note").In(pq.Array([]string{mda2.Note})))
	if err := ModelEqual(TbAas{mda2}, mds); err != nil {
		t.Error(err)
		return
	}

}

func TestJoin(t *testing.T) {

	dbConfig := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5433, "test", "abcd", "test")
	db, err := Connect(dbConfig)
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()

	test := NewTester(t)

	var mdas TbAas
	test.Run(db.Select(&mdas))

	// len(mdas) should be 0
	if err := ModelsEqual(db, mdas); err != nil {
		t.Error(err)
		return
	}

	defer db.Truncate(mdas).Run()
	mda1 := NewTbAa(1, "b1", "n1")
	mda2 := NewTbAa(2, "b2", "n2")
	test.Scan(db.Insert(mda1).Returning("id").Run(), &mda1.ID)
	test.Scan(db.Insert(mda2).Returning("id").Run(), &mda2.ID)

	q := `SELECT a.note, a.embed_ab, a.embed_aa, a.id FROM test.tb_aa AS a JOIN test.tb_aa AS b ON a.id=b.id`

	mds := TbAas{}
	if err := db.DB.Select(&mds, q); err != nil {
		t.Error(err)
		return
	}
	if err := ModelEqual(append(TbAas{}, mda1, mda2), mds); err != nil {
		t.Error(err)
		return
	}

	mds1 := TbAcs{}
	if err := db.DB.Select(&mds1, q); err != nil {
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
	if err := db.DB.Select(&mds2, q); err == errors.New("missing destination name note in *pgsql.Tbs") {
		t.Error(err)
		return
	}
}

type database interface {
	Select(md interface{}, cols ...string) Query
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

func ModelEqual(want, got interface{}) error {

	if !cmp.Equal(want, got) {
		return fmt.Errorf("\nwant:%+v,\ngot :%+v", want, got)
	}
	return nil

}

func BenchmarkDB(b *testing.B) {
	db := setupBenchData()
	defer db.Close()
	defer db.Truncate(mdas).Run()

	b.Run("pgsqlDB-01", benchmark01pgsqlDB(b, db))
	b.Run("sqlxDB-01", benchmark01sqlxDB(b, db))
	b.Run("sqlxDB-02", benchmark02sqlxDB(b, db))
	b.Run("sqlDB-01", benchmark01sqlDB(b, db))

}

func benchmark01pgsqlDB(b *testing.B, db *DB) func(b *testing.B) {

	f := func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			db.Select(&TbAas{}).Run()
		}
	}
	return f
}

func benchmark01sqlxDB(b *testing.B, db *DB) func(b *testing.B) {
	sqlxDB := db.DB
	f := func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sqlxDB.Select(&TbAas{}, "SELECT id, embed_aa, embed_ab, note FROM tb_aa")
		}
	}
	return f
}

func benchmark02sqlxDB(b *testing.B, db *DB) func(b *testing.B) {
	q := db.Select(&TbAas{}).SQL()
	sqlxDB := db.DB
	f := func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sqlxDB.Select(&TbAas{}, q)
		}
	}
	return f
}

func benchmark01sqlDB(b *testing.B, db *DB) func(b *testing.B) {
	sqlDB := db.DB.DB
	f := func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sqlSelect(sqlDB)
		}
	}
	return f
}

func setupBenchData() *DB {
	dbConfig := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5433, "test", "abcd", "test")
	db, err := Connect(dbConfig)
	if err != nil {
		logger.Panic(err)
	}
	for i := 0; i < 1000; i++ {
		if err := db.Insert(mda1).Run(); err != nil {
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
