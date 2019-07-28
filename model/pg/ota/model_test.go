package model

import (
	"cloud/lib/logger"
	"cloud/server/ota/config"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-pg/pg"
	"github.com/google/go-cmp/cmp"
	"cloud/lib/null"
)

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	fmt.Println(q.FormattedQuery())
	return c, nil
}

func TestBasicCRUD(t *testing.T) {

	config.ParseConfig(os.Getenv("GOPATH")+"/src/cloud/server/ota/config/config.local.json", &config.Config)
	cfg := &config.Config

	pgdb, err := cfg.DB.Pg().Connect()
	if err != nil {
		return
	}
	pgdb.AddQueryHook(dbLogger{})

	tx, err := pgdb.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer tx.Rollback()

	var mds []MsFw
	if err := tx.Model(&mds).Select(); err != nil {
		t.Error(err)
		return
	}

	// insert
	md1 := NewMsFw("v1", 1, "bk1", "obj1", time.Now().UTC(), "tag1")
	md2 := NewMsFw("v2", 2, "bk2", "obj2", time.Now().UTC(), "tag2")
	if err := tx.Insert(&md1); err != nil {
		t.Error(err)
		return
	}
	if err := tx.Insert(&md2); err != nil {
		t.Error(err)
		return
	}

	mds1 := append(mds, md1, md2)
	if err := ModelsEqual(tx, mds1); err != nil {
		t.Error(err)
		return
	}

	// count
	n, err := tx.Model(&MsFw{}).Count()
	if err != nil {
		t.Error(err)
		return
	}
	if n != len(mds1) {
		t.Error(n)
		return
	}

	var md MsFw
	md.ID = md1.ID

	// select by primary key
	if err := tx.Select(&md); err != nil {
		t.Error(err)
		return
	}
	if err := ModelEqual(md1, md); err != nil {
		t.Error(err)
		return
	}

	// select by id
	md = MsFw{}
	if err := tx.Model(&md).Where("id=?", md1.ID).Select(); err != nil {
		t.Error(err)
		return
	}
	if err := ModelEqual(md1, md); err != nil {
		t.Error(err)
		return
	}

	// update by primary key
	md2.Version = "test22"
	md2.Com = 3
	if err := tx.Update(&md2); err != nil {
		t.Error(err)
		return
	}

	// delete by primary key
	if err := tx.Delete(&md1); err != nil {
		t.Error(err)
		return
	}

	if err := ModelsEqual(tx, append(mds, md2)); err != nil {
		t.Error(err)
		return
	}

}

func TestWhere(t *testing.T) {

	config.ParseConfig(os.Getenv("GOPATH")+"/src/cloud/server/ota/config/config.local.json", &config.Config)
	cfg := &config.Config
	pgdb, err := cfg.DB.Pg().Connect()
	if err != nil {
		return
	}
	pgdb.AddQueryHook(dbLogger{})

	tx, err := pgdb.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	defer tx.Rollback()

	var mds MsFws
	if err := tx.Model(&mds).Select(); err != nil {
		t.Error(err)
		return
	}

	// insert
	md1 := NewMsFw("v1", 1, "bk1", "obj1", time.Now().UTC(), "tag1")
	md2 := NewMsFw("v2", 2, "bk2", "obj2", time.Now().UTC(), "tag2")
	md21 := NewMsFw("v21", 2, "bk2", "obj2", time.Now().UTC(), "tag2")
	if err := tx.Insert(&md1); err != nil {
		t.Error(err)
		return
	}
	if err := tx.Insert(&md2); err != nil {
		t.Error(err)
		return
	}
	if err := tx.Insert(&md21); err != nil {
		t.Error(err)
		return
	}

	if err := ModelsEqual(tx, append(mds, md1, md2, md21)); err != nil {
		t.Error(err)
		return
	}

	// count
	n, err := tx.Model(&MsFw{}).Where("com=? AND version=?", md2.Com, md2.Version).Count()
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
	if err := tx.Model(&md).Where("com=? AND version=?", md21.Com, md21.Version).Select(); err != nil {
		t.Error(err)
		return
	}
	if err := ModelEqual(md, md21); err != nil {
		t.Error(err)
		return
	}

	// no update
	if _, err := tx.Model(&MsFw{ID: md2.ID, Tag:  null.NewString("fake")}).Where("version = ?", "fake").Update(); err != nil {
		t.Error(err)
		return
	}

	if err := ModelsEqual(tx, append(mds, md1, md2, md21)); err != nil {
		t.Error(err)
		return
	}

	// update all field except for id
	md22 := md21
	md22.ID = md21.ID + 1
	md22.Tag = null.NewString("")
	if _, err := tx.Model(&md22).Where("version = ?", md21.Version).Update(); err != nil {
		t.Error(err)
		return
	}
	md22.ID = md21.ID
	if err := ModelsEqual(tx, append(mds, md1, md2, md22)); err != nil {
		t.Error(err)
		return
	}


	md22.Version = ""
	if _, err := tx.Model(&md22).Where("version = ?", md21.Version).Update(); err != nil {
		t.Error(err)
		return
	}

	// delete by tag
	if _, err := tx.Model(&MsFw{}).Where("tag = ?", md22.Tag).Delete(); err != nil {
		t.Error(err)
		return
	}

	if err := ModelsEqual(tx, append(mds, md1, md2)); err != nil {
		t.Error(err)
		return
	}
}

func ModelsEqual(tx *pg.Tx, want MsFws) error {
	var got MsFws
	if err := tx.Model(&got).Select(); err != nil {
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
