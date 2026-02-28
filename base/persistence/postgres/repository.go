package postgres

import (
	"github.com/fundathyon/base/persistence/contracts"
	"github.com/fundathyon/base/persistence/models"
	"gorm.io/gorm"
)

type PostgresRepository[E contracts.IEntity, M contracts.IModel] struct {
	Connection *gorm.DB
	OnMutation models.OnMutationFunc
}
