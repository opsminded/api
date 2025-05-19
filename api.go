//go:generate go tool oapi-codegen -config cfg.yaml https://raw.githubusercontent.com/opsminded/spec/refs/heads/main/openapi.json
package api

import (
	"context"
	"errors"
	"log"

	"github.com/opsminded/graphlib/v2"
	"github.com/opsminded/service"
)

type API struct {
	service *service.Service
}

var _ StrictServerInterface = (*API)(nil)

func New(s *service.Service) StrictServerInterface {
	return &API{
		service: s,
	}
}

func (api *API) Summary(ctx context.Context, request SummaryRequestObject) (SummaryResponseObject, error) {
	sum := api.service.Summary()
	summary := Summary{
		TotalEdges:        sum.TotalEdges,
		TotalVertices:     sum.TotalVertices,
		UnhealthyVertices: []Vertex{},
	}

	for _, v := range sum.UnhealthyVertices {
		vertex := Vertex{
			Key:     v.Key,
			Label:   v.Label,
			Healthy: v.Healthy,
		}
		summary.UnhealthyVertices = append(summary.UnhealthyVertices, vertex)
	}

	return Summary200JSONResponse(summary), nil
}

func (api *API) GetVertex(ctx context.Context, request GetVertexRequestObject) (GetVertexResponseObject, error) {
	p, err := api.service.GetVertex(request.Key)

	if errors.As(err, &graphlib.VertexNotFoundErr{}) {
		nf := NotFoundJSONResponse{Code: 404, Error: err.Error()}
		return GetVertex404JSONResponse{NotFoundJSONResponse: nf}, nil
	}
	if err != nil {
		ise := InternalServerErrorJSONResponse{Code: 500, Error: err.Error()}
		return GetVertex500JSONResponse{InternalServerErrorJSONResponse: ise}, nil
	}
	v := Vertex{
		Key:     p.Key,
		Label:   p.Label,
		Healthy: p.Healthy,
	}
	return GetVertex200JSONResponse(v), nil
}

func (api *API) GetVertexAttributes(ctx context.Context, request GetVertexAttributesRequestObject) (GetVertexAttributesResponseObject, error) {
	// TODO: Implement this function

	x := VertexAttrubutes_Value{}
	x.FromVertexAttrubutesValue0("hello")

	y := VertexAttrubutes_Value{}
	y.FromVertexAttrubutesValue1(234234)

	return GetVertexAttributes200JSONResponse{
		struct {
			Description string                 `json:"description"`
			Type        string                 `json:"type"`
			Value       VertexAttrubutes_Value `json:"value"`
		}{
			Description: "description",
			Type:        "string",
			Value:       x,
		},
		struct {
			Description string                 `json:"description"`
			Type        string                 `json:"type"`
			Value       VertexAttrubutes_Value `json:"value"`
		}{
			Description: "strange number",
			Type:        "integer",
			Value:       y,
		},
	}, nil
}

func (api *API) GetVertexDependents(ctx context.Context, request GetVertexDependentsRequestObject) (GetVertexDependentsResponseObject, error) {
	pall := false
	if request.Params.All != nil {
		pall = *request.Params.All
	}

	serviceSub, err := api.service.VertexDependents(request.Key, pall)
	if errors.As(err, &graphlib.VertexNotFoundErr{}) {
		nf := NotFoundJSONResponse{Code: 404, Error: err.Error()}
		return GetVertexDependents404JSONResponse{NotFoundJSONResponse: nf}, nil
	}

	if err != nil {
		ise := InternalServerErrorJSONResponse{Code: 500, Error: err.Error()}
		return GetVertexDependents500JSONResponse{InternalServerErrorJSONResponse: ise}, nil
	}

	sub := Subgraph{
		Title: "Dependentes de " + request.Key,
		All:   pall,
		Principal: Vertex{
			Key:     serviceSub.Principal.Key,
			Label:   serviceSub.Principal.Label,
			Healthy: serviceSub.Principal.Healthy,
		},
		Edges:      []Edge{},
		Vertices:   []Vertex{},
		Highlights: []Vertex{},
	}

	for _, e := range serviceSub.SubGraph.Edges {
		edge := Edge{
			Key:    e.Key,
			Source: e.Source,
			Target: e.Target,
		}
		sub.Edges = append(sub.Edges, edge)
	}
	for _, v := range serviceSub.SubGraph.Vertices {
		vertex := Vertex{
			Key:     v.Key,
			Label:   v.Label,
			Healthy: v.Healthy,
		}
		sub.Vertices = append(sub.Vertices, vertex)
	}
	return GetVertexDependents200JSONResponse(sub), nil
}

