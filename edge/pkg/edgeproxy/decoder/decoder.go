package decoder

import (
	"errors"
	"io"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/runtime/serializer/streaming"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	restclientwatch "k8s.io/client-go/rest/watch"
)

// Manager interface provides methods to get the corresponding Decoder based on the resource type.
type Manager interface {
	// generate decoder based on contentType and groupVersion
	GetDecoder(contentType string, gv schema.GroupVersion) (runtime.Decoder, error)
	// generate decoder based on contentType and groupVersion
	GetUnstructuredDecoder(contentType string, gv schema.GroupVersion) (runtime.Decoder, error)

	GetUnstructuredListDecoder(contentType string, gv schema.GroupVersion) (runtime.Decoder, error)
	// generate watch decoder based on contentType, groupVersion and readCloser
	GetStreamDecoder(contentType string, gv schema.GroupVersion, rc io.ReadCloser) (watch.Decoder, error)

	GetUnstructuredStreamDecoder(contentType string, gv schema.GroupVersion, rc io.ReadCloser) (watch.Decoder, error)
}

var DefaultDecoderMgr = &mgr{
	codecFactory: serializer.NewCodecFactory(scheme.Scheme),
}

type mgr struct {
	codecFactory serializer.CodecFactory
}

func (dm *mgr) GetDecoder(contentType string, gv schema.GroupVersion) (runtime.Decoder, error) {
	decoder, _, err := dm.getDecoder(contentType, gv)
	return decoder, err
}
func (dm *mgr) GetUnstructuredListDecoder(contentType string, gv schema.GroupVersion) (runtime.Decoder, error) {
	decoder, _, err := dm.getDecoder(contentType, gv)
	if err != nil {
		return decoder, err
	}
	unstructuredListDecoder := &unstructuredListDecoder{
		decoder: decoder,
	}
	return unstructuredListDecoder, err
}

func (dm *mgr) GetUnstructuredDecoder(contentType string, gv schema.GroupVersion) (runtime.Decoder, error) {
	decoder, _, err := dm.getDecoder(contentType, gv)
	if err != nil {
		return decoder, err
	}
	unstructuredObjDecoder := &unstructuredDecoder{
		decoder: decoder,
	}
	return unstructuredObjDecoder, err
}

func (dm *mgr) getDecoder(contentType string, gv schema.GroupVersion) (runtime.Decoder, runtime.SerializerInfo, error) {
	mediaTypes := dm.codecFactory.SupportedMediaTypes()
	info, ok := runtime.SerializerInfoForMediaType(mediaTypes, contentType)
	if !ok {
		if len(contentType) != 0 || len(mediaTypes) == 0 {
			return nil, info, errors.New("content type and midiaTypes'dm length are empty")
		}
		info = mediaTypes[0]
	}
	decoder := dm.codecFactory.DecoderToVersion(info.Serializer, gv)
	return decoder, info, nil
}

func (dm *mgr) GetStreamDecoder(contentType string, gv schema.GroupVersion, rc io.ReadCloser) (watch.Decoder, error) {
	objDecoder, info, err := dm.getDecoder(contentType, gv)
	if err != nil {
		return nil, err
	}
	frameReader := info.StreamSerializer.Framer.NewFrameReader(rc)
	watchEventDecoder := streaming.NewDecoder(frameReader, info.StreamSerializer)
	watchDecoder := restclientwatch.NewDecoder(watchEventDecoder, objDecoder)
	return watchDecoder, nil
}

func (dm *mgr) GetUnstructuredStreamDecoder(contentType string, gv schema.GroupVersion, rc io.ReadCloser) (watch.Decoder, error) {
	objDecoder, info, err := dm.getDecoder(contentType, gv)
	if err != nil {
		return nil, err
	}
	unstructuredObjDecoder := &unstructuredDecoder{
		decoder: objDecoder,
	}
	frameReader := info.StreamSerializer.Framer.NewFrameReader(rc)
	watchEventDecoder := streaming.NewDecoder(frameReader, info.StreamSerializer)
	watchDecoder := restclientwatch.NewDecoder(watchEventDecoder, unstructuredObjDecoder)
	return watchDecoder, nil
}

type unstructuredDecoder struct {
	decoder runtime.Decoder
}

func (u *unstructuredDecoder) Decode(data []byte, defaults *schema.GroupVersionKind, into runtime.Object) (runtime.Object, *schema.GroupVersionKind, error) {
	obj := &unstructured.Unstructured{}
	return u.decoder.Decode(data, defaults, obj)
}

type unstructuredListDecoder struct {
	decoder runtime.Decoder
}

func (u *unstructuredListDecoder) Decode(data []byte, defaults *schema.GroupVersionKind, into runtime.Object) (runtime.Object, *schema.GroupVersionKind, error) {
	obj := &unstructured.UnstructuredList{}
	return u.decoder.Decode(data, defaults, obj)
}
