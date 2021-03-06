package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

type RestaurantDetails struct {
	Fields   []Field
	MenuUrl  string
	Name     string
	ThumbUrl string
}

type RestaurantStat struct {
	LastCount      int `json:"last_count"`
	CuisineId      int `json:"cuisine_id"`
	CountRemaining int `json:"count_remaining"`
}

type Restaurant struct {
	Name     string `json:"name"`
	Location struct {
		Address  string `json:"address"`
		Locality string `json:"locality"`
	} `json:"location"`
	MenuUrl           string `json:"menu_url"`
	AverageCostForTwo int    `json:"average_cost_for_two"`
	UserRating        struct {
		AggregateRating string `json:"aggregate_rating"`
	} `json:"user_rating"`
	Thumb string `json:"thumb"`
}

type SearchResult struct {
	ResultsFound int    `json:"results_found"`
	ResultsStart string `json:"results_start"`
	ResultsShown int    `json:"results_shown"`
	Restaurants  []struct {
		Restaurant Restaurant `json:"restaurant"`
	} `json:"restaurants"`
}

func GetRestaurantNamesInCityByCuisine(ctx context.Context,
	inputCuisineId int, start int) ([]RestaurantDetails, int, error) {
	req, err := http.NewRequest("GET", "https://developers.zomato.com/api/v2.1/search", nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Key", "a881a2c3cbd4e8320634917542051763")
	query := req.URL.Query()
	query.Add("entity_id", "280")
	query.Add("entity_type", "city")
	query.Add("cuisines", strconv.Itoa(inputCuisineId))
	query.Add("lat", "40.737920")
	query.Add("lon", "-73.992781")
	query.Add("radius", "1000")
	query.Add("start", strconv.Itoa(start))
	query.Add("count", "5")
	query.Add("sort", "cost")
	query.Add("order", "asc")
	req.URL.RawQuery = query.Encode()

	client := urlfetch.Client(ctx)
	//client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, 0, err
	}
	searchRes := SearchResult{}
	err = json.NewDecoder(resp.Body).Decode(&searchRes)
	if err != nil {
		return nil, 0, err
	}

	restaurantNameSlice := make([]RestaurantDetails, 0)
	for _, restaurant := range searchRes.Restaurants {
		fields := make([]Field, 0)
		fields = append(fields, Field{
			Title: "Location", Value: restaurant.Restaurant.Location.Address, Short: true})
		fields = append(fields, Field{
			Title: "Average Cost for Two", Value: strconv.Itoa(restaurant.Restaurant.AverageCostForTwo), Short: true})
		fields = append(fields, Field{
			Title: "Rating", Value: restaurant.Restaurant.UserRating.AggregateRating, Short: true})

		restaurantNameSlice = append(restaurantNameSlice, RestaurantDetails{Fields: fields,
			MenuUrl: restaurant.Restaurant.MenuUrl, Name: restaurant.Restaurant.Name,
			ThumbUrl: restaurant.Restaurant.Thumb})
	}
	return restaurantNameSlice, searchRes.ResultsFound, nil
}
