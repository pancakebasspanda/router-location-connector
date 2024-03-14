package api

type Router struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	LocationID  int    `json:"location_id"`
	RouterLinks []int  `json:"router_links"`
}

type Location struct {
	ID       int    `json:"id"`
	Postcode string `json:"postcode"`
	Name     string `json:"name"`
}

type RouterLocationData struct {
	Routers   []Router   `json:"routers"`
	Locations []Location `json:"locations"`
}

type RouterLocationLink struct {
	UniqueID   string // sort two locations alphabetically and concatenate them
	Connection string
}
