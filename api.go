//go:generate go tool oapi-codegen -config cfg.yaml https://raw.githubusercontent.com/opsminded/spec/refs/heads/main/openapi.json
package api

import (
	"context"
	"errors"

	"github.com/opsminded/graphlib"
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
			Label:   v.Label,
			Healthy: v.Health,
		}
		summary.UnhealthyVertices = append(summary.UnhealthyVertices, vertex)
	}

	return Summary200JSONResponse(summary), nil
}

func (api *API) GetVertex(ctx context.Context, request GetVertexRequestObject) (GetVertexResponseObject, error) {
	p, err := api.service.GetVertex(request.Label)
	if errors.As(err, &graphlib.VertexNotFoundError{}) {
		nf := NotFoundJSONResponse{Code: 404, Error: err.Error()}
		return GetVertex404JSONResponse{NotFoundJSONResponse: nf}, nil
	}
	if err != nil {
		ise := InternalServerErrorJSONResponse{Code: 500, Error: err.Error()}
		return GetVertex500JSONResponse{InternalServerErrorJSONResponse: ise}, nil
	}
	v := Vertex{
		Label:   p.Label,
		Healthy: p.Health,
	}
	return GetVertex200JSONResponse(v), nil
}

func (api *API) GetVertexDependents(ctx context.Context, request GetVertexDependentsRequestObject) (GetVertexDependentsResponseObject, error) {
	pall := false
	if request.Params.All != nil {
		pall = *request.Params.All
	}

	serviceSub, err := api.service.GetVertexDependencies(request.Label, false)
	if errors.As(err, &graphlib.VertexNotFoundError{}) {
		nf := NotFoundJSONResponse{Code: 404, Error: err.Error()}
		return GetVertexDependents404JSONResponse{NotFoundJSONResponse: nf}, nil
	}

	if err != nil {
		ise := InternalServerErrorJSONResponse{Code: 500, Error: err.Error()}
		return GetVertexDependents500JSONResponse{InternalServerErrorJSONResponse: ise}, nil
	}

	sub := Subgraph{
		Title: "Dependentes de " + request.Label,
		All:   pall,
		Principal: Vertex{
			Label:   serviceSub.Principal.Label,
			Healthy: serviceSub.Principal.Health,
		},
		Edges:      []Edge{},
		Vertices:   []Vertex{},
		Highlights: []Vertex{},
	}

	for _, e := range serviceSub.SubGraph.Edges {
		edge := Edge{
			Label:       e.Label,
			Source:      e.Source.Label,
			Destination: e.Destination.Label,
		}
		sub.Edges = append(sub.Edges, edge)
	}
	for _, v := range serviceSub.SubGraph.Vertices {
		vertex := Vertex{
			Label:   v.Label,
			Healthy: v.Health,
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

	serviceSub, err := api.service.GetVertexDependencies(request.Label, pall)

	if errors.As(err, &graphlib.VertexNotFoundError{}) {
		nf := NotFoundJSONResponse{Code: 404, Error: err.Error()}
		return GetVertexDependencies404JSONResponse{NotFoundJSONResponse: nf}, nil
	}

	if err != nil {
		ise := InternalServerErrorJSONResponse{Code: 500, Error: err.Error()}
		return GetVertexDependencies500JSONResponse{InternalServerErrorJSONResponse: ise}, nil
	}

	sub := Subgraph{
		Title: "Dependencias de " + request.Label,
		All:   serviceSub.All,
		Principal: Vertex{
			Label:   serviceSub.Principal.Label,
			Healthy: serviceSub.Principal.Health,
		},
		Edges:      []Edge{},
		Vertices:   []Vertex{},
		Highlights: []Vertex{},
	}

	for _, e := range serviceSub.SubGraph.Edges {
		edge := Edge{
			Label:       e.Label,
			Source:      e.Source.Label,
			Destination: e.Destination.Label,
		}
		sub.Edges = append(sub.Edges, edge)
	}

	for _, v := range serviceSub.SubGraph.Vertices {
		vertex := Vertex{
			Label:   v.Label,
			Healthy: v.Health,
		}
		sub.Vertices = append(sub.Vertices, vertex)
	}

	return GetVertexDependencies200JSONResponse(sub), nil
}

func (api *API) GetVertexLineages(ctx context.Context, request GetVertexLineagesRequestObject) (GetVertexLineagesResponseObject, error) {
	sub := Subgraph{
		Title:      "Linhas de " + request.Label,
		Principal:  Vertex{},
		Edges:      []Edge{},
		Vertices:   []Vertex{},
		Highlights: []Vertex{},
	}
	return GetVertexLineages200JSONResponse(sub), nil
}

func (api *API) GetVertexNeighbors(ctx context.Context, request GetVertexNeighborsRequestObject) (GetVertexNeighborsResponseObject, error) {
	p, err := api.service.GetVertex(request.Label)
	if errors.As(err, &graphlib.VertexNotFoundError{}) {
		nf := NotFoundJSONResponse{Code: 404, Error: err.Error()}
		return GetVertexNeighbors404JSONResponse{NotFoundJSONResponse: nf}, nil
	}
	if err != nil {
		ise := InternalServerErrorJSONResponse{Code: 500, Error: err.Error()}
		return GetVertexNeighbors500JSONResponse{InternalServerErrorJSONResponse: ise}, nil
	}

	serviceSub, err := api.service.Neighbors(request.Label)
	if errors.As(err, &graphlib.VertexNotFoundError{}) {
		nf := NotFoundJSONResponse{Code: 404, Error: err.Error()}
		return GetVertexNeighbors404JSONResponse{NotFoundJSONResponse: nf}, err
	}

	if err != nil {
		ise := InternalServerErrorJSONResponse{Code: 500, Error: err.Error()}
		return GetVertexNeighbors500JSONResponse{InternalServerErrorJSONResponse: ise}, nil
	}

	ss := Subgraph{
		Principal: Vertex{
			Label:   p.Label,
			Healthy: p.Health,
		},
		Edges:      []Edge{},
		Vertices:   []Vertex{},
		Highlights: []Vertex{},
	}

	for _, e := range serviceSub.SubGraph.Edges {
		edge := Edge{
			Label:       e.Label,
			Source:      e.Source.Label,
			Destination: e.Destination.Label,
		}
		ss.Edges = append(ss.Edges, edge)
	}

	for _, v := range serviceSub.SubGraph.Vertices {
		vertex := Vertex{
			Label:   v.Label,
			Healthy: v.Health,
		}
		ss.Vertices = append(ss.Vertices, vertex)
	}

	return GetVertexNeighbors200JSONResponse(ss), nil
}

func (api *API) GetPath(ctx context.Context, request GetPathRequestObject) (GetPathResponseObject, error) {
	serviceSub, err := api.service.Path(request.Label, request.Destination)

	if errors.As(err, &graphlib.VertexNotFoundError{}) {
		nf := NotFoundJSONResponse{Code: 404, Error: err.Error()}
		return GetPath404JSONResponse{NotFoundJSONResponse: nf}, nil
	}

	if err != nil {
		ise := InternalServerErrorJSONResponse{Code: 500, Error: err.Error()}
		return GetPath500JSONResponse{InternalServerErrorJSONResponse: ise}, nil
	}

	sub := Subgraph{
		Title: "Caminho entre " + request.Label + " e " + request.Destination,
		Principal: Vertex{
			Label:   serviceSub.Principal.Label,
			Healthy: serviceSub.Principal.Health,
		},
		Edges:      []Edge{},
		Vertices:   []Vertex{},
		Highlights: []Vertex{},
	}

	for _, e := range serviceSub.SubGraph.Edges {
		edge := Edge{
			Label:       e.Label,
			Source:      e.Source.Label,
			Destination: e.Destination.Label,
		}
		sub.Edges = append(sub.Edges, edge)
	}

	for _, v := range serviceSub.SubGraph.Vertices {
		vertex := Vertex{
			Label:   v.Label,
			Healthy: v.Health,
		}
		sub.Vertices = append(sub.Vertices, vertex)
	}

	return GetPath200JSONResponse(sub), nil
}
