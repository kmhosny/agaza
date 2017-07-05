package api

import (
	"encoding/json"
	"agaza/config"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

const (
	apiBase     = "/api"
	rootURI     = apiBase + "/"
)

var (
	configuration = config.LoadConfiguration()
)

//FasthttpAPIHandler is structure wrapper for the fast http objects
type FasthttpAPIHandler struct {
	router *routing.Router
}

//StartAndServeAPIs will start the server and attach the endpoints to the router
func (f *FasthttpAPIHandler) StartAndServeAPIs() {
	f.router = routing.New()
	f.router.Get(rootURI,f.SetHeaders, f.Logger, f.GetRoot)
	fasthttp.ListenAndServe(":"+configuration.APIserverport, f.router.HandleRequest)
}


//GetRoot the root page
func (f *FasthttpAPIHandler) GetRoot(c *routing.Context) error {
	agaza, _ := json.Marshal("agaza")
	c.Write(agaza)
	return nil
}
