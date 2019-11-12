package env

import (
	"fmt"
	"flag"
	"github.com/PharbersDeveloper/bp-go-lib/log"
	"os"
)

const (
	//Project env key
	ProjectName = "PROJECT_NAME"

	//Log env key
	LogTimeFormat = "BP_LOG_TIME_FORMAT"
	LogOutput     = "BP_LOG_OUTPUT"
	LogLevel      = "BP_LOG_LEVEL"

	//kafka env key
	KafkaConfigEnable = "BP_KAFKA_CONFIG_ENABLE"
	KafkaConfigPath = "BP_KAFKA_CONFIG_PATH"
)

func SetLocalEnv() {
	//项目范围内的环境变量
	_ = os.Setenv(ProjectName, "SandBox")

	//log
	_ = os.Setenv(LogTimeFormat, "2006-01-02 15:04:05")
	_ = os.Setenv(LogOutput, "console")
	//_ = os.Setenv(env.LogOutput, "/Users/qianpeng/bplogs/mqtt-message-storage.log")
	_ = os.Setenv(LogLevel, "info")
	_ = os.Setenv("BM_KAFKA_CONF_HOME", fmt.Sprint(os.Getenv("SANDBOX_HOME"), "/resource/kafkaconfig.json"))
	_ = os.Setenv("HDFSAVROCONF", fmt.Sprint(os.Getenv("SANDBOX_HOME"), "/resource/hdfs-avro.json"))
	_ = os.Setenv("STREAMOSS2HDFSAVROCONF", fmt.Sprint(os.Getenv("SANDBOX_HOME"), "/resource/stream-oss2hdfs-avro.json"))
	_ = os.Setenv("EMAIL_TEMPLATE", fmt.Sprint(os.Getenv("SANDBOX_HOME"), "/resource/email-template.txt"))
	_ = os.Setenv("EMAILADDRESS", fmt.Sprint(os.Getenv("SANDBOX_HOME"), "/resource/emails.json"))

	//kafka
	//_ = os.Setenv(KafkaConfigPath, "../resources/kafka_config.json")
}

func SetStartingParameter() {
	mode := flag.String("mode", "prod", "Development Mode")
	flag.Parse()
	_ = os.Setenv("mode", *mode)
	log.NewLogicLoggerBuilder().Build().Info("Starting Parameter，Development Mode ", *mode)
}
