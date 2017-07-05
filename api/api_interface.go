package api
import "sync"

//ServerAPIFactory is the interface for factory design pattern to impelement different 
//servers if needed
type ServerAPIFactory interface {
  StartAndServeAPIs()
}

var (
  fasthttpAPIHandler *FasthttpAPIHandler
  once sync.Once
)

//GetFastHTTPServer serves the API 
func GetFastHTTPServer() ServerAPIFactory{
  once.Do( func (){
      fasthttpAPIHandler = new(FasthttpAPIHandler)
  })
  return fasthttpAPIHandler
}
