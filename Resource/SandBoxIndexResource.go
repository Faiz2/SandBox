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

// SandBoxIndexResource
type SandBoxIndexResource struct {
	SandBoxIndexStorage   *DataStorage.SandBoxIndexStorage
}

// NewSandBoxIndexResource Initialize Parameter And injection Storage Or Resource
func (s SandBoxIndexResource) NewSandBoxIndexResource(args []BmDataStorage.BmStorage) *SandBoxIndexResource {
	var dcs *DataStorage.SandBoxIndexStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "SandBoxIndexStorage" {
			dcs = arg.(*DataStorage.SandBoxIndexStorage)
		}
	}
	return &SandBoxIndexResource{
		SandBoxIndexStorage:	dcs,
	}
}

func (s SandBoxIndexResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	accountId, ok := r.QueryParams["accountId"]
	if ok {
		modelRootID := accountId[0]
		r.QueryParams["account-id"] = []string{modelRootID}
	}

	result := s.SandBoxIndexStorage.GetAll(r, -1, -1)
	return &Response{Res: result}, nil
}

// PaginatedFindAll can be used to load models in chunks
func (s SandBoxIndexResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	var (
		result                      []Model.SandBoxIndex
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
		for _, iter := range s.SandBoxIndexStorage.GetAll(r, int(start), int(sizeI)) {
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

		for _, iter := range s.SandBoxIndexStorage.GetAll(r, int(offsetI), int(limitI)) {
			result = append(result, *iter)
		}
	}

	in := Model.SandBoxIndex{}
	count := s.SandBoxIndexStorage.Count(r, in)

	return uint(count), &Response{Res: result}, nil
}

// FindOne to satisfy `api2go.DataSource` interface
// this method should return the model with the given ID, otherwise an error
func (s SandBoxIndexResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	modelRoot, err := s.SandBoxIndexStorage.GetOne(ID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}

	return &Response{Res: modelRoot}, nil
}

// Create method to satisfy `api2go.DataSource` interface
func (s SandBoxIndexResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(Model.SandBoxIndex)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid_Instance_Given"), "Invalid Instance Given", http.StatusBadRequest)
	}

	id := s.SandBoxIndexStorage.Insert(model)
	model.ID = id

	return &Response{Res: model, Code: http.StatusCreated}, nil
}

// Delete to satisfy `api2go.DataSource` interface
func (s SandBoxIndexResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := s.SandBoxIndexStorage.Delete(id)
	return &Response{Code: http.StatusNoContent}, err
}

//Update stores all changes on the model
func (s SandBoxIndexResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(Model.SandBoxIndex)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid_Instance_Given"), "Invalid Instance Given", http.StatusBadRequest)
	}

	err := s.SandBoxIndexStorage.Update(model)
	return &Response{Res: model, Code: http.StatusNoContent}, err
}
