// Package DataStorage
package DataStorage

import (
	"SandBox/Model"
	"errors"
	"fmt"
	"github.com/PharbersDeveloper/bp-go-lib/log"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/manyminds/api2go"
	"net/http"
)

// FileVersionStorage 注入MongoDB
type FileVersionStorage struct {
	db *BmMongodb.BmMongodb
}

// NewFileVersionStorage initialize parameter
func (s FileVersionStorage) NewFileVersionStorage(args []BmDaemons.BmDaemon) *FileVersionStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &FileVersionStorage{mdb}
}

// GetAll of the model
func (s FileVersionStorage) GetAll(r api2go.Request, skip int, take int) []*Model.FileVersion {
	log.NewLogicLoggerBuilder().Build().Info("Call FileVersion GetAll Parameters ", r.QueryParams)
	in := Model.FileVersion{}
	var out []*Model.FileVersion
	err := s.db.FindMulti(r, &in, &out, skip, take)
	if err == nil {
		for i, iter := range out {
			s.db.ResetIdWithId_(iter)
			out[i] = iter
		}
		return out
	} else {
		return nil
	}
}

// GetOne tasty model
func (s FileVersionStorage) GetOne(id string) (Model.FileVersion, error) {
	log.NewLogicLoggerBuilder().Build().Info("Call FileVersion GetOne Parameters Id ", id)
	in := Model.FileVersion{ID: id}
	out := Model.FileVersion{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf(" FileVersion for id %s not found", id)
	return Model.FileVersion{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a fresh one
func (s *FileVersionStorage) Insert(c Model.FileVersion) string {
	log.NewLogicLoggerBuilder().Build().Info("Call FileVersion Insert Model ", c)
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *FileVersionStorage) Delete(id string) error {
	log.NewLogicLoggerBuilder().Build().Info("Call FileVersion Delete Parameters Id ", id)
	in := Model.FileVersion{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf(" FileVersion with id %s does not exist", id)
	}

	return nil
}

// Update updates an existing model
func (s *FileVersionStorage) Update(c Model.FileVersion) error {
	log.NewLogicLoggerBuilder().Build().Info("Call FileVersion Update Model ", c)
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf(" FileVersion with id does not exist")
	}

	return nil
}

// Count MongoDB Query amount
func (s *FileVersionStorage) Count(req api2go.Request, c Model.FileVersion) int {
	r, _ := s.db.Count(req, &c)
	return r
}
