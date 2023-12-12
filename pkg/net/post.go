package net

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"library/pkg/io"
)

type post_config struct {
	ProcRequest func(map[string]interface{}, interface{}) (interface{}, error)
}

func PostHandler(wri http.ResponseWriter, req *http.Request, args ...interface{}) {
	conf := args[0].(*post_config)
	work := args[2].(chan interface{})

	//switch based on method passed
	switch req.Method {
	case "POST":
		//read request
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Printf("Error reading body: %+v\n", err)
			fmt.Printf("\tbody: %s\n", req.Body)
		}
		data := map[string]interface{}{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Printf("Error unmarshaling data: %+v\n", err)
			fmt.Printf("\tdata: %s\n", body)
			fmt.Fprintf(wri, "%+v", err)
			return
		}

		//process request
		hash := fmt.Sprintf("%s", sha256.Sum256([]byte(fmt.Sprintf("%s", data))))
		c := make(chan interface{})
		work <- map[string]interface{}{hash: []interface{}{data, c, conf.ProcRequest}}
		res := <-c

		//write result
		bytes, err := io.UnescapeJson(res)
		if err != nil {
			fmt.Printf("Error marshaling data: %+v\n", err)
			fmt.Printf("\tdata: %s\n", data)
			fmt.Printf("\tresult: %s\n", res)
			fmt.Fprintf(wri, "%+v", err)
			return
		}

		fmt.Fprintf(wri, "%s", bytes)
	default:
		fmt.Fprintf(wri, "Only POST requests are supported on this endpoint\n")
	}
}
