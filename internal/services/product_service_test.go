package services_test

import (
	"context"
	"testing"

	"POS/internal/repositories"
	"POS/internal/services"
)

func setupProductService() (*services.ProductService, *repositories.MockCategoryRepository) {
	catRepo := repositories.NewMockCategoryRepository()
	prodRepo := repositories.NewMockProductRepository(catRepo)
	svc := services.NewProductService(prodRepo, catRepo)
	return svc, catRepo
}

func makeCategory(t *testing.T, catRepo *repositories.MockCategoryRepository, name string) uint {
	t.Helper()
	ctx := context.Background()
	catSvc := services.NewCategoryService(catRepo)
	cat, err := catSvc.Create(ctx, name, "")
	if err != nil {
		t.Fatalf("gagal bikin kategori: %v", err)
	}
	return cat.ID
}

func TestProductService_Create_Success(t *testing.T) {
	ctx := context.Background()
	svc, catRepo := setupProductService()
	catID := makeCategory(t, catRepo, "Makanan")

	product, err := svc.Create(ctx, services.ProductInput{
		SKU: "MKN-001", Name: "Indomie", Price: 3500, Cost: 2500, Stock: 100, CategoryID: catID,
	})
	if err != nil {
		t.Fatalf("harusnya sukses: %v", err)
	}
	if product.SKU != "MKN-001" {
		t.Errorf("SKU salah: %s", product.SKU)
	}
}

func TestProductService_Create_DuplicateSKU(t *testing.T) {
	ctx := context.Background()
	svc, catRepo := setupProductService()
	catID := makeCategory(t, catRepo, "Makanan")

	_, _ = svc.Create(ctx, services.ProductInput{SKU: "MKN-001", Name: "A", CategoryID: catID})

	_, err := svc.Create(ctx, services.ProductInput{SKU: "MKN-001", Name: "B", CategoryID: catID})
	if err == nil {
		t.Fatal("harusnya error")
	}
	if err.Error() != "SKU sudah dipakai" {
		t.Errorf("pesan error salah: %s", err.Error())
	}
}

func TestProductService_Create_CategoryNotFound(t *testing.T) {
	ctx := context.Background()
	svc, _ := setupProductService()

	_, err := svc.Create(ctx, services.ProductInput{SKU: "MKN-001", Name: "Indomie", CategoryID: 999})
	if err == nil {
		t.Fatal("harusnya error")
	}
}

func TestProductService_Update_Success(t *testing.T) {
	ctx := context.Background()
	svc, catRepo := setupProductService()
	catID := makeCategory(t, catRepo, "Makanan")

	p, _ := svc.Create(ctx, services.ProductInput{
		SKU: "MKN-001", Name: "Indomie", Price: 3500, CategoryID: catID,
	})

	updated, err := svc.Update(ctx, p.ID, services.ProductInput{
		SKU: "MKN-001", Name: "Indomie Spesial", Price: 4000, CategoryID: catID,
	})
	if err != nil {
		t.Fatalf("harusnya sukses: %v", err)
	}
	if updated.Name != "Indomie Spesial" {
		t.Errorf("nama salah: %s", updated.Name)
	}
}

func TestProductService_Delete_Success(t *testing.T) {
	ctx := context.Background()
	svc, catRepo := setupProductService()
	catID := makeCategory(t, catRepo, "Makanan")

	p, _ := svc.Create(ctx, services.ProductInput{SKU: "MKN-001", Name: "A", CategoryID: catID})

	if err := svc.Delete(ctx, p.ID); err != nil {
		t.Fatalf("harusnya sukses: %v", err)
	}
	if _, err := svc.GetByID(ctx, p.ID); err == nil {
		t.Error("setelah delete, GetByID harusnya error")
	}
}

func TestProductService_GetAll_WithFilter(t *testing.T) {
	ctx := context.Background()
	svc, catRepo := setupProductService()
	cat1 := makeCategory(t, catRepo, "Makanan")
	cat2 := makeCategory(t, catRepo, "Minuman")

	_, _ = svc.Create(ctx, services.ProductInput{SKU: "MKN-001", Name: "Indomie", CategoryID: cat1})
	_, _ = svc.Create(ctx, services.ProductInput{SKU: "MKN-002", Name: "Mie Sedaap", CategoryID: cat1})
	_, _ = svc.Create(ctx, services.ProductInput{SKU: "MNM-001", Name: "Aqua", CategoryID: cat2})

	t.Run("filter by category", func(t *testing.T) {
		products, total, err := svc.GetAll(ctx, repositories.ProductFilter{CategoryID: &cat1})
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if total != 2 {
			t.Errorf("harusnya 2, dapat: %d", total)
		}
		if len(products) != 2 {
			t.Errorf("harusnya 2, dapat: %d", len(products))
		}
	})

	t.Run("filter by search", func(t *testing.T) {
		_, total, err := svc.GetAll(ctx, repositories.ProductFilter{Search: "indomie"})
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if total != 1 {
			t.Errorf("harusnya 1, dapat: %d", total)
		}
	})
}
