package services_test

import (
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
	catSvc := services.NewCategoryService(catRepo)
	cat, err := catSvc.Create(name, "")
	if err != nil {
		t.Fatalf("gagal bikin kategori: %v", err)
	}
	return cat.ID
}

func TestProductService_Create_Success(t *testing.T) {
	svc, catRepo := setupProductService()
	catID := makeCategory(t, catRepo, "Makanan")

	product, err := svc.Create(services.ProductInput{
		SKU:        "MKN-001",
		Name:       "Indomie Goreng",
		Price:      3500,
		Cost:       2500,
		Stock:      100,
		CategoryID: catID,
	})
	if err != nil {
		t.Fatalf("harusnya sukses, dapat: %v", err)
	}
	if product.SKU != "MKN-001" {
		t.Errorf("SKU salah, dapat: %s", product.SKU)
	}
	if !product.IsActive {
		t.Error("produk baru harusnya aktif")
	}
}

func TestProductService_Create_DuplicateSKU(t *testing.T) {
	svc, catRepo := setupProductService()
	catID := makeCategory(t, catRepo, "Makanan")

	_, _ = svc.Create(services.ProductInput{
		SKU: "MKN-001", Name: "Indomie", CategoryID: catID,
	})

	_, err := svc.Create(services.ProductInput{
		SKU: "MKN-001", Name: "Mie Sedaap", CategoryID: catID,
	})
	if err == nil {
		t.Fatal("harusnya error karena SKU duplikat")
	}
	if err.Error() != "SKU sudah dipakai" {
		t.Errorf("pesan error salah: %s", err.Error())
	}
}

func TestProductService_Create_CategoryNotFound(t *testing.T) {
	svc, _ := setupProductService()

	_, err := svc.Create(services.ProductInput{
		SKU: "MKN-001", Name: "Indomie", CategoryID: 999,
	})
	if err == nil {
		t.Fatal("harusnya error karena kategori tidak ada")
	}
	if err.Error() != "kategori tidak ditemukan" {
		t.Errorf("pesan error salah: %s", err.Error())
	}
}

func TestProductService_Update_Success(t *testing.T) {
	svc, catRepo := setupProductService()
	catID := makeCategory(t, catRepo, "Makanan")

	p, _ := svc.Create(services.ProductInput{
		SKU: "MKN-001", Name: "Indomie", Price: 3500, CategoryID: catID,
	})

	updated, err := svc.Update(p.ID, services.ProductInput{
		SKU: "MKN-001", Name: "Indomie Goreng Spesial", Price: 4000, CategoryID: catID,
	})
	if err != nil {
		t.Fatalf("harusnya sukses, dapat: %v", err)
	}
	if updated.Name != "Indomie Goreng Spesial" {
		t.Errorf("nama salah: %s", updated.Name)
	}
	if updated.Price != 4000 {
		t.Errorf("harga salah: %f", updated.Price)
	}
}

func TestProductService_Update_DuplicateSKU(t *testing.T) {
	svc, catRepo := setupProductService()
	catID := makeCategory(t, catRepo, "Makanan")

	_, _ = svc.Create(services.ProductInput{SKU: "MKN-001", Name: "A", CategoryID: catID})
	p2, _ := svc.Create(services.ProductInput{SKU: "MKN-002", Name: "B", CategoryID: catID})

	_, err := svc.Update(p2.ID, services.ProductInput{
		SKU: "MKN-001", Name: "B", CategoryID: catID,
	})
	if err == nil {
		t.Fatal("harusnya error karena SKU duplikat")
	}
	if err.Error() != "SKU sudah dipakai" {
		t.Errorf("pesan error salah: %s", err.Error())
	}
}

func TestProductService_Update_SameSKU_NoError(t *testing.T) {
	svc, catRepo := setupProductService()
	catID := makeCategory(t, catRepo, "Makanan")

	p, _ := svc.Create(services.ProductInput{
		SKU: "MKN-001", Name: "Indomie", CategoryID: catID,
	})

	_, err := svc.Update(p.ID, services.ProductInput{
		SKU: "MKN-001", Name: "Indomie Baru", CategoryID: catID,
	})
	if err != nil {
		t.Fatalf("harusnya sukses, dapat: %v", err)
	}
}

func TestProductService_Update_NotFound(t *testing.T) {
	svc, catRepo := setupProductService()
	catID := makeCategory(t, catRepo, "Makanan")

	_, err := svc.Update(999, services.ProductInput{
		SKU: "MKN-001", Name: "A", CategoryID: catID,
	})
	if err == nil {
		t.Fatal("harusnya error")
	}
}

func TestProductService_Delete_Success(t *testing.T) {
	svc, catRepo := setupProductService()
	catID := makeCategory(t, catRepo, "Makanan")

	p, _ := svc.Create(services.ProductInput{SKU: "MKN-001", Name: "A", CategoryID: catID})

	if err := svc.Delete(p.ID); err != nil {
		t.Fatalf("harusnya sukses, dapat: %v", err)
	}
	if _, err := svc.GetByID(p.ID); err == nil {
		t.Error("setelah delete, GetByID harusnya error")
	}
}

func TestProductService_Delete_NotFound(t *testing.T) {
	svc, _ := setupProductService()

	if err := svc.Delete(999); err == nil {
		t.Fatal("harusnya error")
	}
}

func TestProductService_GetAll_WithFilter(t *testing.T) {
	svc, catRepo := setupProductService()
	cat1 := makeCategory(t, catRepo, "Makanan")
	cat2 := makeCategory(t, catRepo, "Minuman")

	_, _ = svc.Create(services.ProductInput{SKU: "MKN-001", Name: "Indomie", CategoryID: cat1})
	_, _ = svc.Create(services.ProductInput{SKU: "MKN-002", Name: "Mie Sedaap", CategoryID: cat1})
	_, _ = svc.Create(services.ProductInput{SKU: "MNM-001", Name: "Aqua", CategoryID: cat2})

	t.Run("filter by category", func(t *testing.T) {
		products, total, err := svc.GetAll(repositories.ProductFilter{CategoryID: &cat1})
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
		_, total, err := svc.GetAll(repositories.ProductFilter{Search: "indomie"})
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if total != 1 {
			t.Errorf("harusnya 1, dapat: %d", total)
		}
	})

	t.Run("preload category", func(t *testing.T) {
		products, _, _ := svc.GetAll(repositories.ProductFilter{})
		for _, p := range products {
			if p.Category == nil {
				t.Errorf("produk %s harusnya punya Category ter-preload", p.SKU)
			}
		}
	})
}
