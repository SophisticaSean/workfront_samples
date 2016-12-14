package main

import (
	"encoding/json"
	"fmt"
	"os"

	URL "net/url"

	"github.com/davecgh/go-spew/spew"
	"github.com/levigross/grequests"
)

type issue struct {
	ID      string `json:"ID"`
	Name    string `json:"name"`
	ObjCode string `json:"objCode"`
}

type issues struct {
	Issues []issue `json:"data"`
}

func main() {

	if len(os.Args) != 4 {
		fmt.Println(`yo dawg, you gotta give me an API Key, and instance, and a projectID (in that order) like this:
			./all_issues_to_tasks some_api_key my_instance_name some_project_id
		`)
	} else {
		apiKey := os.Args[1]
		instance := os.Args[2]
		projectID := os.Args[3]

		url := "https://" + instance + ".attask-ondemand.com/attask/api/optask/search?fields=ID&apiKey=" + apiKey + "&projectID=" + projectID

		ro := &grequests.RequestOptions{}
		resp, err := grequests.Get(url, ro)
		if err != nil {
			panic(err)
		}

		var returnList issues

		if err := json.Unmarshal(resp.Bytes(), &returnList); err != nil {
			panic(err)
		}

		spew.Dump(returnList)

		for _, issue := range returnList.Issues {
			ro := &grequests.RequestOptions{}

			url := `https://` + instance + `.attask-ondemand.com/attask/api-internal/optask/` + issue.ID + `/convertToTask?apiKey=` + apiKey + `&updates={"options":["preservePrimaryContact"],"task":{"name":"` + URL.QueryEscape(issue.Name) + `"}}&method=PUT`
			fmt.Println(url)
			resp, err := grequests.Put(url, ro)
			if err != nil {
				panic(err)
			}

			spew.Dump(resp.Bytes())
			spew.Dump(resp.StatusCode)
		}
	}
}
