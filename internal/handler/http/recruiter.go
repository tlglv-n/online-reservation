package http

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"reservation-system/internal/domain/recruiter"
	"reservation-system/internal/service/reservation"
	"reservation-system/pkg/server/response"
	"reservation-system/pkg/store"
)

type RecruiterHandler struct {
	reservationService *reservation.Service
}

func NewRecruiterHandler(s *reservation.Service) *RecruiterHandler {
	return &RecruiterHandler{reservationService: s}
}

func (h *RecruiterHandler) Routes() chi.Router {
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

// @Summary	list of recruiters from the repository
// @Tags		recruiters
// @Accept		json
// @Produce	json
// @Success	200			{array}		recruiter.Response
// @Failure	500			{object}	response.Object
// @Router		/recruiters 	[get]
func (h *RecruiterHandler) list(w http.ResponseWriter, r *http.Request) {
	res, err := h.reservationService.ListRecruiters(r.Context())
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, res)
}

// @Summary	add a new recruiter to the repository
// @Tags		recruiters
// @Accept		json
// @Produce	json
// @Param		request	body		recruiter.Request	true	"body param"
// @Success	200		{object}	recruiter.Response
// @Failure	400		{object}	response.Object
// @Failure	500		{object}	response.Object
// @Router		/recruiters [post]
func (h *RecruiterHandler) add(w http.ResponseWriter, r *http.Request) {
	req := recruiter.Request{}
	if err := render.Bind(r, &req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	res, err := h.reservationService.AddRecruiter(r.Context(), req)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, res)
}

// @Summary	get the recruiter from the repository
// @Tags		recruiters
// @Accept		json
// @Produce	json
// @Param		id	path		int	true	"path param"
// @Success	200	{object}	recruiter.Response
// @Failure	404	{object}	response.Object
// @Failure	500	{object}	response.Object
// @Router		/recruiters/{id} [get]
func (h *RecruiterHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	res, err := h.reservationService.GetRecruiter(r.Context(), id)
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

// @Summary	update the recruiter in the repository
// @Tags		recruiters
// @Accept		json
// @Produce	json
// @Param		id		path	int				true	"path param"
// @Param		request	body	recruiter.Request	true	"body param"
// @Success	200
// @Failure	400	{object}	response.Object
// @Failure	404	{object}	response.Object
// @Failure	500	{object}	response.Object
// @Router		/recruiters/{id} [put]
func (h *RecruiterHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	req := recruiter.Request{}
	if err := render.Bind(r, &req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	if err := h.reservationService.UpdateRecruiter(r.Context(), id, req); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(w, r, err)
		default:
			response.InternalServerError(w, r, err)
		}
		return
	}
}

// @Summary	delete the recruiter from the repository
// @Tags		recruiters
// @Accept		json
// @Produce	json
// @Param		id	path	int	true	"path param"
// @Success	200
// @Failure	404	{object}	response.Object
// @Failure	500	{object}	response.Object
// @Router		/recruiters/{id} [delete]
func (h *RecruiterHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.reservationService.DeleteRecruiter(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(w, r, err)
		default:
			response.InternalServerError(w, r, err)
		}
		return
	}
}
