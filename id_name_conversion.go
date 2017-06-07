package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"errors"

	"golang.org/x/net/context"

	"google.golang.org/appengine/urlfetch"
)

type CusinesStruct struct {
	Cuisine struct {
		CuisineId   int    `json:"cuisine_id"`
		CuisineName string `json:"cuisine_name"`
	} `json:"cuisine"`
}

type Collection struct {
	Cuisines []CusinesStruct `json:"cuisines"`
}

type Integer struct {
	Value int
}

func ConvertNameID(cusineName string, ctx context.Context) (*Integer, error) {
	baseUrl := "https://developers.zomato.com/api/v2.1/"
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/cuisines", baseUrl), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Key", "a881a2c3cbd4e8320634917542051763")
	query := req.URL.Query()
	query.Add("city_id", "280")
	req.URL.RawQuery = query.Encode()
	client := urlfetch.Client(ctx)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	collection := Collection{}
	err = json.NewDecoder(resp.Body).Decode(&collection)
	if err != nil {
		return nil, err
	}
	var inputCuisineId int = -1
	for _, c := range collection.Cuisines {
		if c.Cuisine.CuisineName == cusineName {
			inputCuisineId = c.Cuisine.CuisineId
		}
	}
	if inputCuisineId > 0 {
		return &Integer{Value: inputCuisineId}, nil
	} else {
		return nil, errors.New("Unknown cuisine name")
	}
}
