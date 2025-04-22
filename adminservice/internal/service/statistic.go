package service

import (
	"adminservice/pkg"
	"net/http"
)

// @Summary		Get season statistic
// @Description	Get season statistic
// @Tags			Statistic
// @Accept			json
// @Produce		json
// @Success		200
// @Router			/statistic/players [get]
func Players(w http.ResponseWriter, r *http.Request) {
	// Checking Auth Responce
	if err := pkg.CheckToken(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("There must be Players statistic"))
}

// @Summary		Get seasons statistic
// @Description	Get seasons statistic
// @Tags			Statistic
// @Accept			json
// @Produce		json
// @Success		200
// @Router			/statistic/seasons [get]
func Seasons(w http.ResponseWriter, r *http.Request) {
	// Checking Auth Responce
	if err := pkg.CheckToken(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("There must be Season statistic"))
}
