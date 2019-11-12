// Package Resource
package Resource

import (
	"SandBox/DataStorage"
	"SandBox/Model"
	"errors"
	"github.com/alfredyang1986/BmServiceDef/BmDataStorage"
	"github.com/manyminds/api2go"
	"net/http"
	"reflect"
	"strconv"
)

// FileMetaDatumResource
type FileMetaDatumResource struct {
	FileMetaDatumStorage   	*DataStorage.FileMetaDatumStorage
	SandBoxIndexStorage		*DataStorage.SandBoxIndexStorage
}

// NewFileMetaDatumResource Initialize Parameter And injection Storage Or Resource
func (s FileMetaDatumResource) NewFileMetaDatumResource(args []BmDataStorage.BmStorage) *FileMetaDatumResource {
	var dcs *DataStorage.FileMetaDatumStorage
	var sbi *DataStorage.SandBoxIndexStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "FileMetaDatumStorage" {
			dcs = arg.(*DataStorage.FileMetaDatumStorage)
		} else if tp.Name() == "SandBoxIndexStorage" {
			sbi = arg.(*DataStorage.SandBoxIndexStorage)
		}
	}
	return &FileMetaDatumResource{
		FileMetaDatumStorage: dcs,
		SandBoxIndexStorage: sbi,
	}
}

func (s FileMetaDatumResource) FindAll(r api2go.Request) (api2go.Responder, error) {

	sandBoxIndicesID, sok := r.QueryParams["sandBoxIndicesID"]

	if sok {
		modelRootID := sandBoxIndicesID[0]

		modelRoot, err := s.SandBoxIndexStorage.GetOne(modelRootID)

		if err != nil {
			return  &Response{}, nil
		}

		r.QueryParams["ids"] = modelRoot.FileMetaDatumIDs

		result := s.FileMetaDatumStorage.GetAll(r, -1, -1)

		return &Response{Res: result}, nil
	}

	if len(r.QueryParams["group-id"]) > 0 {
		result := s.FileMetaDatumStorage.GetAll(r, -1, -1)
		return &Response{Res: result}, nil
	}

	// TODO: 目前只对Group进行查询
	//companyID, cok := r.QueryParams["companyID"]
	//if cok {
	//	fmt.Println(companyID)
	//}

	return &Response{}, nil
}

// PaginatedFindAll can be used to load models in chunks
func (s FileMetaDatumResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	var (
		result                      []Model.FileMetaDatum
		number, size, offset, limit string
	)

	numberQuery, ok := r.QueryParams["page[number]"]
	if ok {
		number = numberQuery[0]
	}
	sizeQuery, ok := r.QueryParams["page[size]"]
	if ok {
		size = sizeQuery[0]
	}
	offsetQuery, ok := r.QueryParams["page[offset]"]
	if ok {
		offset = offsetQuery[0]
	}
	limitQuery, ok := r.QueryParams["page[limit]"]
	if ok {
		limit = limitQuery[0]
	}

	if size != "" {
		sizeI, err := strconv.ParseInt(size, 10, 64)
		if err != nil {
			return 0, &Response{}, err
		}

		numberI, err := strconv.ParseInt(number, 10, 64)
		if err != nil {
			return 0, &Response{}, err
		}

		start := sizeI * (numberI - 1)
		for _, iter := range s.FileMetaDatumStorage.GetAll(r, int(start), int(sizeI)) {
			result = append(result, *iter)
		}

	} else {
		limitI, err := strconv.ParseUint(limit, 10, 64)
		if err != nil {
			return 0, &Response{}, err
		}

		offsetI, err := strconv.ParseUint(offset, 10, 64)
		if err != nil {
			return 0, &Response{}, err
		}

		for _, iter := range s.FileMetaDatumStorage.GetAll(r, int(offsetI), int(limitI)) {
			result = append(result, *iter)
		}
	}

	in := Model.FileMetaDatum{}
	count := s.FileMetaDatumStorage.Count(r, in)

	return uint(count), &Response{Res: result}, nil
}

// FindOne to satisfy `api2go.DataSource` interface
// this method should return the model with the given ID, otherwise an error
func (s FileMetaDatumResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	modelRoot, err := s.FileMetaDatumStorage.GetOne(ID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}

	return &Response{Res: modelRoot}, nil
}

// Create method to satisfy `api2go.DataSource` interface
func (s FileMetaDatumResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(Model.FileMetaDatum)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid_Instance_Given"), "Invalid Instance Given", http.StatusBadRequest)
	}

	id := s.FileMetaDatumStorage.Insert(model)
	model.ID = id

	return &Response{Res: model, Code: http.StatusCreated}, nil
}

// Delete to satisfy `api2go.DataSource` interface
func (s FileMetaDatumResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := s.FileMetaDatumStorage.Delete(id)
	return &Response{Code: http.StatusNoContent}, err
}

//Update stores all changes on the model
func (s FileMetaDatumResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(Model.FileMetaDatum)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid_Instance_Given"), "Invalid Instance Given", http.StatusBadRequest)
	}
	err := s.FileMetaDatumStorage.Update(model)
	return &Response{Res: model, Code: http.StatusNoContent}, err
}
