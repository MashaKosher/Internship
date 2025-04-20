package service

import "net/http"

// @Summary		Get season statistic
// @Description	Get season statistic
// @Tags			Statistic
// @Accept			json
// @Produce		json
// @Success		200
// @Router			/statistic/players [get]
func Players(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Players statistic"))
}

// @Summary		Get seasons statistic
// @Description	Get seasons statistic
// @Tags			Statistic
// @Accept			json
// @Produce		json
// @Success		200
// @Router			/statistic/seasons [get]
func Seasons(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Season statistic"))
}
