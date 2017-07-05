package api

import (
	"time"
	"agaza/logger"
	"github.com/qiangxue/fasthttp-routing"
)

//Logger log the request
func (a *FasthttpAPIHandler) Logger(c *routing.Context) error {
	logger.Trace.Println(string(c.Method())+":", c.URI().String(), time.Now(), c.QueryArgs(), string(c.PostBody()))
	return nil
}

//SetHeaders set common headers
func (a *FasthttpAPIHandler) SetHeaders(c *routing.Context) error {
	c.SetContentType("application/json")
	return nil
}