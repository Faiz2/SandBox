package test

import (
	"SandBox/Patterns/State"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/alfredyang1986/blackmirror/bmlog"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAuth(t *testing.T) {
	t.Parallel()
	Convey("SandBox Auth Test", t, func() {
		db := &BmMongodb.BmMongodb {
			Host:     "127.0.0.1",
			Port:     "27017",
			User:     "",
			Pass:     "",
			Database: "pharbers-sandbox",
		}
		rd := &BmRedis.BmRedis {
			Host:     "127.0.0.1",
			Port:     "6379",
			Password: "",
			Database: "",
		}
		context := State.AuthContext{}
		who := map[string]string {
			"account": "5cd51df9f4ce43ee2495d4dd",
			"group": "5cc3f7f7ceb3c45854b80e25",
		}
		context.NewAuthContext("5d3299f4421aa93290f1c91a", "r", who,  db, rd)
		result, err := context.DoExecute()
		bmlog.StandardLogger().Info(result)
		bmlog.StandardLogger().Info(err)
	})
}
