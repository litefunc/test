package msfw

import "time"

type Service struct {
	DB DB
}

type DB interface {
	Begin() (Tx, error)
	MsfwStore
	MsfwShareComStore
	MsfwShareGroupStore
	GroupMsfwStore
}

type Tx interface {
	Commit() error
	Rollback() error
	MsfwStore
	MsfwShareComStore
	MsfwShareGroupStore
	GroupMsfwStore
}

type Msfw struct {
	ID uint64 `json:"id"`
	MsfwUnique
	Bucket string    `json:"bucket"`
	Obj    string    `json:"obj"`
	Time   time.Time `json:"time"`
	Tag    string    `json:"tag"`
}

type MsfwUnique struct {
	Com     uint64 `json:"com"`
	Version string `json:"ver"`
}

type Msfws []Msfw

func NewMsfw(version string, com uint64, bk, obj string, time time.Time, tag string) Msfw {
	msfwkey := MsfwUnique{Version: version, Com: com}
	return Msfw{MsfwUnique: msfwkey, Bucket: bk, Obj: obj, Time: time, Tag: tag}
}

func NewMsfwUnique(com uint64, version string) MsfwUnique {
	return MsfwUnique{Com: com, Version: version}
}

type MsfwStore interface {
	InsertMsfw(Msfw) (uint64, error)
	GetMsfws(condition map[string]interface{}, cols ...string) (Msfws, error)
	GetOneMsfw(condition map[string]interface{}, cols ...string) (Msfw, error)
	UpdateMsfw(set, condition map[string]interface{}) error
	DeleteMsfw(condition map[string]interface{}) error
	CountMsfw(condition map[string]interface{}) (uint64, error)
}

type MsfwShareCom struct {
	ID  uint64 `json:"id"`
	Com uint64 `json:"com"`
}

type MsfwShareComs []MsfwShareCom

func NewMsfwShareCom(id, com uint64) MsfwShareCom {
	return MsfwShareCom{ID: id, Com: com}
}

type MsfwShareComStore interface {
	InsertMsfwShareCom(MsfwShareCom) error
	GetMsfwShareComs(condition map[string]interface{}, cols ...string) (MsfwShareComs, error)
	GetOneMsfwShareCom(condition map[string]interface{}, cols ...string) (MsfwShareCom, error)
	UpdateMsfwShareCom(set, condition map[string]interface{}) error
	DeleteMsfwShareCom(condition map[string]interface{}) error
	CountMsfwShareCom(condition map[string]interface{}) (uint64, error)
}

type MsfwShareGroup struct {
	ID    uint64 `json:"id"`
	Group uint64 `json:"gp"`
}

type MsfwShareGroups []MsfwShareGroup

func NewMsfwShareGroup(id, gp uint64) MsfwShareGroup {
	return MsfwShareGroup{ID: id, Group: gp}
}

type MsfwShareGroupStore interface {
	InsertMsfwShareGroup(MsfwShareGroup) error
	GetMsfwShareGroups(condition map[string]interface{}, cols ...string) (MsfwShareGroups, error)
	GetOneMsfwShareGroup(condition map[string]interface{}, cols ...string) (MsfwShareGroup, error)
	UpdateMsfwShareGroup(set, condition map[string]interface{}) error
	DeleteMsfwShareGroup(condition map[string]interface{}) error
	CountMsfwShareGroup(condition map[string]interface{}) (uint64, error)
}

type GroupMsfw struct {
	Gp uint64 `json:"gp"`
	ID uint64 `json:"id"`
}

type GroupMsfws []GroupMsfw

func NewGroupMsfw(gp, id uint64) GroupMsfw {
	return GroupMsfw{Gp: gp, ID: id}
}

type GroupMsfwStore interface {
	InsertGroupMsfw(GroupMsfw) error
	GetGroupMsfws(condition map[string]interface{}, cols ...string) (GroupMsfws, error)
	GetOneGroupMsfw(condition map[string]interface{}, cols ...string) (GroupMsfw, error)
	UpdateGroupMsfw(set, condition map[string]interface{}) error
	DeleteGroupMsfw(condition map[string]interface{}) error
	CountGroupMsfw(condition map[string]interface{}) (uint64, error)
}
