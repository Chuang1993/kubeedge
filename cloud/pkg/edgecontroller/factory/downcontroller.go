package factory

type DownController struct {
	resources []*DownFlow
	router    *Router
}

func (dc *DownController) Start() {
	for _, dr := range dc.resources {
		dr.Router = dc.router
		go dr.Start()
	}
}
