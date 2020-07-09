package watch

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/kubeedge/beehive/pkg/core/model"
)

var (
	msgOperationToWatchEvent = map[string]watch.EventType{
		model.DeleteOperation:        watch.Deleted,
		model.InsertOperation:        watch.Added,
		model.UpdateOperation:        watch.Modified,
		model.ResponseErrorOperation: watch.Error,
	}
)

func ConvertMsgToEvent(msg model.Message) watch.Event {
	action := msgOperationToWatchEvent[msg.Router.Operation]
	obj := &unstructured.Unstructured{}
	var content []byte
	switch msg.Content.(type) {
	case []byte:
		content = msg.GetContent().([]byte)
	default:
		content, _ = json.Marshal(msg.Content)

	}
	json.Unmarshal(content, obj)
	event := watch.Event{
		Type:   action,
		Object: obj,
	}
	return event
}
