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
	http.HandleFunc("/five", GetFive)
}

func GetFive(w http.ResponseWriter, r *http.Request) {
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
	if interactiveRequestMessage.Actions[0].Name == "next" {
		GetNextFive(w, r, &v)
	} else {
		GetPrevFive(w, r, &v)
	}

}

func GetPrevFive(w http.ResponseWriter, r *http.Request, v *RestaurantStat) {
	ctx := appengine.NewContext(r)
	lastCount := v.LastCount
	cusId := v.CuisineId

	start := lastCount - 10
	lastCount = start + 5
	restaurantNames, _, err := GetRestaurantNamesInCityByCuisine(ctx, cusId, start)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	} else {
		resultRemaining := v.CountRemaining + 5
		restaurantStat := &RestaurantStat{LastCount: lastCount, CuisineId: cusId, CountRemaining: resultRemaining}
		resp, err := GetResponseElement(restaurantNames, restaurantStat)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s\n", resp)
	}
}

func GetNextFive(w http.ResponseWriter, r *http.Request, v *RestaurantStat) {
	ctx := appengine.NewContext(r)
	lastCount := v.LastCount
	cusId := v.CuisineId

	start := lastCount
	lastCount = start + 5
	restaurantNames, _, err := GetRestaurantNamesInCityByCuisine(ctx, cusId, start)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	} else {
		resultRemaining := v.CountRemaining - 5
		restaurantStat := &RestaurantStat{LastCount: lastCount, CuisineId: cusId, CountRemaining: resultRemaining}
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
	start := 0
	restaurantNames, resultsFound, err := GetRestaurantNamesInCityByCuisine(ctx, inputCusineId.Value, start)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	} else {
		lastCount := start + 5
		resultRemaining := resultsFound - lastCount
		restaurantStat := &RestaurantStat{LastCount: lastCount, CuisineId: inputCusineId.Value, CountRemaining: resultRemaining}
		resp, err := GetResponseElement(restaurantNames, restaurantStat)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s\n", resp)
	}
}
