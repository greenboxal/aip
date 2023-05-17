package rest

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type ResourcesAPI struct {
	chi.Router

	db forddb.Database
}

func NewResourcesAPI(db forddb.Database) *ResourcesAPI {
	api := &ResourcesAPI{
		Router: chi.NewMux(),
		db:     db,
	}

	api.Get("/{resource}", api.ListResource)
	api.Get("/{resource}/{id}", api.GetResource)
	api.Put("/{resource}/{id}", api.CreateOrUpdateResource)
	api.Delete("/{resource}/{id}", api.DeleteResource)

	return api
}

func (a *ResourcesAPI) ListResource(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	resourceTypeName := chi.URLParam(request, "resource")
	resourceType := forddb.LookupTypeByName(resourceTypeName)

	if resourceType == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	resources, err := a.db.List(ctx, resourceType.GetResourceID())

	if err != nil {
		panic(err)
	}

	data, err := json.Marshal(resources)

	if err != nil {
		panic(err)
	}

	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(data)
}

func (a *ResourcesAPI) GetResource(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	resourceTypeName := chi.URLParam(request, "resource")
	resourceIdName := chi.URLParam(request, "id")

	resourceType := forddb.LookupTypeByName(resourceTypeName)

	if resourceType == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	resourceId := resourceType.CreateID(resourceIdName)
	resource, err := a.db.Get(ctx, resourceType.GetResourceID(), resourceId)

	if err != nil {
		panic(err)
	}

	if resource == nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	data, err := json.Marshal(resource)

	if err != nil {
		panic(err)
	}

	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(data)
}

func (a *ResourcesAPI) CreateOrUpdateResource(writer http.ResponseWriter, request *http.Request) {
	resourceTypeName := chi.URLParam(request, "resource")
	resourceIdName := chi.URLParam(request, "id")

	resourceType := forddb.LookupTypeByName(resourceTypeName)

	if resourceType == nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	resourceId := resourceType.CreateID(resourceIdName)
	resource := resourceType.CreateInstance().(forddb.BasicResource)

	data, err := io.ReadAll(request.Body)

	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, resource); err != nil {
		panic(err)
	}

	if resource.GetResourceBasicID() != resourceId {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := forddb.Put(request.Context(), a.db, resource)

	if err == forddb.ErrVersionMismatch {
		writer.WriteHeader(http.StatusConflict)
		return
	} else if err != nil {
		panic(err)
	}

	data, err = json.Marshal(result)

	if err != nil {
		panic(err)
	}

	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(data)
}

func (a *ResourcesAPI) DeleteResource(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	resourceTypeName := chi.URLParam(request, "resource")
	resourceIdName := chi.URLParam(request, "id")

	resourceType := forddb.LookupTypeByName(resourceTypeName)

	if resourceType == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	resourceId := resourceType.CreateID(resourceIdName)
	resource, err := a.db.Get(ctx, resourceType.GetResourceID(), resourceId)

	if err != nil {
		panic(err)
	}

	if resource == nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	resource, err = a.db.Delete(ctx, resource)

	if err != nil {
		panic(err)
	}

	if resource == nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	data, err := json.Marshal(resource)

	if err != nil {
		panic(err)
	}

	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(data)
}
