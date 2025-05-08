//go:generate go tool oapi-codegen -config cfg.yaml http://0.0.0.0:3030/openapi.json
package api

import (
	"context"

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

	list := []Vertex{}
	for _, v := range sum.UnhealthVertex {
		vertex := Vertex{
			Label: &v.Label,
			Class: &v.Class,
		}
		list = append(list, vertex)
	}

	summary := Summary{
		TotalVertices:     &sum.TotalVertex,
		UnhealthyVertices: &list,
	}

	return Summary200JSONResponse(summary), nil
}

func (api *API) GetVertex(ctx context.Context, request GetVertexRequestObject) (GetVertexResponseObject, error) {

	p, err := api.service.GetVertex(request.Label)
	if err != nil {
		return nil, err
	}

	v := Vertex{
		Label: &p.Label,
		Class: &p.Class,
	}

	return GetVertex200JSONResponse(v), nil
}

func (api *API) GetVertexDependants(ctx context.Context, request GetVertexDependantsRequestObject) (GetVertexDependantsResponseObject, error) {
	pall := false
	if request.Params.All != nil {
		pall = *request.Params.All
	}
	serviceSub := api.service.GetVertexDependants(request.Label, pall)

	title := "Dependentes de " + request.Label
	edges := []Edge{}
	vertices := []Vertex{}

	sub := Subgraph{
		Title: &title,
		All:   &serviceSub.All,
		Principal: &Vertex{
			Label: &serviceSub.Principal.Label,
			Class: &serviceSub.Principal.Class,
		},
		Edges:    &[]Edge{},
		Vertices: &[]Vertex{},
	}

	for _, e := range serviceSub.Edges {
		edge := Edge{
			Label:       &e.Label,
			Class:       &e.Class,
			Source:      &e.Source,
			Destination: &e.Destination,
		}
		edges = append(edges, edge)
	}

	for _, v := range serviceSub.Vertices {
		vertex := Vertex{
			Label: &v.Label,
			Class: &v.Class,
		}
		vertices = append(vertices, vertex)
	}
	sub.Edges = &edges
	sub.Vertices = &vertices

	return GetVertexDependants200JSONResponse(sub), nil
}

func (api *API) GetVertexDependencies(ctx context.Context, request GetVertexDependenciesRequestObject) (GetVertexDependenciesResponseObject, error) {
	pall := false
	if request.Params.All != nil {
		pall = *request.Params.All
	}

	serviceSub := api.service.GetVertexDependencies(request.Label, pall)

	title := "Dependencias de " + request.Label
	edges := []Edge{}
	vertices := []Vertex{}

	sub := Subgraph{
		Title: &title,
		All:   &serviceSub.All,
		Principal: &Vertex{
			Label: &serviceSub.Principal.Label,
			Class: &serviceSub.Principal.Class,
		},
		Edges:    &[]Edge{},
		Vertices: &[]Vertex{},
	}

	for _, e := range serviceSub.Edges {
		edge := Edge{
			Label:       &e.Label,
			Class:       &e.Class,
			Source:      &e.Source,
			Destination: &e.Destination,
		}
		edges = append(edges, edge)
	}

	for _, v := range serviceSub.Vertices {
		vertex := Vertex{
			Label: &v.Label,
			Class: &v.Class,
		}
		vertices = append(vertices, vertex)
	}

	sub.Edges = &edges
	sub.Vertices = &vertices

	return GetVertexDependencies200JSONResponse(sub), nil
}

func (api *API) GetVertexLineages(ctx context.Context, request GetVertexLineagesRequestObject) (GetVertexLineagesResponseObject, error) {
	return GetVertexLineages200JSONResponse{}, nil
}

func (api *API) GetVertexNeighbors(ctx context.Context, request GetVertexNeighborsRequestObject) (GetVertexNeighborsResponseObject, error) {

	p, err := api.service.GetVertex(request.Label)
	if err != nil {
		return nil, err
	}

	principal := Vertex{
		Label: &p.Label,
		Class: &p.Class,
	}

	edges := []Edge{}
	vertices := []Vertex{}

	serviceSub := api.service.Neighbors(request.Label)

	for _, e := range serviceSub.Edges {
		edge := Edge{
			Label:       &e.Label,
			Class:       &e.Class,
			Source:      &e.Source,
			Destination: &e.Destination,
		}

		edges = append(edges, edge)
	}

	for _, v := range serviceSub.Vertices {
		vertex := Vertex{
			Label: &v.Label,
			Class: &v.Class,
		}

		vertices = append(vertices, vertex)
	}

	ss := Subgraph{
		Principal: &principal,
		Edges:     &edges,
		Vertices:  &vertices,
	}

	return GetVertexNeighbors200JSONResponse(ss), nil
}

func (api *API) GetPath(ctx context.Context, request GetPathRequestObject) (GetPathResponseObject, error) {

	serviceSub := api.service.Path(request.Label, request.Destination)

	title := "Caminho entre " + request.Label + " e " + request.Destination
	edges := []Edge{}
	vertices := []Vertex{}

	sub := Subgraph{
		Title: &title,
		All:   &serviceSub.All,
		Principal: &Vertex{
			Label: &serviceSub.Principal.Label,
			Class: &serviceSub.Principal.Class,
		},
		Edges:    &[]Edge{},
		Vertices: &[]Vertex{},
	}

	for _, e := range serviceSub.Edges {
		edge := Edge{
			Label:       &e.Label,
			Class:       &e.Class,
			Source:      &e.Source,
			Destination: &e.Destination,
		}
		edges = append(edges, edge)
	}

	for _, v := range serviceSub.Vertices {
		vertex := Vertex{
			Label: &v.Label,
			Class: &v.Class,
		}
		vertices = append(vertices, vertex)
	}

	sub.Edges = &edges
	sub.Vertices = &vertices

	return GetPath200JSONResponse(sub), nil
}
