package relation

import (
	"github.com/kubeedge/kubeedge/edge/pkg/common/dbm"
	"k8s.io/klog"
)

const (
	RelationTableName = "relation"
)

// Relation matadata object
type Relation struct {
	GroupResource string `orm:"column(groupresource); size(256); pk"`
	Kind          string `orm:"column(kind); size(256)"`
	List          string `orm:"column(list); size(256)"`
}

func InsertOrUpdate(relation *Relation) error {
	_, err := dbm.DBAccess.Raw("INSERT OR REPLACE INTO relation (groupresource, kind, list) VALUES (?,?,?)", relation.GroupResource, relation.Kind, relation.List).Exec() // will update all field
	klog.V(4).Infof("Update result %v", err)
	return err
}

func QueryAll() (*[]Relation, error) {
	relations := new([]Relation)
	_, err := dbm.DBAccess.QueryTable(RelationTableName).All(relations)
	if err != nil {
		return nil, err
	}
	return relations, nil
}
