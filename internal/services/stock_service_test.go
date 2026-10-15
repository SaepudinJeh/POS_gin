package services_test

import (
	"context"
	"testing"

	"POS/internal/models"
	"POS/internal/repositories"
	"POS/internal/services"
)

func setupStockService() (*services.StockService, *repositories.MockProductRepository, *repositories.MockCategoryRepository) {
	catRepo := repositories.NewMockCategoryRepository()
	prodRepo := repositories.NewMockProductRepository(catRepo)
	stockRepo := repositories.NewMockStockRepository(prodRepo)
	svc := services.NewStockService(stockRepo, prodRepo)
	return svc, prodRepo, catRepo
}

func seedProduct(t *testing.T, prodRepo *repositories.MockProductRepository, catRepo *repositories.MockCategoryRepository, stock int) uint {
	t.Helper()
	ctx := context.Background()
	cat := &models.Category{Name: "Makanan"}
	_ = catRepo.Create(ctx, cat)

	prod := &models.Product{
		SKU: "MKN-001", Name: "Indomie", Price: 3500,
		Stock: stock, CategoryID: cat.ID, IsActive: true,
	}
	_ = prodRepo.Create(ctx, prod)
	return prod.ID
}

func TestStockService_CreateMutation_In(t *testing.T) {
	ctx := context.Background()
	svc, prodRepo, catRepo := setupStockService()
	prodID := seedProduct(t, prodRepo, catRepo, 10)

	mut, err := svc.CreateMutation(ctx, services.StockInput{
		ProductID: prodID, Type: models.StockTypeIn, Quantity: 50, Note: "Restock", UserID: 1,
	})
	if err != nil {
		t.Fatalf("harusnya sukses: %v", err)
	}
	if mut.StockBefore != 10 || mut.StockAfter != 60 {
		t.Errorf("stock salah: before=%d after=%d", mut.StockBefore, mut.StockAfter)
	}

	p, _ := prodRepo.FindByID(ctx, prodID)
	if p.Stock != 60 {
		t.Errorf("product stock harusnya 60, dapat: %d", p.Stock)
	}
}

func TestStockService_CreateMutation_Out(t *testing.T) {
	ctx := context.Background()
	svc, prodRepo, catRepo := setupStockService()
	prodID := seedProduct(t, prodRepo, catRepo, 100)

	mut, err := svc.CreateMutation(ctx, services.StockInput{
		ProductID: prodID, Type: models.StockTypeOut, Quantity: 30, UserID: 1,
	})
	if err != nil {
		t.Fatalf("harusnya sukses: %v", err)
	}
	if mut.StockAfter != 70 {
		t.Errorf("stock after harusnya 70, dapat: %d", mut.StockAfter)
	}
}

func TestStockService_CreateMutation_OutInsufficient(t *testing.T) {
	ctx := context.Background()
	svc, prodRepo, catRepo := setupStockService()
	prodID := seedProduct(t, prodRepo, catRepo, 10)

	_, err := svc.CreateMutation(ctx, services.StockInput{
		ProductID: prodID, Type: models.StockTypeOut, Quantity: 50, UserID: 1,
	})
	if err == nil {
		t.Fatal("harusnya error")
	}
	if err.Error() != "stok tidak mencukupi" {
		t.Errorf("pesan error salah: %s", err.Error())
	}
}

func TestStockService_CreateMutation_Adjustment(t *testing.T) {
	ctx := context.Background()
	svc, prodRepo, catRepo := setupStockService()
	prodID := seedProduct(t, prodRepo, catRepo, 100)

	mut, err := svc.CreateMutation(ctx, services.StockInput{
		ProductID: prodID, Type: models.StockTypeAdjustment, Quantity: 85, Note: "Opname", UserID: 1,
	})
	if err != nil {
		t.Fatalf("harusnya sukses: %v", err)
	}
	if mut.StockAfter != 85 {
		t.Errorf("stock harusnya 85, dapat: %d", mut.StockAfter)
	}
}

func TestStockService_GetByProductID(t *testing.T) {
	ctx := context.Background()
	svc, prodRepo, catRepo := setupStockService()
	prodID := seedProduct(t, prodRepo, catRepo, 10)

	_, _ = svc.CreateMutation(ctx, services.StockInput{
		ProductID: prodID, Type: models.StockTypeIn, Quantity: 20, UserID: 1,
	})
	_, _ = svc.CreateMutation(ctx, services.StockInput{
		ProductID: prodID, Type: models.StockTypeOut, Quantity: 5, UserID: 1,
	})

	mutations, err := svc.GetByProductID(ctx, prodID)
	if err != nil {
		t.Fatalf("harusnya sukses: %v", err)
	}
	if len(mutations) != 2 {
		t.Errorf("harusnya 2, dapat: %d", len(mutations))
	}
}