func (api *API) GetVertexDependencies(ctx context.Context, request GetVertexDependenciesRequestObject) (GetVertexDependenciesResponseObject, error) {
	pall := false
	if request.Params.All != nil {
		pall = *request.Params.All
	}

	log.Println("api GetVertexDependencies", request.Key, pall)
	serviceSub, err := api.service.VertexDependencies(request.Key, pall)

	if errors.As(err, &graphlib.VertexNotFoundErr{}) {
		nf := NotFoundJSONResponse{Code: 404, Error: err.Error()}
		return GetVertexDependencies404JSONResponse{NotFoundJSONResponse: nf}, nil
	}

	if err != nil {
		ise := InternalServerErrorJSONResponse{Code: 500, Error: err.Error()}
		return GetVertexDependencies500JSONResponse{InternalServerErrorJSONResponse: ise}, nil
	}

	sub := Subgraph{
		Title: "Dependencias de " + request.Key,
		All:   pall,
		Principal: Vertex{
			Key:     serviceSub.Principal.Key,
			Label:   serviceSub.Principal.Label,
			Healthy: serviceSub.Principal.Healthy,
		},
		Edges:      []Edge{},
		Vertices:   []Vertex{},
		Highlights: []Vertex{},
	}

	for _, e := range serviceSub.SubGraph.Edges {
		edge := Edge{
			Key:    e.Key,
			Source: e.Source,
			Target: e.Target,
		}
		sub.Edges = append(sub.Edges, edge)
	}

	for _, v := range serviceSub.SubGraph.Vertices {
		vertex := Vertex{
			Key:     v.Key,
			Label:   v.Label,
			Healthy: v.Healthy,
		}
		sub.Vertices = append(sub.Vertices, vertex)
	}

	return GetVertexDependencies200JSONResponse(sub), nil
}

func (api *API) GetVertexNeighbors(ctx context.Context, request GetVertexNeighborsRequestObject) (GetVertexNeighborsResponseObject, error) {
	p, err := api.service.GetVertex(request.Key)
	if errors.As(err, &graphlib.VertexNotFoundErr{}) {
		nf := NotFoundJSONResponse{Code: 404, Error: err.Error()}
		return GetVertexNeighbors404JSONResponse{NotFoundJSONResponse: nf}, nil
	}
	if err != nil {
		ise := InternalServerErrorJSONResponse{Code: 500, Error: err.Error()}
		return GetVertexNeighbors500JSONResponse{InternalServerErrorJSONResponse: ise}, nil
	}

	serviceSub, err := api.service.VertexNeighbors(request.Key)
	if errors.As(err, &graphlib.VertexNotFoundErr{}) {
		nf := NotFoundJSONResponse{Code: 404, Error: err.Error()}
		return GetVertexNeighbors404JSONResponse{NotFoundJSONResponse: nf}, err
	}

	if err != nil {
		ise := InternalServerErrorJSONResponse{Code: 500, Error: err.Error()}
		return GetVertexNeighbors500JSONResponse{InternalServerErrorJSONResponse: ise}, nil
	}

	ss := Subgraph{
		Principal: Vertex{
			Key:     p.Key,
			Label:   p.Label,
			Healthy: p.Healthy,
		},
		Edges:      []Edge{},
		Vertices:   []Vertex{},
		Highlights: []Vertex{},
	}

	for _, e := range serviceSub.SubGraph.Edges {
		edge := Edge{
			Key:    e.Key,
			Source: e.Source,
			Target: e.Target,
		}
		ss.Edges = append(ss.Edges, edge)
	}

	for _, v := range serviceSub.SubGraph.Vertices {
		vertex := Vertex{
			Key:     v.Key,
			Label:   v.Label,
			Healthy: v.Healthy,
		}
		ss.Vertices = append(ss.Vertices, vertex)
	}

	return GetVertexNeighbors200JSONResponse(ss), nil
}

