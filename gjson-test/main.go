package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"log"
)

var data = `
{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
  ]
}`

var (
	cfgFlags = &genericclioptions.ConfigFlags{}
	client   = InitClient()
)

// InitClient 初始化Kubernetes客户端
func InitClient() *kubernetes.Clientset {
	cfgFlags = genericclioptions.NewConfigFlags(true)
	config, err := cfgFlags.ToRawKubeConfigLoader().ClientConfig()
	if err != nil {
		log.Fatal(err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func main() {
	InitSharedInformerFactory()
	//namePattern := "^my"
	//b, _ := os.ReadFile("list.json")
	//ret := gjson.Get(string(b), "list.#.metadata.name")
	//for _, r := range ret.Array() {
	//	if m, err := regexp.MatchString(namePattern, r.String()); err == nil && m {
	//		fmt.Println(r.String())
	//	}
	//}

	pod, err := fact.Core().V1().Pods().Lister().Pods("default").Get("my-nginx-6fdfc68494-llwsm")
	if err != nil {
		log.Println(err)
	}

	jsonStr, err := json.Marshal(pod)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(jsonStr), "---------------------")
	//v1.Pod{}
	resultData := gjson.Get(string(jsonStr), "metadata")
	if !resultData.Exists() {
		log.Println("无法找到对应内容, ")
	}
	fmt.Println(resultData.Raw)
}

var fact informers.SharedInformerFactory

func InitSharedInformerFactory() {
	fact = informers.NewSharedInformerFactoryWithOptions(client, 0)
	fact.Core().V1().Pods().Informer().AddEventHandler(&PodHandler{})
	ch := make(chan struct{})
	fact.Start(ch)
	fact.WaitForCacheSync(ch)
}

type PodHandler struct {
}

func (h *PodHandler) OnAdd(obj interface{}) {

}

func (h *PodHandler) OnUpdate(oldObj, newObj interface{}) {

}

func (h *PodHandler) OnDelete(obj interface{}) {

}
