package v1

import (
	"net/http"
)

// GetTest
// @Summary		tests by param
// @Tags			Test
// @Description	get tests by param
// @Produce		json
// @Param			param	query		string	true	"param"
// @Success		200			{object}	dto.Test
// @Failure		404			{object}	dto.Test
// @Failure		400			{object}	http.response
// @Failure		500 		{object} 	http.response
// @Router			/v1/test [get]
func (h *Handler) GetTest(w http.ResponseWriter, r *http.Request) {
}
