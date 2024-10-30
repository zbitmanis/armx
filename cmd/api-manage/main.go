package main

import (
   "fmt"
   "log"
   "github.com/zbitmanis/armx"
)

//const apiRouteMangerURL="https://jsonplaceholder.typicode.com/posts"
  const (
  capiAminURL="http://localhost:9180/"
  capiRouteURL="/apisix/admin/routes"
  capiUpsteamPath="/apisix/admin/upstream"
  capiKeyName="X-API-KEY"
  capiKey="************"
  capiId="90b94276"
  )

func main() {

 aroutes,err := armx.getEnrichedRoutes(capiAminURL,capiKeyName,capiKey)
 
 for _, el := range aroutes.Route {
  if el.Value.Id == capiId {
    fmt.Println("Patching: ", capiId)
    armx.patchRouteHost (capiAminURL, capiKeyName, capiKey, el.Value.Id, el.Value.UpstreamFQHost)
   }
 }
 armx.inspectRoutes (capiAminURL, capiKeyName, capiKey)
 
 if err != nil {
    log.Fatalln(err)
    fmt.Println(err)
    return
 }
 //showRoutes(aroutes)
 

 //fmt.Println(resp)
}
