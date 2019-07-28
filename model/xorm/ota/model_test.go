package model

import (
	"cloud/lib/logger"
	"cloud/server/ota/config"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/google/go-cmp/cmp"
	_ "github.com/lib/pq"
)

func TestBasicCRUD(t *testing.T) {

	config.ParseConfig(os.Getenv("GOPATH")+"/src/cloud/server/ota/config/config.local.json", &config.Config)
	cfg := &config.Config
	dbConfig := config.GetPgsqlConfig(cfg.DB)

	engine, err := xorm.NewEngine("postgres", dbConfig)
	if err != nil {
		t.Error(err)
	}

	pg := NewPg(engine)

	tx := getTx(pg.Begin())
	defer tx.Rollback()

	mds, err := GetAllMsFws(tx)
	if err != nil {
		t.Error(err)
		return
	}

	// insert
	md1 := NewMsFw("v1", 1, "bk1", "obj1", time.Now().UTC(), "tag1")
	md2 := NewMsFw("v2", 2, "bk2", "obj2", time.Now().UTC(), "tag2")
	if err := InsertMsFw(tx, &md1, &md2); err != nil {
		t.Error(err)
		return
	}

	mds1 := append(mds, md1, md2)
	if err := ModelsEqual(tx, mds1); err != nil {
		t.Error(err)
		return
	}

	// count
	n, err := tx.Count(MsFw{})
	if err != nil {
		t.Error(err)
		return
	}
	if n != int64(len(mds1)) {
		t.Error(n)
		return
	}

	var md MsFw

	// select by id
	md, err = GetMsFwById(tx, md1.ID)
	if err != nil {
		t.Error(err)
		return
	}
	if err := ModelEqual(md1, md); err != nil {
		t.Error(err)
		return
	}

	// update by id
	md2.Version = "test22"
	md2.Com = 3
	if err := UpdateMsFwById(tx, md2); err != nil {
		t.Error(err)
		return
	}

	// delete
	if err := DeleteMsFw(tx, md1); err != nil {
		t.Error(err)
		return
	}

	if err := ModelsEqual(tx, append(mds, md2)); err != nil {
		t.Error(err)
		return
	}

	// delete by id
	if err := DeleteMsFwById(tx, md2.ID); err != nil {
		t.Error(err)
		return
	}

	if err := ModelsEqual(tx, mds); err != nil {
		t.Error(err)
		return
	}

}

func TestWhere(t *testing.T) {

	config.ParseConfig(os.Getenv("GOPATH")+"/src/cloud/server/ota/config/config.local.json", &config.Config)
	cfg := &config.Config
	dbConfig := config.GetPgsqlConfig(cfg.DB)

	engine, err := xorm.NewEngine("postgres", dbConfig)
	if err != nil {
		t.Error(err)
	}
	engine.ShowSQL(true)

	pg := NewPg(engine)

	tx := getTx(pg.Begin())
	defer tx.Rollback()

	mds, err := GetAllMsFws(tx)
	if err != nil {
		t.Error(err)
		return
	}

	// insert
	md1 := NewMsFw("v1", 1, "bk1", "obj1", time.Now().UTC(), "tag1")
	md2 := NewMsFw("v2", 2, "bk2", "obj2", time.Now().UTC(), "tag2")
	md21 := NewMsFw("v21", 2, "bk2", "obj2", time.Now().UTC(), "tag2")
	if err := InsertMsFw(tx, &md1, &md2, &md21); err != nil {
		t.Error(err)
		return
	}

	if err := ModelsEqual(tx, append(mds, md1, md2, md21)); err != nil {
		t.Error(err)
		return
	}

	// count
	n, err := CountMsFwByUnique(tx, md2.MsFwUnique)
	if err != nil {
		t.Error(err)
		return
	}
	if n != 1 {
		t.Error(n)
		return
	}

	var md MsFw

	// select by com and version
	md = MsFw{}
	if _, err := tx.Where("com=? AND version = ?", md21.Com, md21.Version).Get(&md); err != nil {
		t.Error(err)
		return
	}
	if err := ModelEqual(md, md21); err != nil {
		t.Error(err)
		return
	}

	// no update
	if _, err := tx.Where("version = ?", "fake").Update(MsFw{ID: md2.ID, Tag: "fake"}); err != nil {
		t.Error(err)
		return
	}

	if err := ModelsEqual(tx, append(mds, md1, md2, md21)); err != nil {
		t.Error(err)
		return
	}

	// update non-zero field by version
	md2.Version = "test2"
	md2.Com = 2
	md22 := md2
	md22.ID = 0
	md22.Tag = ""
	if _, err := tx.Where("version = ?", "v2").Update(md22); err != nil {
		t.Error(err)
		return
	}

	if err := ModelsEqual(tx, append(mds, md1, md2, md21)); err != nil {
		t.Error(err)
		return
	}

	// update all field except for id
	md22 = md21
	md22.ID = md21.ID + 1
	md22.Tag = ""
	if _, err := tx.Where("version = ?", md21.Version).AllCols().Update(md22); err != nil {
		t.Error(err)
		return
	}
	md22.ID = md21.ID
	if err := ModelsEqual(tx, append(mds, md1, md2, md22)); err != nil {
		t.Error(err)
		return
	}

	// update all field
	md22.ID = md21.ID + 1
	if _, err := tx.Where("version = ?", md21.Version).Cols(md22.Cols()...).Update(md22); err != nil {
		t.Error(err)
		return
	}

	if err := ModelsEqual(tx, append(mds, md1, md2, md22)); err != nil {
		t.Error(err)
		return
	}

	// delete by tag
	if _, err := tx.Where("tag = ?", md22.Tag).Cols(md22.Cols()...).Delete(MsFw{}); err != nil {
		t.Error(err)
		return
	}

	if err := ModelsEqual(tx, append(mds, md1, md2)); err != nil {
		t.Error(err)
		return
	}

	// delete by non-zero field
	if _, err := tx.Delete(MsFw{Tag: md2.Tag}); err != nil {
		t.Error(err)
		return
	}

	if err := ModelsEqual(tx, append(mds, md1)); err != nil {
		t.Error(err)
		return
	}

}

func ModelsEqual(tx *Tx, want MsFws) error {
	var got MsFws
	if err := tx.Asc("id").Find(&got); err != nil {
		logger.Error(err)
		return err
	}

	if len(want) != len(got) {
		return fmt.Errorf("\nwant:%+v,\ngot :%+v", want, got)
	}

	for i := range want {
		if want[i].Time.Unix() != got[i].Time.Unix() {
			return fmt.Errorf("\nwant:%+v,\ngot :%+v", want, got)
		}
		got[i].Time = want[i].Time
	}

	if !cmp.Equal(want, got) {
		return fmt.Errorf("\nwant:%+v,\ngot :%+v", want, got)
	}
	return nil

}

func ModelEqual(want, got MsFw) error {

	if want.Time.Unix() != got.Time.Unix() {
		return fmt.Errorf("\nwant:%+v,\ngot :%+v", want, got)
	}
	got.Time = want.Time

	if !cmp.Equal(want, got) {
		return fmt.Errorf("\nwant:%+v,\ngot :%+v", want, got)
	}
	return nil

}
