import os
from jinja2 import Template

f = open("service.txt", "r")
contents = f.read()


tpl = Template(contents)

t1 = '''
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
'''
t2 = '''
type MsfwShareCom struct {
	ID  uint64 `json:"id"`
	Com uint64 `json:"com"`
}
'''
t3 = '''
type MsfwShareGroup struct {
	ID    uint64 `json:"id"`
	Group uint64 `json:"gp"`
}
'''
t4 = '''
type GroupMsfw struct {
	Gp uint64 `json:"gp"`
	ID uint64 `json:"id"`
}
'''

n1 = '''
func NewMsfw(version string, com uint64, bk, obj string, time time.Time, tag string) Msfw {
	msfwkey := MsfwUnique{Version: version, Com: com}
	return Msfw{MsfwUnique: msfwkey, Bucket: bk, Obj: obj, Time: time, Tag: tag}
}

func NewMsfwUnique(com uint64, version string) MsfwUnique {
	return MsfwUnique{Com: com, Version: version}
}
'''

n2 = '''
func NewMsfwShareCom(id, com uint64) MsfwShareCom {
	return MsfwShareCom{ID: id, Com: com}
}
'''

n3 = '''
func NewMsfwShareGroup(id, gp uint64) MsfwShareGroup {
	return MsfwShareGroup{ID: id, Group: gp}
}
'''

n4 = '''
func NewGroupMsfw(gp, id uint64) GroupMsfw {
	return GroupMsfw{Gp: gp, ID: id}
}
'''


md1 = {'def': t1, 'serial': True, 'Model': 'Msfw',
       'Models': 'Msfws', 'model': 'md', 'new': n1}
md2 = {'def': t2, 'serial': False, 'Model': 'MsfwShareCom',
       'Models': 'MsfwShareComs', 'model': 'md', 'new': n2}
md3 = {'def': t3, 'serial': False, 'Model': 'MsfwShareGroup',
       'Models': 'MsfwShareGroups', 'model': 'md', 'new': n3}
md4 = {'def': t4, 'serial': False, 'Model': 'GroupMsfw',
       'Models': 'GroupMsfws', 'model': 'md', 'new': n4}

models = [md1, md2, md3, md4]
d = {
    'package': 'msfw',
    'models': models,
    'import': '',
    'method': ''
}

di = os.path.join('../out/service/msfw')
if not os.path.exists(di):
    os.makedirs(di)


tpl.stream(d).dump(os.path.join(di, 'service.go'))
