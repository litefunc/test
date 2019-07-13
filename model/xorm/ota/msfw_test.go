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

func TestMsFw(t *testing.T) {

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
	logger.Debug(mds)

	// insert
	md1 := NewMsFw("v1", 1, "bk1", "obj1", time.Now().UTC(), "tag1")
	md2 := NewMsFw("v2", 2, "bk2", "obj2", time.Now().UTC(), "tag2")
	if err := InsertMsFw(tx, &md1, &md2); err != nil {
		t.Error(err)
		return
	}

	if err := MsFwsEqual(tx, append(mds, md1, md2)); err != nil {
		t.Error(err)
		return
	}

	// select by id
	md11, err := GetMsFwById(tx, md1.ID)
	if err != nil {
		t.Error(err)
		return
	}
	if err := MsFwEqual(md1, md11); err != nil {
		t.Error(err)
		return
	}

	// select by version
	// md11 = MsFw{}
	// if _, err := tx.Where("version = ?", md1.Version).Get(&md11); err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// if err := MsFwEqual(md1, md11); err != nil {
	// 	t.Error(err)
	// 	return
	// }

	// update by id
	md2.Version = "test22"
	md2.Com = 3
	if err := UpdateMsFw(tx, md2); err != nil {
		t.Error(err)
		return
	}

	// update by version
	// md2.Version = "test2"
	// md2.Com = 2
	// if _, err := tx.Where("version = ?", md2.Version).Update(md2); err != nil {
	// 	t.Error(err)
	// 	return
	// }

	// if err := MsFwsEqual(tx, append(mds, md1, md2)); err != nil {
	// 	t.Error(err)
	// 	return
	// }

	// delete
	if err := DeleteMsFw(tx, md1); err != nil {
		t.Error(err)
		return
	}

	if err := MsFwsEqual(tx, append(mds, md2)); err != nil {
		t.Error(err)
		return
	}
}

func MsFwsEqual(tx *Tx, want MsFws) error {
	var got MsFws
	if err := tx.Find(&got); err != nil {
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

func MsFwEqual(want, got MsFw) error {

	if want.Time.Unix() != got.Time.Unix() {
		return fmt.Errorf("\nwant:%+v,\ngot :%+v", want, got)
	}
	got.Time = want.Time

	if !cmp.Equal(want, got) {
		return fmt.Errorf("\nwant:%+v,\ngot :%+v", want, got)
	}
	return nil

}
