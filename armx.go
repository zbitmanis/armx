package armx

import (
   "encoding/json"
   "fmt"
   "bytes"
   "io"
   "log"
   "strings"
   "net/http"
)

type apiValue struct {
  CreateTime int `json:"create_time"`
  PluginConfigId string `json:"plugin_config_id"`
  Status int `json:"status"`
  Uris []string `json:"uris"`
  UpstreamId string `json:"upstream_id"`
  Labels map[string]interface{} `json:"labels"`
  UpdateTime int `json:"update_time"`
  Host string `json:"host"`
  Desc string `json:"desc"`
  Priority int `json:"priority"`
  Name string `json:"name"`
  Plugins map[string]interface{}  `json:"plugins"`
  Id string `json:"id"`
  Upstream apiUpstream
  UpstreamFQHost string
  UpstreamServiceName string
  UpstreamNameSpace string
  UpstreamServicePort string

}

type apiUpstreamValue struct {
  CreateTime int `json:"create_time"`
  Id string `json:"id"`
  Labels map[string]interface{} `json:"labels"`
  UpdateTime int `json:"update_time"`
  Type string `json:"type"`
  PassHost string `json:"pass_host"`
  Nodes []map[string]interface{}  `json:"nodes"`
  HashOn string `json:"hash_on"`
  Desc string `json:"desc"`
  Name string `json:"name"`
  Scheme string `json:"scheme"`
}

type apiUpstream struct {
  Key   string `json:"key"`
  CreatedIndex  int `json:"createdIndex"`
  Value apiUpstreamValue `json:"value"`
}

type apiUpstreams struct {
  Total int `json:"total"`
  Route []apiUpstream `json:"list"`
}

type apiRoute struct {
  Key   string `json:"key"`
  CreatedIndex  int `json:"createdIndex"`
  Value apiValue `json:"value"`
}

type apiRoutes struct {
  Total int `json:"total"`
  Route []apiRoute `json:"list"`
}


//func ( r * apiRoute ) getRoutes (url string) error {
func getRoutes (apiAminURL string, apiKeyName string, apiKey string)(resp apiRoutes, err error) {
   var aroutes apiRoutes
   var url=apiAminURL+"/apisix/admin/routes"  
   client := http.Client{}
   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      log.Fatalln(err)
   }
   req.Header.Set("Content-Type","application/json")
   req.Header.Set(apiKeyName,apiKey)
   
   res , err := client.Do(req)
   if err != nil {
      log.Fatalln(err)
   }
     
   body, err := io.ReadAll(res.Body)
   if err != nil {
      log.Fatalln(err)
   }
   
   defer res.Body.Close()
   
   if err := json.Unmarshal(body,&aroutes); err != nil {
      log.Fatalln(err)
   }
  
   //showRoutes(aroutes)

   return aroutes, nil
}

func getUpstream (apiAminURL string, apiKeyName string, apiKey string, apiUpstreamKey string)(resp apiUpstream, err error) {
   var aupstream apiUpstream
   var url=apiAminURL+"/apisix/admin/upstreams/"+apiUpstreamKey
   
   client := http.Client{}
   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      log.Fatalln(err)
   }
   req.Header.Set("Content-Type","application/json")
   req.Header.Set(apiKeyName,apiKey)
   
   res , err := client.Do(req)
   if err != nil {
      log.Fatalln(err)
   }
     
   body, err := io.ReadAll(res.Body)
   if err != nil {
      log.Fatalln(err)
   }
   
   defer res.Body.Close()
   
   if err := json.Unmarshal(body,&aupstream); err != nil {
     log.Fatalln(err)
   }

   return aupstream, nil
}

func getEnrichedRoutes(apiAminURL string, apiKeyName string, apiKey string)(resp apiRoutes, err error) {
  aroutes,err := getRoutes (apiAminURL , apiKeyName , apiKey)
  const  cClusterName="cluster.local" 
  if err != nil {
    log.Fatalln(err)
  }
  
  for i:=0 ; i< len(aroutes.Route); i++  {
    el := &aroutes.Route[i]
    aupstream, err := getUpstream (apiAminURL, apiKeyName, apiKey, el.Value.UpstreamId)  
    if err != nil {
      log.Fatalln(err)
    }
    el.Value.Upstream = aupstream
    sArr := strings.Split(aupstream.Value.Name, "_")
    if len(sArr) > 1 {
      el.Value.UpstreamFQHost = sArr[1]+"."+sArr[0]+".svc."+cClusterName+":"+sArr[2]
      el.Value.UpstreamServiceName = sArr[1]
      el.Value.UpstreamNameSpace   = sArr[0]
      el.Value.UpstreamServicePort = sArr[2]
    }
  }
  return aroutes, nil

}

func patchRouteHost (apiAminURL string, apiKeyName string, apiKey string, apiRouteKey string, apiHost string) (err error) {
   var url=apiAminURL+"/apisix/admin/routes/"+apiRouteKey

   value := make(map[string]map[string]map[string]string)
   plugins := make(map[string]map[string]string)
   params := map[string]string{"host": apiHost}
   plugins["proxy-rewrite"]=params
   value["plugins"]=plugins
   
   //jsonStr := []byte(`{"value":{"plugins":{"proxy-rewrite":{"host":"httpbin-apisix.itest.svc.cluster.local:80"}}}}`)
   jsonStr, err := json.Marshal(value) 

   if err != nil {
      log.Fatalln(err)
   }
   fmt.Printf("Patching route:",apiRouteKey,"with: ",string(jsonStr))

   client := http.Client{}
   req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonStr))
   if err != nil {
      log.Fatalln(err)
   }
   req.Header.Set("Content-Type","application/json")
   req.Header.Set(apiKeyName,apiKey)
   
   res , err := client.Do(req)
   if err != nil {
      log.Fatalln(err)
   }
     
   body, err := io.ReadAll(res.Body)
   if err != nil {
      log.Fatalln(err)
   }
   
   defer res.Body.Close()
   
   fmt.Println(string(body)) 
   return nil
}

func inspectRoutes (apiAminURL string, apiKeyName string, apiKey string)(resp apiRoutes, err error) {
  aroutes,err := getRoutes (apiAminURL , apiKeyName , apiKey)
  
  if err != nil {
    log.Fatalln(err)
  }
  
  for idx, el := range aroutes.Route {
    fmt.Println("el: ",idx)
    prettyPrint(el) 
    aupstream, err := getUpstream (apiAminURL, apiKeyName, apiKey, el.Value.UpstreamId)  
    if err != nil {
      log.Fatalln(err)
    }
    prettyPrint(aupstream) 
  }
  return aroutes, nil
}

func showRoutes  (aroutes apiRoutes) {
   for idx, el := range aroutes.Route {
     fmt.Println("el: ",idx)
     fmt.Println("Route:", el.Key)
     prettyPrint(el) 
     fmt.Println("Route values for the route:", el.Key)
     prettyPrint(el.Value) 
   }
}

func showUpstream  (aUp apiUpstream) {
     prettyPrint(aUp) 
}
