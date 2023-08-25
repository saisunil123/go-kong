package kong

import (
	"context"
	"encoding/json"
	"fmt"
)

type RBACEndpointPermissionResponse struct {
	CreatedAt *int      `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	Workspace *string   `json:"workspace,omitempty" yaml:"workspace,omitempty"`
	Endpoint  *string   `json:"endpoint,omitempty" yaml:"endpoint,omitempty"`
	Actions   []*string `json:"actions,omitempty" yaml:"actions,omitempty"`
	Negative  *bool     `json:"negative,omitempty" yaml:"negative,omitempty"`
	Role      *RBACRole `json:"role,omitempty" yaml:"role,omitempty"`
	Comment   *string   `json:"comment,omitempty" yaml:"comment,omitempty"`
}

// MarshalJSON marshals an endpoint permission into a suitable form for the Kong admin API
func (e *Response) MarshalJSON() ([]byte, error) {
	type ep struct {
		CreatedAt *int      `json:"created_at,omitempty" yaml:"created_at,omitempty"`
		Workspace *string   `json:"workspace,omitempty" yaml:"workspace,omitempty"`
		Endpoint  *string   `json:"endpoint,omitempty" yaml:"endpoint,omitempty"`
		Actions   []*string `json:"actions,omitempty" yaml:"actions,omitempty"`
		Negative  *bool     `json:"negative,omitempty" yaml:"negative,omitempty"`
		Role      *RBACRole `json:"role,omitempty" yaml:"role,omitempty"`
		Comment   *string   `json:"comment,omitempty" yaml:"comment,omitempty"`
	}

	return json.Marshal(&ep{
		CreatedAt: e.CreatedAt,
		Workspace: e.Workspace,
		Endpoint:  e.Endpoint,
		Actions:   e.Actions,
		Negative:  e.Negative,
		Comment:   e.Comment,
	})
}

// AbstractRBACEndpointPermissionService handles RBACEndpointPermissions in Kong.
type AbstractRBACEndpointPermissionService interface {
	// Create creates a RBACEndpointPermission in Kong.
	Create(ctx context.Context, ep *RBACEndpointPermission) (*RBACEndpointPermissionResponse, error)
	// Get fetches a RBACEndpointPermission in Kong.
	Get(ctx context.Context, roleNameOrID *string, workspaceNameOrID *string,
		endpointName *string) (*RBACEndpointPermissionResponse, error)
	// Update updates a RBACEndpointPermission in Kong.
	Update(ctx context.Context, ep *RBACEndpointPermission) (*RBACEndpointPermissionResponse, error)
	// Delete deletes a EndpointPermission in Kong
	Delete(ctx context.Context, roleNameOrID *string, workspaceNameOrID *string, endpoint *string) error
	// ListAllForRole fetches a list of all RBACEndpointPermissions in Kong for a given role.
	ListAllForRole(ctx context.Context, roleNameOrID *string) ([]*RBACEndpointPermissionResponse, error)
}

// RBACEndpointPermissionService handles RBACEndpointPermissions in Kong.
type RBACEndpointPermissionService service

// Create creates a RBACEndpointPermission in Kong.
func (s *RBACEndpointPermissionService) Create(ctx context.Context,
	ep *RBACEndpointPermission,
) (*RBACEndpointPermissionResponse, error) {
	if ep == nil {
		return nil, fmt.Errorf("cannot create a nil endpointpermission")
	}
	if ep.Role == nil || ep.Role.ID == nil {
		return nil, fmt.Errorf("cannot create endpoint permission with role or role id undefined")
	}

	method := "POST"
	endpoint := fmt.Sprintf("/rbac/roles/%v/endpoints", *ep.Role.ID)
	req, err := s.client.NewRequest(method, endpoint, nil, ep)
	if err != nil {
		return nil, err
	}

	var createdEndpointPermissionResponse RBACEndpointPermissionResponse

	_, err = s.client.Do(ctx, req, &createdEndpointPermissionResponse)
	if err != nil {
		return nil, err
	}
	return &createdEndpointPermissionResponse, nil
}

// Get fetches a RBACEndpointPermission in Kong.
func (s *RBACEndpointPermissionService) Get(ctx context.Context,
	roleNameOrID *string, workspaceNameOrID *string, endpointName *string,
) (*RBACEndpointPermissionResponse, error) {
	if isEmptyString(endpointName) {
		return nil, fmt.Errorf("endpointName cannot be nil for Get operation")
	}
	if *endpointName == "*" {
		endpointName = String("/" + *endpointName)
	}
	endpoint := fmt.Sprintf("/rbac/roles/%v/endpoints/%v%v", *roleNameOrID, *workspaceNameOrID, *endpointName)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var EndpointPermissionResponse RBACEndpointPermissionResponse
	_, err = s.client.Do(ctx, req, &EndpointPermissionResponse)
	if err != nil {
		return nil, err
	}
	return &EndpointPermissionResponse, nil
}

// Update updates a RBACEndpointPermission in Kong.
func (s *RBACEndpointPermissionService) Update(ctx context.Context,
	ep *RBACEndpointPermission,
) (*RBACEndpointPermissionResponse, error) {
	if ep == nil {
		return nil, fmt.Errorf("cannot update a nil EndpointPermission")
	}
	if ep.Workspace == nil {
		return nil, fmt.Errorf("cannot update an EndpointPermission with workspace as nil")
	}
	if ep.Role == nil || ep.Role.ID == nil {
		return nil, fmt.Errorf("cannot create endpoint permission with role or role id undefined")
	}

	if isEmptyString(ep.Endpoint) {
		return nil, fmt.Errorf("ID cannot be nil for Update operation")
	}

	endpointName := ep.Endpoint
	if *endpointName == "*" {
		endpointName = String("/" + *endpointName)
	}
	endpoint := fmt.Sprintf("/rbac/roles/%v/endpoints/%v%v",
		*ep.Role.ID, *ep.Workspace, *endpointName)
	req, err := s.client.NewRequest("PATCH", endpoint, nil, ep)
	if err != nil {
		return nil, err
	}

	var updatedEndpointPermissionResponse RBACEndpointPermissionResponse
	_, err = s.client.Do(ctx, req, &updatedEndpointPermissionResponse)
	if err != nil {
		return nil, err
	}
	return &updatedEndpointPermissionResponse, nil
}

// Delete deletes a EndpointPermission in Kong
func (s *RBACEndpointPermissionService) Delete(ctx context.Context,
	roleNameOrID *string, workspaceNameOrID *string, endpointName *string,
) error {
	if endpointName == nil {
		return fmt.Errorf("cannot update a nil EndpointPermission")
	}
	if workspaceNameOrID == nil {
		return fmt.Errorf("cannot update an EndpointPermission with workspace as nil")
	}
	if roleNameOrID == nil {
		return fmt.Errorf("cannot update an EndpointPermission with role as nil")
	}

	if *endpointName == "*" {
		endpointName = String("/" + *endpointName)
	}
	reqEndpoint := fmt.Sprintf("/rbac/roles/%v/endpoints/%v/%v",
		*roleNameOrID, *workspaceNameOrID, *endpointName)
	req, err := s.client.NewRequest("DELETE", reqEndpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// ListAllForRole fetches a list of all RBACEndpointPermissions in Kong for a given role.
func (s *RBACEndpointPermissionService) ListAllForRole(ctx context.Context,
	roleNameOrID *string,
) ([]*RBACEndpointPermissionResponse, error) {
	data, _, err := s.client.list(ctx, fmt.Sprintf("/rbac/roles/%v/endpoints", *roleNameOrID), nil)
	if err != nil {
		return nil, err
	}
	var eps []*RBACEndpointPermissionResponse
	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, err
		}
		var ep RBACEndpointPermissionResponse
		err = json.Unmarshal(b, &ep)
		if err != nil {
			return nil, err
		}
		eps = append(eps, &ep)
	}

	return eps, nil
}
