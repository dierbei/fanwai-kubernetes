package main

import (
	"bytes"
	"context"
	_ "embed"
	"log"
	"os/exec"

	"githup.com/dierbei/fanwai-kubernetes/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
)

//go:embed tpls/deploy.yaml
var deployTpl string

type JSONPatch struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

type JSONPatchList []*JSONPatch

func AddJsonPatch(jps ...*JSONPatch) JSONPatchList {
	list := make([]*JSONPatch, len(jps))
	for index, jp := range jps {
		list[index] = jp
	}
	return list
}

func main() {
	// 动态客户端
	//dynamicClient := config.NewKubernetesConfig().InitDynamicClient()
	// 客户端
	clientSet := config.NewKubernetesConfig().InitClient()

	// 是否存在对应deployment
	_, err := clientSet.AppsV1().Deployments("default").Get(context.Background(), "my-nginx", metav1.GetOptions{})
	if err != nil {
		log.Println(err)
	}

	// 模拟序列化前端数据
	list := AddJsonPatch(&JSONPatch{
		Op:   "add",
		Path: "/spec/template/spec/containers/1",
		Value: map[string]interface{}{
			"name":  "my-nginx",
			"image": "nginx:1.18-alpine",
		},
	})
	b, err := json.Marshal(list)
	if err != nil {
		log.Println(err)
	}

	// patch更新
	_, err = clientSet.AppsV1().Deployments("default").
		Patch(context.Background(), "my-nginx", types.JSONPatchType, b, metav1.PatchOptions{})
	if err != nil {
		log.Println(err)
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
