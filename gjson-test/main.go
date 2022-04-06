package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"os"
	"regexp"
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

func main() {
	namePattern := "^my"
	b, _ := os.ReadFile("list.json")
	ret := gjson.Get(string(b), "list.#.metadata.name")
	for _, r := range ret.Array() {
		if m, err := regexp.MatchString(namePattern, r.String()); err == nil && m {
			fmt.Println(r.String())
		}
	}
}
