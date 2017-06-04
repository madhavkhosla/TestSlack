package main

import (
	"encoding/json"
	"fmt"
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

func GetRestaurantNamesInCityByCuisine(ctx context.Context, inputCuisineId int, start int) ([]RestaurantDetails, int, error) {
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
	fmt.Println(req.URL.String())
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
	restaurantResultStart, err := strconv.Atoi(searchRes.ResultsStart)
	lastCount := restaurantResultStart + 5
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
	return restaurantNameSlice, lastCount, nil
}
