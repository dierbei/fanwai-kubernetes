package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"os/exec"
)

//go:embed tpls/deploy.yaml
var deployTpl string

func main() {
	//dynamicClient := config.NewKubernetesConfig().InitDynamicClient()
	//deployGVR := schema.GroupVersionResource{
	//	Group:    "apps",
	//	Version:  "v1",
	//	Resource: "deployments",
	//}

	fmt.Println(kustomize("deploy"))
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
