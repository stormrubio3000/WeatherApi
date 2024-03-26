Simple weather api using here:https://openweathermap.org/api.

API key should be stored in the root dir of the project in a file named conf.json in the below format.
{
    "api-key": "api-key-here"
}

App checks runs on localhost:8080 and the single endpoint is on localhost:8080/. The endpoint requires lat and lon as params to work.
Example format for the endpoint will look like localhost:8080/?lon=44.34&lat=10.99

Return value is below
{
    "Conditions": "",
    "Tempature": ""
}

Conditions will list all of the weather conditions that are currently active in the location selected.
Tempature is checking for under 50(f) for cold and over 85(f) for hot with the temps in between as moderate.
