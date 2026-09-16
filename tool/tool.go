package main

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
) 

type Weather struct{
	City string  `json:"city"`
	Temperature  string `json:"temperature"`
}

type InputParams struct {
	City string `json:"city" jsonschema:"description = name of city"`
}

func GetWeather(_ context.Context,params *InputParams)(string, error){
	WeatherSet := []Weather{
		{City: "北京", Temperature: "28"},
		{City: "上海", Temperature: "35"},
		{City: "大连", Temperature: "0"},
		{City: "深圳", Temperature: "29"},
	}

	for _, t := range WeatherSet{
		if t.City == params.City{
			return t.Temperature, nil
		}
	}
	return "", nil
}

func CreateTool()tool.InvokableTool{
	GetWeatherTool := utils.NewTool(
		&schema.ToolInfo{
			Name: "getweather",
			Desc: "get the weather of a city",
			ParamsOneOf: schema.NewParamsOneOfByParams(
				map[string]*schema.ParameterInfo{
					"city": &schema.ParameterInfo{
						Type: schema.String,
						Required: true,
					},
				},
			),
		},GetWeather)

	return GetWeatherTool

}
