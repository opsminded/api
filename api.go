//go:generate go tool oapi-codegen -config cfg.yaml https://raw.githubusercontent.com/opsminded/spec/refs/heads/main/openapi.json
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
	summary := Summary{
		TotalVertices: sum.TotalVertex,
	}
	return Summary200JSONResponse(summary), nil
}

func (api *API) GetVertex(ctx context.Context, request GetVertexRequestObject) (GetVertexResponseObject, error) {
	p, err := api.service.GetVertex(request.Label)
	if err != nil {
		return nil, err
	}

	v := Vertex{
		Label: p.Label,
	}

	return GetVertex200JSONResponse(v), nil
}

func (api *API) GetVertexDependants(ctx context.Context, request GetVertexDependantsRequestObject) (GetVertexDependantsResponseObject, error) {
	pall := false
	if request.Params.All != nil {
		pall = *request.Params.All
	}

	serviceSub := api.service.GetVertexDependencies(request.Label, false)

	sub := Subgraph{
		Title: "Dependentes de " + request.Label,
		All:   pall,
		Principal: Vertex{
			Label: serviceSub.Principal.Label,
		},
		Edges:    []Edge{},
		Vertices: []Vertex{},
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
			Label: v.Label,
		}
		sub.Vertices = append(sub.Vertices, vertex)
	}

	return GetVertexDependants200JSONResponse(sub), nil
}

func (api *API) GetVertexDependencies(ctx context.Context, request GetVertexDependenciesRequestObject) (GetVertexDependenciesResponseObject, error) {
	pall := false
	if request.Params.All != nil {
		pall = *request.Params.All
	}

	serviceSub := api.service.GetVertexDependencies(request.Label, pall)

	sub := Subgraph{
		Title: "Dependencias de " + request.Label,
		All:   serviceSub.All,
		Principal: Vertex{
			Label: serviceSub.Principal.Label,
		},
		Edges:    []Edge{},
		Vertices: []Vertex{},
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
			Label: v.Label,
		}
		sub.Vertices = append(sub.Vertices, vertex)
	}

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

	serviceSub := api.service.Neighbors(request.Label)

	ss := Subgraph{
		Principal: Vertex{
			Label: p.Label,
		},
		Edges:    []Edge{},
		Vertices: []Vertex{},
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
			Label: v.Label,
		}
		ss.Vertices = append(ss.Vertices, vertex)
	}

	return GetVertexNeighbors200JSONResponse(ss), nil
}

func (api *API) GetPath(ctx context.Context, request GetPathRequestObject) (GetPathResponseObject, error) {
	serviceSub := api.service.Path(request.Label, request.Destination)
	sub := Subgraph{
		Title: "Caminho entre " + request.Label + " e " + request.Destination,
		Principal: Vertex{
			Label: serviceSub.Principal.Label,
		},
		Edges:    []Edge{},
		Vertices: []Vertex{},
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
			Label: v.Label,
		}
		sub.Vertices = append(sub.Vertices, vertex)
	}

	return GetPath200JSONResponse(sub), nil
}
