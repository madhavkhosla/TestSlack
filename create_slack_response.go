package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type ResponseElement struct {
	Title      string   `json:"title"`
	TitleLink  string   `json:"title_link"`
	Fields     []Field  `json:"fields"`
	Color      string   `json:"color"`
	ThumbUrl   string   `json:"thumb_url"`
	Actions    []Action `json:"actions"`
	CallbackId string   `json:"callback_id"`
}

type SlackResponse struct {
	ResponseType string            `json:"response_type"`
	Text         string            `json:"text"`
	Attachments  []ResponseElement `json:"attachments"`
}

func GetResponseElement(restaurantNames []RestaurantDetails, lastCount int, cuisineId int) ([]byte, error) {
	responseElements := make([]ResponseElement, 0)
	for _, r := range restaurantNames {
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
}
