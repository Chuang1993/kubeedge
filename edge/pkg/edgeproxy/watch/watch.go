package watch

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/klog"

	"github.com/kubeedge/beehive/pkg/common/util"
	"github.com/kubeedge/beehive/pkg/core/model"
	"github.com/kubeedge/kubeedge/common/constants"
)

var (
	brocasters sync.Map
)

func SendMsg(msg model.Message) error {
	_, msgType, _, err := util.ParseResourceEdge(msg.GetResource(), msg.GetOperation())
	if err != nil {
		klog.Error("error")
		return err
	}
	var content []byte
	switch msg.Content.(type) {
	case []byte:
		content = msg.GetContent().([]byte)
	default:
		content, err = json.Marshal(msg.Content)
		if err != nil {
			klog.Errorf("marshal message content failed: %v", err)
			return err
		}
	}

	if msgType == constants.ResourceTypeService {
		return sendSingleObject(msg)
	} else if msgType == constants.ResourceTypeServiceList {
		var svcs []v1.Service
		err := json.Unmarshal(content, &svcs)
		if err != nil {
			return err
		}
		accesser := meta.NewAccessor()
		for i := range svcs {
			svc := svcs[i]
			rv, _ := accesser.ResourceVersion(&svc)

			singmsg := model.NewMessage("").BuildRouter(msg.Router.Source, msg.Router.Group, fmt.Sprintf("%s/%s/%s", svc.Namespace, constants.ResourceTypeService, svc.Name), msg.Router.Operation).SetResourceVersion(rv)
			singmsg.Content = svc
			sendSingleObject(*singmsg)
		}

	}

	return nil
}

func sendSingleObject(msg model.Message) error {
	_, msgType, _, err := util.ParseResourceEdge(msg.GetResource(), msg.GetOperation())
	if err != nil {
		klog.Error("error")
		return err
	}
	var brocaster *MsgBrocaster
	value, ok := brocasters.Load(msgType)
	if !ok {
		brocaster = NewMsgBrocaster()
		brocasters.Store(msgType, brocaster)
	} else {
		brocaster = value.(*MsgBrocaster)
	}
	brocaster.Brocast(msg)
	return nil
}

func AddWatcher(brocastKey string, watcherkey string, resourceVersion uint64, wc IncomingWatcher, timeout time.Duration) {
	var brocaster *MsgBrocaster
	value, ok := brocasters.Load(brocastKey)
	if !ok {
		brocaster = NewMsgBrocaster()
		brocasters.Store(brocastKey, brocaster)
	} else {
		brocaster = value.(*MsgBrocaster)
	}
	brocaster.AddWatcher(watcherkey, wc, timeout, resourceVersion)

}
