package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"strconv"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
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
	http.HandleFunc("/", handler)
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

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	cusineName := r.FormValue("text")
	//fmt.Println(cusineName)
	baseUrl := "https://developers.zomato.com/api/v2.1/"
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/cuisines", baseUrl), nil)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Key", "a881a2c3cbd4e8320634917542051763")
	query := req.URL.Query()
	query.Add("city_id", "280")
	req.URL.RawQuery = query.Encode()
	client := urlfetch.Client(ctx)
	//client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	collection := Collection{}
	err = json.NewDecoder(resp.Body).Decode(&collection)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	var inputCusineId int = -1
	for _, c := range collection.Cuisines {
		if c.Cuisine.CuisineName == cusineName {
			inputCusineId = c.Cuisine.CuisineId
			//fmt.Fprintf(w, "%d\n", c.Cuisine.CuisineId)
		}
	}
	if inputCusineId != -1 {
		restaurantNames, lastCount, err := GetRestaurantNamesInCityByCuisine(ctx, inputCusineId, 0)
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
			val, _ := json.Marshal(Value{LastCount: strconv.Itoa(lastCount), CuisineId: strconv.Itoa(inputCusineId)})
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
				Text:         fmt.Sprintf("The Restaurants for %s cuisine are ", cusineName),
				Attachments:  responseElements,
			}
			resp, err := json.Marshal(response)
			if err != nil {
				fmt.Fprintf(w, err.Error())
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "%s\n", resp)
		}
	} else {
		fmt.Fprintf(w, "Cuisine input is invalid")
	}
}
