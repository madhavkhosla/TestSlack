package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"strconv"

	"google.golang.org/appengine"
)

type Value struct {
	LastCount string `json:"last_count"`
	CuisineId string `json:"cuisine_id"`
}

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
	v := Value{}
	err = json.Unmarshal([]byte(interactiveRequestMessage.Actions[0].Value), &v)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	ctx := appengine.NewContext(r)
	lastCount, _ := strconv.Atoi(v.LastCount)
	cusId, _ := strconv.Atoi(v.CuisineId)
	restaurantNames, lastCount, err := GetRestaurantNamesInCityByCuisine(ctx, cusId, lastCount)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	} else {
		resp, err := GetResponseElement(restaurantNames, lastCount, cusId)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s\n", resp)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, fmt.Sprintf("Count %s %s", v.CuisineId, v.LastCount))
	}
}

func GetRestaurants(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	cusineName := r.FormValue("text")
	inputCusineId, err := ConvertNameID(cusineName, ctx)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	restaurantNames, lastCount, err := GetRestaurantNamesInCityByCuisine(ctx, inputCusineId.Value, 0)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	} else {
		resp, err := GetResponseElement(restaurantNames, lastCount, inputCusineId.Value)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s\n", resp)
	}
}
