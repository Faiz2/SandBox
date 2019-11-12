// Package Model FileVersion
package Model

import (
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

// FileVersion 文件版本
type FileVersion struct{
	ID		string			`json:"-"`
	Id_		bson.ObjectId	`json:"-" bson:"_id"`
	Parent	string			`json:"parent" bson:"parent"`
	Size	int64			`json:"size" bson:"size"`
	Where	string			`json:"where" bson:"where"`
	Kind	string			`json:"kind" bson:"kind"`
	Tag		string			`json:"tag" bson:"tag"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (f FileVersion) GetID() string {
	return f.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (f *FileVersion) SetID(id string) error {
	f.ID = id
	return nil
}

func (f *FileVersion) GetConditionsBsonM(parameters map[string][]string) bson.M {
	rst := make(map[string]interface{})
	for k, v := range parameters {
		switch k {
		case "ids":
			r := make(map[string]interface{})
			var ids []bson.ObjectId
			for i := 0; i < len(v); i++ {
				ids = append(ids, bson.ObjectIdHex(v[i]))
			}
			r["$in"] = ids
			rst["_id"] = r
		case "code":
			v, _ := strconv.Atoi(v[0])
			rst[k] = v

		}
	}
	return rst
}