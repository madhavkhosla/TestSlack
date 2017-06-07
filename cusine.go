package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"strconv"

	"github.com/madhavkhosla/TestSlack/transactions"
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

type CusinesStruct struct {
	Cuisine struct {
		CuisineId   int    `json:"cuisine_id"`
		CuisineName string `json:"cuisine_name"`
	} `json:"cuisine"`
}

type Collection struct {
	Cuisines []CusinesStruct `json:"cuisines"`
}

func init() {
	http.HandleFunc("/", GetRestaurants)
	http.HandleFunc("/nextfive", GetNextFive)
	//http.ListenAndServe(":8081", nil)
}

func GetNextFive(w http.ResponseWriter, r *http.Request) {
	payload := r.FormValue("payload")
	interactiveRequestMessage := InteractiveMessageRequest{}
	err := json.Unmarshal([]byte(payload), &interactiveRequestMessage)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	v := Value{}
	err = json.Unmarshal([]byte(interactiveRequestMessage.Actions[0].Value), &v)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	//newStart, _ := strconv.Atoi(interactiveRequestMessage.Actions[0].Value)

	ctx := appengine.NewContext(r)
	lastCount, _ := strconv.Atoi(v.LastCount)
	cusId, _ := strconv.Atoi(v.CuisineId)
	restaurantNames, lastCount, err := GetRestaurantNamesInCityByCuisine(ctx, cusId, lastCount)
	//restaurantNames, err := GetRestaurantNamesInCityByCuisine(inputCusineId)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	} else {

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
		val, _ := json.Marshal(Value{LastCount: strconv.Itoa(lastCount), CuisineId: v.CuisineId})
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
	inputCusineId, err := transactions.ConvertNameID(cusineName, ctx)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	//if inputCusineId != -1 {
	restaurantNames, lastCount, err := GetRestaurantNamesInCityByCuisine(ctx, inputCusineId.Value, 0)
	//restaurantNames, err := GetRestaurantNamesInCityByCuisine(inputCusineId)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	} else {

		//responseElements := make([]ResponseElement, 0)
		//for _, r := range restaurantNames {
		//
		//	//responseTextFormatted := fmt.Sprintf("Name: %s\n Address: %s\n Menu url: %s\n Cost for 2: %s\n, Rating: %s\n",
		//	//	r.Name, r.Address, r.MenuUrl, r.AverageCostForTwo, r.AggregateRating)
		//	responseElements = append(responseElements, ResponseElement{
		//		Title:     fmt.Sprintf("%s", r.Name),
		//		TitleLink: r.MenuUrl,
		//		Fields:    r.Fields,
		//		Color:     "#36a64f",
		//		ThumbUrl:  r.ThumbUrl,
		//	})
		//}
		//val, _ := json.Marshal(Value{LastCount: strconv.Itoa(lastCount), CuisineId: strconv.Itoa(inputCusineId.Value)})
		//action := Action{
		//	Name:  "Get",
		//	Text:  "Get",
		//	Type:  "button",
		//	Value: string(val)}
		//buttonAttachment := ResponseElement{Title: "Get next 5 restaurants", Actions: []Action{
		//	action,
		//}, CallbackId: "zomato_next5", Color: "#36a64f"}
		//responseElements = append(responseElements, buttonAttachment)
		//response := SlackResponse{
		//	ResponseType: "ephemeral",
		//	Text:         fmt.Sprintf("The Restaurants for %s cuisine are ", cusineName),
		//	Attachments:  responseElements,
		//}
		//resp, err := json.Marshal(response)
		//if err != nil {
		//	fmt.Fprintf(w, err.Error())
		//}
		resp, err := GetResponseElement(restaurantNames, lastCount, inputCusineId.Value)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s\n", resp)
	}
	//} else {
	//	fmt.Fprintf(w, "Cuisine input is invalid")
	//}
}
