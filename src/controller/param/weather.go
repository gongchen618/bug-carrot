package param

type ResponseWeather struct {
	Result []weatherResult `json:"results"`
}

type weatherResult struct {
	Location weatherResultLocation `json:"location"`
	Daily    []weatherResultDaily  `json:"daily"`
}

type weatherResultLocation struct {
	Name string `yaml:"name"`
}

type weatherResultDaily struct {
	Date      string `json:"date" dismiss:"true"`
	TextDay   string `json:"text_day" text:"早间"`
	TextNight string `json:"text_night" text:"晚间"`
	High      string `json:"high" text:"最高温度"`
	Low       string `json:"low" text:"最低温度"`
	Precip    string `json:"precip" text:"降水概率"`
	Rainfall  string `json:"rainfall" text:"降雨量"`
	WindScale string `json:"wind_scale" text:"风力等级"`
}
