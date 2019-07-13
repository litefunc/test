package model

type Model struct {
	Table    string
	Model    string
	Models   string
	GetOneBy []By
	GetBy    []By
	GetByIn  []By
}

func New(tb, md, mds string, get, gets, getsin []By) Model {
	return Model{
		Table:    tb,
		Model:    md,
		Models:   mds,
		GetOneBy: get,
		GetBy:    gets,
		GetByIn:  getsin,
	}
}

type By struct {
	Name   string
	Params string
	Query  string
	Args   string
}

func NewBy(name, params, query, args string) By {
	return By{
		Name:   name,
		Params: params,
		Query:  query,
		Args:   args,
	}
}

var (
	by1      = NewBy("Id", "id uint64", "id = ?", "id")
	by2      = NewBy("Unique", "u MsFwUnique", "com = ? AND version = ?", "u.Com, u.Version")
	by3      = NewBy("Com", "com uint64", "com = ?", "com")
	by4      = NewBy("Coms", "coms []uint64", "com = ?", "coms")
	getOneBy = []By{by1, by2, by3}
	getByIn  = []By{by4}
	Md       = New("cloud.msfw", "MsFw", "MsFws", getOneBy, nil, getByIn)
)
