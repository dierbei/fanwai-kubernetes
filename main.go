package main

import (
	"bytes"
	"context"
	_ "embed"
	"io"
	"log"
	"os/exec"

	"githup.com/dierbei/fanwai-kubernetes/config"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	syaml "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/restmapper"
)

//go:embed tpls/deploy.yaml
var deployTpl string

func main() {
	// 动态客户端
	dynamicClient := config.NewKubernetesConfig().InitDynamicClient()
	// 客户端
	clientSet := config.NewKubernetesConfig().InitClient()

	// kustomize读取
	strData := kustomize("deploy")
	// 解码器
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewBufferString(strData), len(strData))

	for {
		// 获取到每一个对象 yaml中以 --- 分割
		var rawObj runtime.RawExtension
		err := decoder.Decode(&rawObj)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}

		// 获得 groupKindVersion
		obj, gvk, err := syaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).
			Decode(rawObj.Raw, nil, nil)
		if err != nil {
			log.Fatal(err)
		}

		// 创建非结构对象
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			log.Fatal(err)
		}
		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

		// 获取groupResource列表
		gr, err := restmapper.GetAPIGroupResources(clientSet.Discovery())
		if err != nil {
			log.Fatal(err)
		}
		mapper := restmapper.NewDiscoveryRESTMapper(gr)

		// 根据指定group version 获取GVR
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		gvr := schema.GroupVersionResource{
			Group:    mapping.Resource.Group,
			Version:  mapping.Resource.Version,
			Resource: mapping.Resource.Resource,
		}

		// 动态客户端创建
		_, err = dynamicClient.Resource(gvr).Namespace(unstructuredObj.GetNamespace()).
			Create(context.Background(), unstructuredObj, v1.CreateOptions{})
		if err != nil {
			log.Println(err)
		}
	}
}

func kustomize(path string) string {
	var buf = &bytes.Buffer{}
	cmd := exec.Command("D:/Go1.17/bin/kustomize.exe", "build", path)
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		log.Println(err)
	}

	return buf.String()
}
