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

// FileMetaDatumStorage 注入MongoDB
type FileMetaDatumStorage struct {
	db *BmMongodb.BmMongodb
}

// NewFileMetaDatumStorage initialize parameter
func (s FileMetaDatumStorage) NewFileMetaDatumStorage(args []BmDaemons.BmDaemon) *FileMetaDatumStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &FileMetaDatumStorage{mdb}
}

// GetAll of the model
func (s FileMetaDatumStorage) GetAll(r api2go.Request, skip int, take int) []*Model.FileMetaDatum {
	log.NewLogicLoggerBuilder().Build().Info("Call FileMetaDatum GetAll Parameters ", r.QueryParams)
	in := Model.FileMetaDatum{}
	var out []*Model.FileMetaDatum
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
func (s FileMetaDatumStorage) GetOne(id string) (Model.FileMetaDatum, error) {
	log.NewLogicLoggerBuilder().Build().Info("Call FileMetaDatum GetOne Parameters Id ", id)
	in := Model.FileMetaDatum{ID: id}
	out := Model.FileMetaDatum{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf("FileMetaDatum for id %s not found", id)
	return Model.FileMetaDatum{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a fresh one
func (s *FileMetaDatumStorage) Insert(c Model.FileMetaDatum) string {
	log.NewLogicLoggerBuilder().Build().Info("Call FileMetaDatum Insert Model ", c)
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *FileMetaDatumStorage) Delete(id string) error {
	log.NewLogicLoggerBuilder().Build().Info("Call FileMetaDatum Delete Parameters Id ", id)
	in := Model.FileMetaDatum{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("FileMetaDatum with id %s does not exist", id)
	}

	return nil
}

// Update updates an existing model
func (s *FileMetaDatumStorage) Update(c Model.FileMetaDatum) error {
	log.NewLogicLoggerBuilder().Build().Info("Call FileMetaDatum Update Model ", c)
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("FileMetaDatum with id does not exist")
	}

	return nil
}

// Count MongoDB Query amount
func (s *FileMetaDatumStorage) Count(req api2go.Request, c Model.FileMetaDatum) int {
	r, _ := s.db.Count(req, &c)
	return r
}
