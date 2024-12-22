package model

type WeewxLoopPayload struct {
	DateTime             string `json:"dateTime"`
	DaymaxwindMph        string `json:"daymaxwind_mph"`
	ExtraHumid2          string `json:"extraHumid2"`
	ExtraTemp2F          string `json:"extraTemp2_F"`
	HeapFreeByte         string `json:"heap_free_byte"`
	LuminosityLux        string `json:"luminosity_lux"`
	OutHumidity          string `json:"outHumidity"`
	OutTempF             string `json:"outTemp_F"`
	PDayRainIn           string `json:"p_dayRain_in"`
	PMonthRainIn         string `json:"p_monthRain_in"`
	PRainIn              string `json:"p_rain_in"`
	PRainRateInchPerHour string `json:"p_rainRate_inch_per_hour"`
	PStormRainIn         string `json:"p_stormRain_in"`
	PWeekRainIn          string `json:"p_weekRain_in"`
	PYearRainIn          string `json:"p_yearRain_in"`
	PressureInHg         string `json:"pressure_inHg"`
	RelbarometerInHg     string `json:"relbarometer_inHg"`
	Uv                   string `json:"UV"`
	UvradiationWpm2      string `json:"uvradiation_Wpm2"`
	Wh31Ch2BattCount     string `json:"wh31_ch2_batt_count"`
	Wh31Ch2SigCount      string `json:"wh31_ch2_sig_count"`
	WindDir              string `json:"windDir"`
	WindGustMph          string `json:"windGust_mph"`
	WindSpeedMph         string `json:"windSpeed_mph"`
	Ws90SigCount         string `json:"ws90_sig_count"`
	TxBatteryStatus      string `json:"txBatteryStatus"`
	RxCheckPercent       string `json:"rxCheckPercent"`
	UsUnits              string `json:"usUnits"`
	RadiationWpm2        string `json:"radiation_Wpm2"`
	AltimeterInHg        string `json:"altimeter_inHg"`
	AppTempF             string `json:"appTemp_F"`
	BarometerInHg        string `json:"barometer_inHg"`
	CloudbaseFoot        string `json:"cloudbase_foot"`
	DewpointF            string `json:"dewpoint_F"`
	HeatindexF           string `json:"heatindex_F"`
	HumidexF             string `json:"humidex_F"`
	MaxSolarRadWpm2      string `json:"maxSolarRad_Wpm2"`
	RainRateInchPerHour  string `json:"rainRate_inch_per_hour"`
	WindchillF           string `json:"windchill_F"`
}

type Config struct {
	AllskyHost     string `json:"allsky_host"`
	AllskyPort     int    `json:"allsky_port"`
	AllskyClientID string `json:"allsky_client_id"`
	AllskyUsername string `json:"allsky_username"`
	AllskyPassword string `json:"allsky_password"`
	AllskyTopic    string `json:"allsky_topic"`
	WxHost         string `json:"wx_host"`
	WxPort         int    `json:"wx_port"`
	WxClientID     string `json:"wx_client_id"`
	WxUsername     string `json:"wx_username"`
	WxPassword     string `json:"wx_password"`
	WxTopic        string `json:"wx_topic"`
}
