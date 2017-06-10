package main

import (
	"encoding/json"
	"fmt"
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

func GetActionButton(name string, val []byte) Action {
	return Action{
		Name:  name,
		Text:  name,
		Type:  "button",
		Value: string(val)}
}

func GetResponseElement(restaurantNames []RestaurantDetails, restaurantStat *RestaurantStat) ([]byte, error) {
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
	valByteArray, err := json.Marshal(restaurantStat)
	if err != nil {
		return nil, err
	}
	var actionList []Action
	if restaurantStat.CountRemaining <= 0 {
		actionPrev := GetActionButton("prev", valByteArray)
		actionList = append(actionList, actionPrev)
	} else if restaurantStat.LastCount-5 <= 0 {
		actionNext := GetActionButton("next", valByteArray)
		actionList = append(actionList, actionNext)
	} else {
		actionNext := GetActionButton("next", valByteArray)
		actionPrev := GetActionButton("prev", valByteArray)
		actionList = append(actionList, actionNext, actionPrev)
	}

	buttonAttachment := ResponseElement{Title: "Get page",
		Actions: actionList, CallbackId: "zomato_5", Color: "#36a64f"}
	responseElements = append(responseElements, buttonAttachment)

	return createSlackResponse(responseElements)
}

func createSlackResponse(responseElements []ResponseElement) ([]byte, error) {
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
