package relation

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"strings"
)

var mgr Manager

func Init() {
	mgr = &manager{}
	mgr.Init()
	// TODO 通过edgehub与cloudhub通信，由upstreamcontroller完成查询
	go func() {
		config, _ := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
		client, _ := kubernetes.NewForConfig(config)
		_, resourcelists, _ := client.DiscoveryClient.ServerGroupsAndResources()
		for _, resourcelist := range resourcelists {
			gv, err := schema.ParseGroupVersion(resourcelist.GroupVersion)
			if err != nil {
				continue
			}
			for _, resource := range resourcelist.APIResources {
				if strings.Contains(resource.Name, "/") {
					continue
				}
				gr := schema.GroupResource{
					Group:    gv.Group,
					Resource: resource.Name,
				}.String()
				list, _ := GetList(gr)
				Update(&Relation{
					GroupResource: gr,
					Kind:          resource.Kind,
					List:          list,
				})
			}
		}
	}()
}

func UpdateList(groupResource string, list string) error {
	return mgr.UpdateList(groupResource, list)
}

func Update(relation *Relation) error {
	return mgr.Update(relation)
}
func GetList(groupResource string) (string, bool) {
	return mgr.GetList(groupResource)
}
func GetKind(groupResource string) (string, bool) {
	return mgr.GetKind(groupResource)
}
