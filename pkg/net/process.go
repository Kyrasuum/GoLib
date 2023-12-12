package net

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type webhandler func(params map[string]interface{}, user interface{}, level int, session interface{}, conf interface{}) (interface{}, error)

type proc_config struct {
	DefaultParams map[string]interface{}
	WebHandlers   map[string]webhandler
	GetUser       func(user_id string, conf interface{}) (interface{}, error)
	GetUserLevel  func(user interface{}) int
	GetSession    func(user_id string, conf interface{}) (interface{}, error)
}

func ProcParams(request map[string]interface{}, conf proc_config) (params map[string]interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error on validating parameters: \n%+v\n%+v\n\n", request, r.(error))
			err = fmt.Errorf("Invalid Parameter")
		}
	}()

	//defaults
	params = map[string]interface{}{}
	for param, val := range conf.DefaultParams {
		params[param] = val
	}

	//validate parameters
	for param, val := range request {
		if _, ok := params[param]; ok {
			if reflect.TypeOf(val) == reflect.TypeOf(params[param]) {
				params[param] = val
			} else {
				return nil, fmt.Errorf("Invalid Parameter:\nGot %T on %s, expected %T", val, param, params[param])
			}
		} else {
			return nil, fmt.Errorf("Invalid Parameter: %s", param)
		}
	}

	return params, nil
}

func ProcRequest(request map[string]interface{}, conf proc_config) (response interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	//read parameters
	params, err := ProcParams(request, conf)
	if err != nil {
		return nil, err
	}

	//get user information
	curr_user, err := conf.GetUser(params["user_id"].(string), conf)
	if err != nil {
		return "", err
	}

	//get session
	curr_sess, err := conf.GetSession(params["sess_id"].(string), conf)
	if err != nil {
		return "", err
	}

	//get level
	level := 0
	level = conf.GetUserLevel(curr_user)

	//switch based on action
	if handler, ok := conf.WebHandlers[request["action"].(string)]; ok {
		response, err = handler(params, curr_user, level, curr_sess, conf)
		if err != nil {
			return nil, fmt.Errorf("Error: %+v", err)
		}
	} else {
		return nil, fmt.Errorf("Invalid Endpoint")
	}
	//return response
	switch response.(type) {
	case []byte:
		msg := []map[string]interface{}{}
		err = json.Unmarshal(response.([]byte), &msg)
		if err != nil {
			return response, nil
		}
		return msg, nil
	default:
		return response, nil
	}
}
