package http

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"reservation-system/internal/domain/candidate"
	"reservation-system/internal/service/reservation"
	"reservation-system/pkg/server/response"
	"reservation-system/pkg/store"
)

type CandidateHandler struct {
	reservationService *reservation.Service
}

func NewCandidateHandler(s *reservation.Service) *CandidateHandler {
	return &CandidateHandler{reservationService: s}
}

func (h *CandidateHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.list)
	r.Post("/", h.add)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.get)
		r.Put("/", h.update)
		r.Delete("/", h.delete)
	})

	return r
}

// @Summary	list of candidates from the repository
// @Tags		candidates
// @Accept		json
// @Produce	json
// @Success	200			{array}		candidate.Response
// @Failure	500			{object}	response.Object
// @Router		/candidates 	[get]
func (h *CandidateHandler) list(w http.ResponseWriter, r *http.Request) {
	res, err := h.reservationService.ListCandidates(r.Context())
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, res)
}

// @Summary	add a new candidate to the repository
// @Tags		candidates
// @Accept		json
// @Produce	json
// @Param		request	body		candidate.Request	true	"body param"
// @Success	200		{object}	candidate.Response
// @Failure	400		{object}	response.Object
// @Failure	500		{object}	response.Object
// @Router		/candidates [post]
func (h *CandidateHandler) add(w http.ResponseWriter, r *http.Request) {
	req := candidate.Request{}
	if err := render.Bind(r, &req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	res, err := h.reservationService.AddCandidate(r.Context(), req)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, res)
}

// @Summary	get the candidate from the repository
// @Tags		candidates
// @Accept		json
// @Produce	json
// @Param		id	path		int	true	"path param"
// @Success	200	{object}	candidate.Response
// @Failure	404	{object}	response.Object
// @Failure	500	{object}	response.Object
// @Router		/candidates/{id} [get]
func (h *CandidateHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	res, err := h.reservationService.GetCandidate(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(w, r, err)
		default:
			response.InternalServerError(w, r, err)
		}
		return
	}

	response.OK(w, r, res)
}

// @Summary	update the candidate in the repository
// @Tags		candidates
// @Accept		json
// @Produce	json
// @Param		id		path	int				true	"path param"
// @Param		request	body	candidate.Request	true	"body param"
// @Success	200
// @Failure	400	{object}	response.Object
// @Failure	404	{object}	response.Object
// @Failure	500	{object}	response.Object
// @Router		/candidates/{id} [put]
func (h *CandidateHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	req := candidate.Request{}
	if err := render.Bind(r, &req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	if err := h.reservationService.UpdateCandidate(r.Context(), id, req); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(w, r, err)
		default:
			response.InternalServerError(w, r, err)
		}
		return
	}
}

// @Summary	delete the candidate from the repository
// @Tags		candidates
// @Accept		json
// @Produce	json
// @Param		id	path	int	true	"path param"
// @Success	200
// @Failure	404	{object}	response.Object
// @Failure	500	{object}	response.Object
// @Router		/candidates/{id} [delete]
func (h *CandidateHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.reservationService.DeleteCandidate(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(w, r, err)
		default:
			response.InternalServerError(w, r, err)
		}
		return
	}
}
