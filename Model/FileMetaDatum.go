// Package Model FileMetaDatum
package Model

import (
	"errors"
	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/mgo.v2/bson"
)

// FileMetaDatum 文件源数据
type FileMetaDatum struct {
	ID				string			`json:"-"`
	Id_ 			bson.ObjectId	`json:"-" bson:"_id"`
	Name			string			`json:"name" bson:"name"`
	Extension		string			`json:"extension" bson:"extension"`
	Created 		int64			`json:"created" bson:"created"`
	Kind 			string			`json:"kind" bson:"kind"`
	OwnerID			string			`json:"owner-id" bson:"owner-id"`
	OwnerName		string			`json:"owner-name" bson:"owner-name"`
	GroupID			string			`json:"group-id" bson:"group-id"`
	Mod 			string			`json:"mod" bson:"mod"`
	FileVersionIDs	[]string		`json:"-" bson:"file-version-ids"`
	FileVersions	[]*FileVersion	`json:"-"`
	JobID			[]string		`json:"job-id" bson:"job-id"`
	TraceID			string 			`json:"trace-id" bson:"trace-id"`
	Length			int				`json:"length" bson:"length"`

	// TODO：结构不对，应该抽出来一张表，为了能及时演示
	SampleData		[]string 		`json:"-" bson:"sample-data"`
	Schema 			string 			`json:"-" bson:"schema"`

}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (f FileMetaDatum) GetID() string {
	return f.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (f *FileMetaDatum) SetID(id string) error {
	f.ID = id
	return nil
}

// GetReferences 设置关联关系详情参考JSONAPI定义的Reference
// https://jsonapi.org/format/#document-resource-object-relationships
func (f FileMetaDatum) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "fileVersions",
			Name: "fileVersions",
		},
	}
}

// GetReferencedIDs 获取关联ID
func (f FileMetaDatum) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}

	for _, kID := range f.FileVersionIDs {
		result = append(result, jsonapi.ReferenceID{
			ID:   kID,
			Type: "fileVersions",
			Name: "fileVersions",
		})
	}
	
	return result
}

// GetReferencedStructs 设置Reference内容
func (f FileMetaDatum) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}

	for key := range f.FileVersions {
		result = append(result, f.FileVersions[key])
	}
	
	return result
}

// SetToManyReferenceIDs 设置关联ID
func (f *FileMetaDatum) SetToManyReferenceIDs(name string, IDs []string) error {
	if name == "fileVersions" {
		f.FileVersionIDs = IDs
		return nil
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// AddToManyIDs
func (f *FileMetaDatum) AddToManyIDs(name string, IDs []string) error {
	if name == "fileVersions" {
		f.FileVersionIDs = append(f.FileVersionIDs, IDs...)
		return nil
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// GetConditionsBsonM MongoDB的拼接Condition
func (f *FileMetaDatum) GetConditionsBsonM(parameters map[string][]string) bson.M {
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
		case "group-id":
			rst[k] = v[0]
		case "kind":
			rst[k] = v[0]
		}
	}
	return rst
}