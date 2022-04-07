package lib

import "fmt"

func ShowLabels(labels map[string]string) string {
	var ret string
	for k, v := range labels {
		ret += fmt.Sprintf("%s=%s\n", k, v)
	}
	return ret
}
