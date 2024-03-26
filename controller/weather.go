package controller

import (
	"encoding/json"
	"net/http"
	"os"
	"weather-api/model"

	"github.com/gin-gonic/gin"
)

func getApiKey() string {
	config := model.Config{}
	//Grab api key from conf.json
	file, _ := os.Open("conf.json")
	defer file.Close()
	err := json.NewDecoder(file).Decode(&config)

	if err != nil {
		panic(err)
	}
	return config.Key
}

func ShowWeather(c *gin.Context) {
	//Grab the latitude and longitude from the url in the get request and send to data manip func
	weather, err := getWeather(c.Query("lon"), c.Query("lat"))

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, weather)
}

func getWeather(long string, lat string) (model.WeatherOut, error) {
	var weather model.WeatherGet
	//Grab API key from helper function
	apikey := getApiKey()
	//Get the full set of data from the weather api as an http response
	res, err := http.Get("https://api.openweathermap.org/data/2.5/weather?lat=" + lat + "&lon=" + long + "&appid=" + apikey)

	if err != nil {
		panic(err)
	}

	//Decode json body into our weather struct for processing
	err = json.NewDecoder(res.Body).Decode(&weather)
	if err != nil {
		panic(err)
	}

	weatherout := model.WeatherOut{Conditions: "", Tempature: ""}

	//weather conditions are an array so we go through and add all to out outgoing struct.
	for _, weath := range weather.Weather {
		if weatherout.Conditions == "" {
			weatherout.Conditions = weath.Main
			continue
		}
		weatherout.Conditions = weatherout.Conditions + " and " + weath.Main
	}

	//Small switch to check temp for hot and cold. 85+(f) is hot and 50-(f) is cold
	switch {
	case weather.Main.FeelsLike > 298:
		weatherout.Tempature = " Tempature is hot."
	case weather.Main.FeelsLike < 283:
		weatherout.Tempature = " Tempature is cold."
	default:
		weatherout.Tempature = " Tempature is moderate."
	}

	return weatherout, nil
}
