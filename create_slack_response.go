package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func GetResponseElement(restaurantNames []RestaurantDetails, lastCount int, cuisineId int) ([]byte, error) {
	responseElements := make([]ResponseElement, 0)
	for _, r := range restaurantNames {

		//responseTextFormatted := fmt.Sprintf("Name: %s\n Address: %s\n Menu url: %s\n Cost for 2: %s\n, Rating: %s\n",
		//	r.Name, r.Address, r.MenuUrl, r.AverageCostForTwo, r.AggregateRating)
		responseElements = append(responseElements, ResponseElement{
			Title:     fmt.Sprintf("%s", r.Name),
			TitleLink: r.MenuUrl,
			Fields:    r.Fields,
			Color:     "#36a64f",
			ThumbUrl:  r.ThumbUrl,
		})
	}
	val, _ := json.Marshal(Value{LastCount: strconv.Itoa(lastCount), CuisineId: strconv.Itoa(cuisineId)})
	action := Action{
		Name:  "Get",
		Text:  "Get",
		Type:  "button",
		Value: string(val)}
	buttonAttachment := ResponseElement{Title: "Get next 5 restaurants", Actions: []Action{
		action,
	}, CallbackId: "zomato_next5", Color: "#36a64f"}
	responseElements = append(responseElements, buttonAttachment)
	response := SlackResponse{
		ResponseType: "ephemeral",
		Text:         fmt.Sprintf("The Restaurants for cuisine are "),
		Attachments:  responseElements,
	}
	resp, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return resp, nil
	//w.Header().Set("Content-Type", "application/json")
	//fmt.Fprintf(w, "%s\n", resp)
}
