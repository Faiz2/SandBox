// Handler
package Handler

import (
	http2 "SandBox/Util/http"
	"encoding/json"
	"fmt"
	"github.com/PharbersDeveloper/bp-go-lib/log"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

type topData struct {
	Attributes  map[string]interface{}	`json:"attributes"`
	Relationships relationship 			`json:"relationships"`
}

type attributes struct {
	Name string `json:"name"`
}

type data struct {
	ID string `json:"id"`
}

type relationshipData struct {
	Data data `json:"data"`
}
type relationship struct {
	Employee relationshipData `json:"employee"`
	Group	 relationshipData `json:"group"`
	Company  relationshipData `json:"company"`
}

type oAuthResult struct {
	Data topData `json:"data"`
}

type GetAccountDetailHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	db         *BmMongodb.BmMongodb
	rd         *BmRedis.BmRedis
}

func (h GetAccountDetailHandler) NewGetAccountDetailHandler(args ...interface{}) GetAccountDetailHandler {
	var m *BmMongodb.BmMongodb
	var r *BmRedis.BmRedis
	var hm string
	var md string
	var ag []string
	for i, arg := range args {
		if i == 0 {
			sts := arg.([]BmDaemons.BmDaemon)
			for _, dm := range sts {
				tp := reflect.ValueOf(dm).Interface()
				tm := reflect.ValueOf(tp).Elem().Type()
				if tm.Name() == "BmMongodb" {
					m = dm.(*BmMongodb.BmMongodb)
				}
				if tm.Name() == "BmRedis" {
					r = dm.(*BmRedis.BmRedis)
				}
			}
		} else if i == 1 {
			md = arg.(string)
		} else if i == 2 {
			hm = arg.(string)
		} else if i == 3 {
			lst := arg.([]string)
			for _, str := range lst {
				ag = append(ag, str)
			}
		}
	}

	return GetAccountDetailHandler{ Method: md, HttpMethod: hm, Args: ag, db: m, rd: r }
}

// GetAccountDetail oss or hdfs link return json
func (h GetAccountDetailHandler) GetAccountDetail(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	// @Alex 差个OAuth的Token验证
	w.Header().Add("Content-Type", "application/json")
	params := map[string]string{}
	res, _ := ioutil.ReadAll(r.Body)
	result := map[string]interface{}{}
	enc := json.NewEncoder(w)
	err := json.Unmarshal(res, &params)
	if err != nil {
		log.NewLogicLoggerBuilder().Build().Error(err)
		result["status"] = "error"
		result["msg"] = "解析失败"
		err = enc.Encode(result)
		log.NewLogicLoggerBuilder().Build().Error(err)
		return 1
	}
	accountId, aok := params["accountId"]

	if !aok {
		log.NewLogicLoggerBuilder().Build().Error("Account 参数缺失")
		result["status"] = "error"
		result["msg"] = "Account 参数缺失"
		err = enc.Encode(result)
		log.NewLogicLoggerBuilder().Build().Error(err)
		return 1
	}

	// 拼接转发的URL
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}

	version := strings.Split(r.URL.Path, "/")[1]
	resource := fmt.Sprint(h.Args[0], "/", version, "/", "accounts/", accountId)
	mergeURL := strings.Join([]string{scheme, resource}, "")

	// 转发
	response := http2.Get(mergeURL, r.Header)
	oar := oAuthResult{}
	json.Unmarshal(response, &oar)

	resource = fmt.Sprint(h.Args[0], "/", version, "/", "employees/", oar.Data.Relationships.Employee.Data.ID)
	mergeURL = strings.Join([]string{scheme, resource}, "")

	// 转发
	response = http2.Get(mergeURL, r.Header)
	oar = oAuthResult{}
	json.Unmarshal(response, &oar)
	result["account-name"] = oar.Data.Attributes["name"]
	result["group-id"] = oar.Data.Relationships.Group.Data.ID

	resource = fmt.Sprint(h.Args[0], "/", version, "/", "groups/", oar.Data.Relationships.Group.Data.ID)
	mergeURL = strings.Join([]string{scheme, resource}, "")
	// 转发
	response = http2.Get(mergeURL, r.Header)
	oar = oAuthResult{}
	json.Unmarshal(response, &oar)

	result["company-id"] = oar.Data.Relationships.Company.Data.ID

	resource = fmt.Sprint(h.Args[0], "/", version, "/", "companies/", oar.Data.Relationships.Company.Data.ID)
	mergeURL = strings.Join([]string{scheme, resource}, "")
	// 转发
	response = http2.Get(mergeURL, r.Header)
	oar = oAuthResult{}
	json.Unmarshal(response, &oar)

	result["company-name"] = oar.Data.Attributes["name"]
	result["status"] = "ok"

	enc.Encode(result)
	return 0
}

func (h GetAccountDetailHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h GetAccountDetailHandler) GetHandlerMethod() string {
	return h.Method
}