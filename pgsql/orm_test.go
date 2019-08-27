package pgsql

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type Tb struct {
	TableName struct{} `json:"-" db:"test.tb_aa"`
	A         int64    `json:"id" db:"id,pk,serial"`
	B         int      `json:"embed_aa" db:"embed_aa"`
	C         string   `json:"embed_ab" db:"embed_ab"`
	D         string   `json:"note" db:"note"`
}

type Tbs []Tb

type TbAas []TbAa

type TbAa struct {
	ID int64 `json:"id" db:",pk,serial"`
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

var (
	mda1 = NewTbAa(1, "b1", "n1")
	mda2 = NewTbAa(2, "b2", "n2")
	mdas = TbAas{mda1, mda2}

	md1 = NewTb(1, "b1", "n1")
	md2 = NewTb(2, "b2", "n2")
	mds = Tbs{md1, md2}

	mdb1 = NewTbAb(1, "pkb1", 1, "b1", "n1")
	mdb2 = NewTbAb(2, "pkb2", 2, "b2", "n2")
	mdbs = TbAbs{mdb1, mdb2}
)

func TestGetTableInfo(t *testing.T) {

	pks1 := []string{"id"}
	pks2 := []string{"pka", "pkb"}
	cols1 := []string{"id", "embed_aa", "embed_ab", "note"}
	cols2 := []string{"pka", "pkb", "embed_aa", "embed_ab", "note"}
	serial1 := "id"
	serial2 := ""
	f := fmt.Sprintf

	for _, v := range []struct {
		name   string
		md     interface{}
		tb     string
		pks    []string
		cols   []string
		serial string
	}{
		{"Tb", md1, "test.tb_aa", pks1, cols1, serial1},
		{"Tbs", mds, "test.tb_aa", pks1, cols1, serial1},
		{"TbAa", mda1, "tb_aa", pks1, cols1, serial1},
		{"TbAas", mdas, "tb_aa", pks1, cols1, serial1},
		{"TbAb", mdb1, "tb_ab", pks2, cols2, serial2},
		{"TbAbs", mdbs, "tb_ab", pks2, cols2, serial2},
	} {
		t.Run(f(`getTable=%s`, v.name), testGetTable(v.md, v.tb))
		t.Run(f(`getPks=%s`, v.name), testGetPks(v.md, v.pks))
		t.Run(f(`getCols=%s`, v.name), testGetCols(v.md, v.cols))
		t.Run(f(`getSerial=%s`, v.name), testGetSerial(v.md, v.serial))
	}
}

func testGetTable(md interface{}, want string) func(t *testing.T) {
	f := func(t *testing.T) {
		if got := GetTable(md); got != want {
			t.Errorf(`want:%s, got:%s`, want, got)
		}
	}
	return f
}

func testGetPks(md interface{}, want []string) func(t *testing.T) {
	f := func(t *testing.T) {
		if got := GetPks(md); !cmp.Equal(got, want) {
			t.Errorf(`want:%v, got:%v`, want, got)
		}
	}
	return f
}

func testGetCols(md interface{}, want []string) func(t *testing.T) {
	f := func(t *testing.T) {
		if got := GetCols(md); !cmp.Equal(got, want) {
			t.Errorf(`want:%v, got:%v`, want, got)
		}
	}
	return f
}

func testGetSerial(md interface{}, want string) func(t *testing.T) {
	f := func(t *testing.T) {
		if got := GetSerial(md); !cmp.Equal(got, want) {
			t.Errorf(`want:%v, got:%v`, want, got)
		}
	}
	return f
}
