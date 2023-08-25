package kong

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// AbstractFileResourceService handles workspace FileResources in Kong.
type AbstractFileResourceService interface {
	// Create creates a workspace FileService in Kong.
	Create(ctx context.Context, FileResource *FileResource) (*FileResource, error)
	// Get fetches a workspace FileResource in Kong.
	Get(ctx context.Context, ID *string) (*FileResource, error)
	// Path fetches a workspace FileResource in Kong.
	GetByPath(ctx context.Context, Path *string) (*FileResource, error)
	// Update updates a workspace FileResource in Kong
	Update(ctx context.Context, FileResource *FileResource) (*FileResource, error)
	// Delete deletes a workspace FileResource in Kong
	Delete(ctx context.Context, ID *string) error
	// List fetches a list of workspace FileResources in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*FileResource, *ListOpt, error)
	// ListAll fetches all workspace FileResources in Kong.
	ListAll(ctx context.Context) ([]*FileResource, error)
}

// FileResourceService handles FileResources in Kong.
type FileResourceService service

func (s *FileResourceService) Create(ctx context.Context,
	fileResource *FileResource,
) (*FileResource, error) {
	queryPath := "/files"
	method := "POST"
	req, err := s.client.NewRequest(method, queryPath, nil, fileResource)
	if err != nil {
		return nil, err
	}

	createdFileResource := FileResource{}
	_, err = s.client.Do(ctx, req, &createdFileResource)
	if err != nil {
		return nil, err
	}
	return &createdFileResource, nil
}

// Get fetches a FileResource in Kong.
func (s *FileResourceService) Get(ctx context.Context,
	ID *string,
) (*FileResource, error) {
	if isEmptyString(ID) {
		return nil, fmt.Errorf("ID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/files/%v", *ID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var fileResource FileResource
	_, err = s.client.Do(ctx, req, &fileResource)
	if err != nil {
		return nil, err
	}
	return &fileResource, nil
}

// GetByPath fetches a workspace FileResource in Kong using it's path (must end with a valid file extension).
func (s *FileResourceService) GetByPath(ctx context.Context,
	path *string,
) (*FileResource, error) {
	if isEmptyString(path) {
		return nil, fmt.Errorf("path cannot be nil for Get operation")
	}

	type QS struct {
		path string `url:"path,omitempty"`
	}

	req, err := s.client.NewRequest("GET", "/files",
		&QS{path: *path}, nil)
	if err != nil {
		return nil, err
	}

	type Response struct {
		Data []FileResource
	}
	var resp Response
	_, err = s.client.Do(ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, NewAPIError(http.StatusNotFound, "Not found")
	}

	return &resp.Data[0], nil
}

// Update updates a workspace FileResource in Kong
func (s *FileResourceService) Update(ctx context.Context,
	fileResource *FileResource,
) (*FileResource, error) {
	if isEmptyString(fileResource.ID) {
		return nil, fmt.Errorf("ID cannot be nil for Update operation")
	}

	endpoint := fmt.Sprintf("/files/%v", *fileResource.ID)
	req, err := s.client.NewRequest("PATCH", endpoint, nil, fileResource)
	if err != nil {
		return nil, err
	}
	type Response struct {
		fileResource FileResource
	}
	var resp Response
	_, err = s.client.Do(ctx, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.fileResource, nil
}

// Delete deletes a FileResource in Kong
func (s *FileResourceService) Delete(ctx context.Context,
	ID *string,
) error {
	if isEmptyString(ID) {
		return fmt.Errorf("ID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/files/%v", *ID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List fetches a list of workspace FileResources in Kong.
// opt can be used to control pagination.
func (s *FileResourceService) List(ctx context.Context,
	opt *ListOpt,
) ([]*FileResource, *ListOpt, error) {
	data, next, err := s.client.list(ctx, "/files", opt)
	if err != nil {
		return nil, nil, err
	}
	var fileResources []*FileResource

	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, nil, err
		}
		var fileResource FileResource
		err = json.Unmarshal(b, &fileResource)
		if err != nil {
			return nil, nil, err
		}
		fileResources = append(fileResources, &fileResource)
	}

	return fileResources, next, nil
}

// ListAll fetches all workspace FileResources in Kong.
// This method can take a while if there
// a lot of FileResources present.
func (s *FileResourceService) ListAll(ctx context.Context) ([]*FileResource, error) {
	var fileResources, data []*FileResource
	var err error
	opt := &ListOpt{Size: pageSize}

	for opt != nil {
		data, opt, err = s.List(ctx, opt)
		if err != nil {
			return nil, err
		}
		fileResources = append(fileResources, data...)
	}
	return fileResources, nil
}
