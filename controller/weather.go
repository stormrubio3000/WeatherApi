package controller

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"weather-api/model"

	"github.com/gin-gonic/gin"
)

func getApiKey() (string, error) {
	config := model.Config{}
	//Grab api key from conf.json
	file, err := os.Open("conf.json")

	if err != nil {
		return "", err
	}

	defer file.Close()
	err = json.NewDecoder(file).Decode(&config)

	if err != nil {
		return "", err
	}

	return config.Key, nil
}

func ShowWeather(c *gin.Context) {
	//Error handling checks for valid lat and log input
	long := c.Query("lon")
	lat := c.Query("lat")

	if _, err := strconv.ParseFloat(lat, 32); err != nil || lat == "" {
		c.JSON(http.StatusBadRequest, "latitude value is not valid")
		return
	}

	if _, err := strconv.ParseFloat(long, 32); err != nil || long == "" {
		c.JSON(http.StatusBadRequest, "longitude value is not valid")
		return
	}

	//Send user lat and long to off to check weather and get back results
	weather, err := getWeather(long, lat)

	//Input errors should already be handled so this would most likely be a error with the weather service we are getting data from.
	if err != nil {
		c.JSON(http.StatusInternalServerError, string(err.Error()))
		return
	}

	c.JSON(http.StatusOK, weather)
}

func getWeather(long string, lat string) (model.WeatherOut, error) {
	weatherout := model.WeatherOut{Conditions: "", Tempature: ""}
	var weather model.WeatherGet
	//Grab API key from helper function
	apikey, err := getApiKey()

	if err != nil {
		return weatherout, err
	}

	//Get the full set of data from the weather api as an http response
	res, err := http.Get("https://api.openweathermap.org/data/2.5/weather?lat=" + lat + "&lon=" + long + "&appid=" + apikey)

	if err != nil {
		return weatherout, err
	}

	//Decode json body into our weather struct for processing
	err = json.NewDecoder(res.Body).Decode(&weather)
	if err != nil {
		return weatherout, err
	}

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