func (api *API) GetPath(ctx context.Context, request GetPathRequestObject) (GetPathResponseObject, error) {
	serviceSub, err := api.service.Path(request.Key, request.Target)

	if errors.As(err, &graphlib.VertexNotFoundErr{}) {
		nf := NotFoundJSONResponse{Code: 404, Error: err.Error()}
		return GetPath404JSONResponse{NotFoundJSONResponse: nf}, nil
	}

	if err != nil {
		ise := InternalServerErrorJSONResponse{Code: 500, Error: err.Error()}
		return GetPath500JSONResponse{InternalServerErrorJSONResponse: ise}, nil
	}

	sub := Subgraph{
		Title: "Caminho entre " + serviceSub.Principal.Label + " e " + request.Target,
		Principal: Vertex{
			Key:     serviceSub.Principal.Key,
			Label:   serviceSub.Principal.Label,
			Healthy: serviceSub.Principal.Healthy,
		},
		Edges:      []Edge{},
		Vertices:   []Vertex{},
		Highlights: []Vertex{},
	}

	for _, e := range serviceSub.SubGraph.Edges {
		edge := Edge{
			Key:    e.Key,
			Source: e.Source,
			Target: e.Target,
		}
		sub.Edges = append(sub.Edges, edge)
	}

	for _, v := range serviceSub.SubGraph.Vertices {
		vertex := Vertex{
			Key:     v.Key,
			Label:   v.Label,
			Healthy: v.Healthy,
		}
		sub.Vertices = append(sub.Vertices, vertex)
	}

	return GetPath200JSONResponse(sub), nil
}

func (api *API) ClearHealthStatus(ctx context.Context, request ClearHealthStatusRequestObject) (ClearHealthStatusResponseObject, error) {
	api.service.ClearGraphHealthyStatus()
	return ClearHealthStatus200Response{}, nil
}

func (api *API) MarkVertexHealthy(ctx context.Context, request MarkVertexHealthyRequestObject) (MarkVertexHealthyResponseObject, error) {
	err := api.service.SetVertexHealth(request.Key, true)
	if errors.As(err, &graphlib.VertexNotFoundErr{}) {
		nf := NotFoundJSONResponse{Code: 404, Error: err.Error()}
		return MarkVertexHealthy404JSONResponse{NotFoundJSONResponse: nf}, nil
	}
	if err != nil {
		ise := InternalServerErrorJSONResponse{Code: 500, Error: err.Error()}
		return MarkVertexHealthy500JSONResponse{InternalServerErrorJSONResponse: ise}, nil
	}
	return MarkVertexHealthy200Response{}, nil
}

func (api *API) MarkVertexUnhealthy(ctx context.Context, request MarkVertexUnhealthyRequestObject) (MarkVertexUnhealthyResponseObject, error) {
	err := api.service.SetVertexHealth(request.Key, false)

	if errors.As(err, &graphlib.VertexNotFoundErr{}) {
		nf := NotFoundJSONResponse{Code: 404, Error: err.Error()}
		return MarkVertexUnhealthy404JSONResponse{NotFoundJSONResponse: nf}, nil
	}
	if err != nil {
		ise := InternalServerErrorJSONResponse{Code: 500, Error: err.Error()}
		return MarkVertexUnhealthy500JSONResponse{InternalServerErrorJSONResponse: ise}, nil
	}

	return MarkVertexUnhealthy200Response{}, nil
}
