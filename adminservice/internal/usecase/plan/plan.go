package plan

import (
	repo "adminservice/internal/adapter/db/sql"
	"adminservice/internal/adapter/kafka/producers"
	"adminservice/internal/entity"
	"adminservice/pkg"
	"log"
	"net/http"
)

type UseCase struct {
	repo repo.PlanRepo
}

func New(r repo.PlanRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (uc *UseCase) PlanSeasons(w http.ResponseWriter, r *http.Request) (entity.Season, error) {
	// Checking Auth Responce
	if err := pkg.CheckToken(r); err != nil {
		return entity.Season{}, err
	}

	var season entity.DetailSeasonJson
	var dbSeason entity.Season

	// parsing season body
	if err := pkg.ParseSeasonBody(r.Body, &season); err != nil {
		return entity.Season{}, err
	}

	// Serialize Season JSON entity in DB entity
	if err := pkg.StoreDeatailSeasonInDBEntity(&season, &dbSeason); err != nil {
		return entity.Season{}, err
	}

	// Finding if seasons are crossing
	if err := uc.repo.FindSeasonCross(&dbSeason); err != nil {
		return entity.Season{}, err
	}

	// If season is not crossing with others we add it to DB
	if err := uc.repo.AddNewSeason(&dbSeason); err != nil {
		return entity.Season{}, err
	}

	log.Println("Season IN DB: ", dbSeason)

	// Produsing new season to Core service
	go producers.SendSeasonInfo(pkg.ParseSeasonToKafkaJSON(dbSeason))

	return dbSeason, nil
}
