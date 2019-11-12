// Handler
package Handler

import (
	"SandBox/Model"
	http2 "SandBox/Util/http"
	"encoding/json"
	"github.com/PharbersDeveloper/bp-go-lib/log"
	"github.com/alfredyang1986/BmServiceDef/BmConfig"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/alfredyang1986/blackmirror/bmkafka"
	"github.com/elodina/go-avro"
	kafkaAvro "github.com/elodina/go-kafka-avro"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/xid"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type PutHDFSHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	db         *BmMongodb.BmMongodb
	rd         *BmRedis.BmRedis
}

func (h PutHDFSHandler) NewPutHDFSHandler(args ...interface{}) PutHDFSHandler {
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

	return PutHDFSHandler{ Method: md, HttpMethod: hm, Args: ag, db: m, rd: r }
}

func (h PutHDFSHandler) StreamOss2HDFS(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	params := map[string]string{}
	data, _ := ioutil.ReadFile(os.Getenv("STREAMOSS2HDFSAVROCONF"))
	rawMetricsSchema := strings.ReplaceAll(strings.ReplaceAll(string(data), "\n", ""), " ", "")
	bkc, _ := bmkafka.GetConfigInstance()
	res, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(res, &params)
	schema, _ := avro.ParseSchema(rawMetricsSchema)
	record := avro.NewGenericRecord(schema)
	// 居然没有三元表达式，这根Go的简单出发点不一样啊
	var (
		traceId string
	)
	if os.Getenv("mode") == "dev" {
		traceId = xid.New().String()
	} else {
		traceId =  params["traceId"]
	}

	record.Set("jobId", "")
	record.Set("traceId", traceId)
	record.Set("ossKey", params["ossKey"])
	record.Set("fileType", params["fileType"])

	encoder := kafkaAvro.NewKafkaAvroEncoder(bkc.SchemaRepositoryUrl)
	recordByteArr, _ := encoder.Encode(record)

	topic := "oss_task_submit"
	bkc.Produce(&topic, recordByteArr)

	log.NewLogicLoggerBuilder().SetTraceId(traceId).
		Build().Info("StreamOss2HDFS Success, File Where ", params["ossKey"])
	return 0
}

func (h PutHDFSHandler) UpdateJobIDWithTraceID(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	params := map[string]interface{}{}
	res, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(res, &params)

	traceId := params["traceId"].(string)

	in := Model.FileMetaDatum{}
	out := Model.FileMetaDatum{}
	cond := bson.M{ "trace-id": traceId }
	err := h.db.FindOneByCondition(&in, &out, cond)
	if err != nil {

		return 1
	} else {
		out.JobID = append(out.JobID, params["jobId"].(string))
		err = h.db.Update(&out)
		if err != nil {
			return 1
		}
	}

	// 删除文件
	//in := Model.FileVersion{}
	//var out []*Model.FileVersion
	//req := api2go.Request{
	//	QueryParams: map[string][]string{},
	//}
	//err := h.db.FindMulti(req, &in, &out, 0, 0)
	//if err != nil { panic(err) }
	//
	////创建OSSClient实例。
	//client, err := oss.New("oss-cn-beijing.aliyuncs.com", "LTAIEoXgk4DOHDGi", "x75sK6191dPGiu9wBMtKE6YcBBh8EI")
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	os.Exit(-1)
	//}
	//
	//bucketName := "pharbers-sandbox"
	//
	//// 获取存储空间。
	//bucket, err := client.Bucket(bucketName)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	os.Exit(-1)
	//}
	//
	//for _, v := range out  {
	//	// 删除单个文件。
	//	err = bucket.DeleteObject(v.Where)
	//	if err != nil {
	//		fmt.Println("Error:", err)
	//		os.Exit(-1)
	//	}
	//}

	return 0
}

func (h PutHDFSHandler) Stream2HDFSFinish(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	params := map[string]string{}
	res, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(res, &params)
	url := h.Args[0]

	traceId := params["traceId"]

	in := Model.FileMetaDatum{}
	out := Model.FileMetaDatum{}
	cond := bson.M{ "trace-id": traceId }
	err := h.db.FindOneByCondition(&in, &out, cond)
	if err != nil {return 1}

	timeLayout := "2006-01-02 15:04:05"
	dataTimeStr := time.Unix(out.Created / 1000, 0).Format(timeLayout)
	pvw := ""
	for _, v := range out.SampleData {
		pvw += v + "<br/>"
	}

	b, _ := ioutil.ReadFile(os.Getenv("EMAIL_TEMPLATE"))
	reg := regexp.MustCompile("\t|\r|\n")
	userName := strings.ReplaceAll(string(b), "**UserName**", out.OwnerName)
	fileName := strings.ReplaceAll(userName, "**FileName**", out.Name)
	fileType := strings.ReplaceAll(fileName, "**FileType**", out.Extension)
	fileLength := strings.ReplaceAll(fileType, "**FileLength**", strconv.Itoa(out.Length) + "行")
	filePreview := strings.ReplaceAll(fileLength, "**FilePreview**", strings.ReplaceAll(pvw, "\"", ""))
	uploadTime := strings.ReplaceAll(filePreview, "**UploadTime**", dataTimeStr)
	html := strings.ReplaceAll(uploadTime, "**HDFSPATH**", "/test/alex/test001/files/jobId=" + out.JobID[0])
	content := reg.ReplaceAllString(html, "")
	//fmt.Println(strings.ReplaceAll(pvw, "\"", ""))

	for _, e := range BmConfig.BmGetConfigMap(os.Getenv(h.Args[1]))["address"].([]interface{}) {
		body := strings.NewReader(`{
			"email": "`+ e.(string) +`",
			"subject": "SandBox文件上传记录",
			"content": "`+ content +`",
			"content-type": "text/html; charset=UTF-8"
		}`)
		http2.Post(url, r.Header, body)
	}
	return 0
}

func (h PutHDFSHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h PutHDFSHandler) GetHandlerMethod() string {
	return h.Method
}
