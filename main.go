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
	frontPost := map[string]interface{}{
		"spec": map[string]interface{}{
			"replicas": 2,
		},
	}
	b, err := json.Marshal(frontPost)
	if err != nil {
		log.Println(err)
	}

	// patch更新
	_, err = clientSet.AppsV1().Deployments("default").
		Patch(context.Background(), "my-nginx", types.StrategicMergePatchType, b, metav1.PatchOptions{})
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
