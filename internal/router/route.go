package router

type route struct {
	name       string
	method     string
	pattern    string
	appHandler AppJSONHandler
}
