package relation

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/kubeedge/beehive/pkg/core"
	"k8s.io/klog"
	"sync"
)

type Manager interface {
	Init()
	UpdateList(groupResource string, list string) error
	//UpdateAll(relations *[]Relation) error
	Update(relation *Relation) error
	GetList(groupResource string) (string, bool)
	GetKind(groupResource string) (string, bool)
}

func InitDBTable(module core.Module) {
	if !module.Enable() {
		klog.Infof("Module %s is disabled, DB cache for it will not be registered", module.Name())
		return
	}
	orm.RegisterModel(new(Relation))
}

type manager struct {
	sync.Once
	topology sync.Map
}

func (m *manager) UpdateList(groupResource string, list string) error {
	v, ok := m.topology.Load(groupResource)
	if !ok {
		return fmt.Errorf("relation does not exit groupresource %s", groupResource)
	}
	r, ok := v.(*Relation)
	if !ok {
		return fmt.Errorf("relation groupresource %s convert to *Relation type error", groupResource)
	}
	r.List = list
	err := m.Update(r)
	return err
}

func (m *manager) Update(relation *Relation) error {
	err := InsertOrUpdate(relation)
	if err != nil {
		return err
	}
	m.topology.Store(relation.GroupResource, relation)
	return nil
}

//func (m *manager) UpdateAll(relations *[]Relation) error {
//	for _, relation := range *relations {
//		err := m.Update(&relation)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

func (m *manager) GetList(groupResource string) (string, bool) {
	v, ok := m.topology.Load(groupResource)
	if !ok {
		return "", ok
	}
	r, ok := v.(*Relation)
	if !ok {
		return "", false
	}
	return r.List, ok
}

func (m *manager) GetKind(groupResource string) (string, bool) {
	v, ok := m.topology.Load(groupResource)
	if !ok {
		return "", ok
	}
	r, ok := v.(*Relation)
	if !ok {
		return "", false
	}
	return r.Kind, ok
}

func (m *manager) Init() {
	m.Do(func() {
		relations, err := QueryAll()
		if err != nil {
			klog.Errorf("relation manager query table error! %v", err)
			return
		}
		for i, relation := range *relations {
			m.topology.Store(relation.GroupResource, &(*relations)[i])
		}
	})
}
