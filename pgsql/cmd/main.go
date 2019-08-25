package main

import (
	"cloud/lib/logger"
	"fmt"

	"test/pgsql"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Tb struct {
	TableName struct{} `json:"-" db:"test.tb_aa"`
	A         uint64   `json:"id" db:"id,pk"`
	B         int      `json:"embed_aa" db:"embed_aa"`
	C         string   `json:"embed_ab" db:"embed_ab"`
	D         string   `json:"note" db:"note"`
}

type Tbs []Tb

type TbAas []TbAa

type TbAa struct {
	ID uint64 `json:"id" db:",pk"`
	Embed
	Note string `json:"note"`
}

type TbAbs []TbAb

type TbAb struct {
	Pka int    `json:"pka" db:",pk"`
	Pkb string `json:"pkb" db:",pk"`
	Embed
	Note string `json:"note"`
}

type Embed struct {
	EmbedAa int    `json:"embed_aa"`
	EmbedAb string `json:"embed_ab"`
}

func NewTb(b int, c, d string) Tb {

	return Tb{
		B: b,
		C: c,
		D: d,
	}
}

func NewTbAa(a int, b, note string) TbAa {

	u := Embed{EmbedAa: a, EmbedAb: b}
	return TbAa{
		Embed: u,
		Note:  note,
	}
}

func NewTbAb(pka int, pkb string, a int, b, note string) TbAb {

	u := Embed{EmbedAa: a, EmbedAb: b}
	return TbAb{
		Pka:   pka,
		Pkb:   pkb,
		Embed: u,
		Note:  note,
	}
}

type info struct {
	tb   string
	pks  []string
	cols []string
	vals []interface{}
}

func getInfo(md interface{}) info {
	return info{
		tb:   pgsql.GetTable(md),
		pks:  pgsql.GetPks(md),
		cols: pgsql.GetCols(md),
		vals: pgsql.GetValues(md),
	}
}

func main() {

	dbConfig := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5431, "test", "abcd", "test")

	db, _ := pgsql.Connect(dbConfig)
	defer db.Close()

	mda1 := NewTbAa(1, "b1", "n1")
	mda2 := NewTbAa(2, "b2", "n2")
	mdas := TbAas{mda1, mda2}

	md1 := NewTb(1, "b1", "n1")
	md2 := NewTb(2, "b2", "n2")
	mds := Tbs{md1, md2}

	mdb1 := NewTbAb(1, "pkb1", 1, "b1", "n1")
	mdb2 := NewTbAb(2, "pkb2", 2, "b2", "n2")
	mdbs := TbAbs{mdb1, mdb2}

	logger.Debug(getInfo(md1))
	logger.Debug(getInfo(mda1))
	logger.Debug(getInfo(mdb1))
	logger.Debug(getInfo(mds))
	logger.Debug(getInfo(mdas))
	logger.Debug(getInfo(mdbs))

	db.Select(mda1).Where("com=? AND gp IN(?)", 1, pq.Array([]uint64{1, 2, 3})).Order("com DESC", "gp ASC").Limit(10).SQL()
	db.Select(&mda1).Where("com=? AND gp IN(?)", 1, pq.Array([]uint64{1, 2, 3})).Order("com DESC", "gp ASC").Limit(10).SQL()
	db.Select(mdas).Where("com=? AND gp IN(?)", 1, pq.Array([]uint64{1, 2, 3})).Order("com DESC", "gp ASC").Limit(10).SQL()
	db.Select(&mdas).Where("com=? AND gp IN(?)", 1, pq.Array([]uint64{1, 2, 3})).Order("com DESC", "gp ASC").Limit(10).SQL()
	db.Select(nil).Where("com=? AND gp IN(?)", 1, pq.Array([]uint64{1, 2, 3})).Order("com DESC", "gp ASC").Limit(10).SQL()

	logger.Debug(pgsql.GetValues(mda1))
	db.Insert(mda1).SQL()

	db.Update(mda1).SQL()
	db.Update(mda1).Set("embed_aa=?, embed_ab=?, note=?", mda1.EmbedAa, mda1.EmbedAb, mda1.Note).SQL()
	db.Delete(mda1).SQL()

	var ch chan int
	getInfo(ch)

}
