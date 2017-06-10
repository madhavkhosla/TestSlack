package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/appengine"
)

type InteractiveMessageRequest struct {
	Actions []Action
}

type Action struct {
	Name  string `json:"name"`
	Text  string `json:"text"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

func init() {
	http.HandleFunc("/", GetRestaurants)
	http.HandleFunc("/nextfive", GetNextFive)
}

func GetNextFive(w http.ResponseWriter, r *http.Request) {
	payload := r.FormValue("payload")
	interactiveRequestMessage := InteractiveMessageRequest{}
	err := json.Unmarshal([]byte(payload), &interactiveRequestMessage)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	v := RestaurantStat{}
	err = json.Unmarshal([]byte(interactiveRequestMessage.Actions[0].Value), &v)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	ctx := appengine.NewContext(r)
	lastCount := v.LastCount
	cusId := v.CuisineId

	restaurantNames, restaurantStat, err := GetRestaurantNamesInCityByCuisine(ctx, cusId, lastCount)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	} else {
		resp, err := GetResponseElement(restaurantNames, restaurantStat)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s\n", resp)
	}
}

func GetRestaurants(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	cusineName := r.FormValue("text")
	inputCusineId, err := ConvertNameID(cusineName, ctx)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	restaurantNames, restaurantStat, err := GetRestaurantNamesInCityByCuisine(ctx, inputCusineId.Value, 0)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	} else {
		resp, err := GetResponseElement(restaurantNames, restaurantStat)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s\n", resp)
	}
}
