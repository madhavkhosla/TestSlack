package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
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
	http.Handle("/", http.FileServer(http.Dir("./www")))
	http.HandleFunc("/oauth", OAuth)
	http.HandleFunc("/init", GetRestaurants)
	http.HandleFunc("/five", GetFive)
}

func OAuth(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	code := r.FormValue("code")
	req, err := http.NewRequest("GET", "https://slack.com/api/oauth.access", nil)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	query := req.URL.Query()
	query.Add("code", code)
	query.Add("client_id", "189197742244.189972746583")
	query.Add("client_secret", "1a86a133e9457e42aa700f2f7c5a665d")
	req.URL.RawQuery = query.Encode()

	client := urlfetch.Client(ctx)
	//client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	fmt.Fprintf(w, "Successfully installed Food Bot NY")
}

func GetFive(w http.ResponseWriter, r *http.Request) {
	payload := r.FormValue("payload")
	interactiveRequestMessage := InteractiveMessageRequest{}
	err := json.Unmarshal([]byte(payload), &interactiveRequestMessage)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Unable to understand user action."))
		return
	}
	v := RestaurantStat{}
	err = json.Unmarshal([]byte(interactiveRequestMessage.Actions[0].Value), &v)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Unable to understand user action."))
		return
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
		fmt.Fprintf(w, fmt.Sprintf("Unable to determine cuisine information. Please contact using support details"))
		return
	} else {
		resultRemaining := v.CountRemaining + 5
		restaurantStat := &RestaurantStat{LastCount: lastCount, CuisineId: cusId, CountRemaining: resultRemaining}
		resp, err := GetResponseElement(restaurantNames, restaurantStat)
		if err != nil {
			fmt.Fprintf(w, fmt.Sprintf("Unable to create a slack response. Please contact using support details"))
			return
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
		fmt.Fprintf(w, fmt.Sprintf("Unable to determine cuisine information. Please contact using support details"))
		return
	} else {
		resultRemaining := v.CountRemaining - 5
		restaurantStat := &RestaurantStat{LastCount: lastCount, CuisineId: cusId, CountRemaining: resultRemaining}
		resp, err := GetResponseElement(restaurantNames, restaurantStat)
		if err != nil {
			fmt.Fprintf(w, fmt.Sprintf("Unable to create a slack response. Please contact using support details"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s\n", resp)
	}
}

func GetRestaurants(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	cusineName := strings.Title(r.FormValue("text"))
	if cusineName == "Help" {
		response := SlackResponse{
			ResponseType: "ephemeral",
			Text:         fmt.Sprintf("The command lists restaurants in NYC.\n Eg: /mycuisine Italian"),
		}
		resp, err := json.Marshal(response)
		if err != nil {
			fmt.Fprintf(w, fmt.Sprintf("Error rendering help on screen. Please contact using support details"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s\n", resp)
		return
	}
	inputCusineId, err := ConvertNameID(cusineName, ctx)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Unable to determine cuisine information. Please contact using support details"))
		return
	}
	start := 0
	restaurantNames, resultsFound, err := GetRestaurantNamesInCityByCuisine(ctx, inputCusineId.Value, start)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Unable to determine cuisine information. Please contact using support details"))
		return
	} else {
		lastCount := start + 5
		resultRemaining := resultsFound - lastCount
		restaurantStat := &RestaurantStat{LastCount: lastCount, CuisineId: inputCusineId.Value, CountRemaining: resultRemaining}
		resp, err := GetResponseElement(restaurantNames, restaurantStat)
		if err != nil {
			fmt.Fprintf(w, fmt.Sprintf("Unable to create a slack response. Please contact using support details"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s\n", resp)
	}
}
