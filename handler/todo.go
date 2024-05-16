package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// ServeHTTP implements http.Handler interface.
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		prev_id, err := strconv.Atoi(r.URL.Query().Get("prev_id"))
		if err != nil {
			prev_id = 0
		}
		size, err := strconv.Atoi(r.URL.Query().Get("size"))
		if err != nil {
			size = 0
		}
		req := &model.ReadTODORequest{PrevID: int64(prev_id), Size: int64(size)}

		res, err := h.Read(r.Context(), req)
		if err != nil {
			panic(err)
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}
	case "POST":
		req := &model.CreateTODORequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			panic(err)
		}
		if req.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Bad request")
			return
		}

		res, err := h.Create(r.Context(), req)
		if err != nil {
			panic(err)
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}

	case "PUT":
		req := &model.UpdateTODORequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			panic(err)
		}
		if req.ID == 0 || req.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Bad request")
			return
		}

		res, err := h.Update(r.Context(), req)
		if err != nil {
			panic(err)
		}

		if err = json.NewEncoder(w).Encode(res); err != nil {
			panic(err)
		}
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	t, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}

	return &model.CreateTODOResponse{TODO: t}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	tt, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		return nil, err
	}
	return &model.ReadTODOResponse{TODOs: tt}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	t, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}

	return &model.UpdateTODOResponse{TODO: t}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}
