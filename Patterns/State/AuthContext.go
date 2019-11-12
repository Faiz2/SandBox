package State

import (
	"SandBox/Model"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/alfredyang1986/blackmirror/bmlog"
)

type AuthContext struct {
	state State
	db *BmMongodb.BmMongodb
	rd *BmRedis.BmRedis
	target string
	who map[string]string
	action string
}

func (c *AuthContext) NewAuthContext(target,
									action string,
									who map[string]string,
									db *BmMongodb.BmMongodb,
									rd *BmRedis.BmRedis) {
	c.target = target
	c.who = who
	c.action = action
	c.db = db
	c.rd = rd
}

func (c *AuthContext) mappingState() {
	if c.db == nil || c.rd == nil {
		panic("MongodbDriver Or RedisDriver is Nil")
	}

	out := Model.FileMetaDatum{}
	err := c.db.FindOne(&Model.FileMetaDatum{ID: c.target}, &out)
	if err != nil {
		bmlog.StandardLogger().Error(err.Error())
		panic(err.Error())
	}

	if c.who["account"] == out.OwnerID {
		c.state = &OwnerState{c, c.db, c.rd, out}
	} else if c.who["group"] == out.GroupID {
		c.state = &GroupState{c, c.db, c.rd, out}
	} else {
		c.state = &OtherState{c, c.db, c.rd, out}
	}
}

func (c *AuthContext) DoExecute() (interface{}, error) {
	c.mappingState()
	return c.state.DoExecute()
	//return nil, nil
}
