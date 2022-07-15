package transaction

import (
	"context"
	"database/sql"
	"strconv"
	"testing"

	"github.com/Melon-Network-Inc/payment-service/pkg/entity"
)

type mockRepository struct {
	items []entity.Transaction
}

func (m mockRepository) Get(ctx context.Context, id string) (entity.Transaction, error) {
	// tests logics in service (what api test does)
	for _, item := range m.items {
		if strconv.FormatUint(uint64(item.Id), 10) == id {
			return item, nil
		}
	}
	return entity.Transaction{}, sql.ErrNoRows // get {Id, etc}
}

// unit test
func TestRegisterHandlers(t *testing.T) {
	// RegisterHandlers(r *mux.Router, service Service, db *gorm.DB, logger log.Logger)
	// expected := "Transaction Received!"
	// actual := RegisterHandlers()
	// if actual != expected {
	// 	t.Errorf("expected %q but got %q", expected, actual)
	// }
}
