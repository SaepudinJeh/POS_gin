package services_test

import (
	"context"
	"testing"

	"POS/internal/repositories"
	"POS/internal/services"
)

func setupCategoryService() *services.CategoryService {
	mockRepo := repositories.NewMockCategoryRepository()
	return services.NewCategoryService(mockRepo)
}

func TestCategoryService_Create(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		inputName   string
		inputDesc   string
		prep        func(svc *services.CategoryService)
		expectError bool
		errorMsg    string
	}{
		{
			name:      "sukses create kategori",
			inputName: "Makanan",
			inputDesc: "Kategori makanan",
		},
		{
			name:      "gagal karena nama duplikat",
			inputName: "Makanan",
			inputDesc: "Duplikat",
			prep: func(svc *services.CategoryService) {
				_, _ = svc.Create(ctx, "Makanan", "Sudah ada")
			},
			expectError: true,
			errorMsg:    "nama kategori sudah ada",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupCategoryService()
			if tt.prep != nil {
				tt.prep(svc)
			}

			cat, err := svc.Create(ctx, tt.inputName, tt.inputDesc)

			if tt.expectError {
				if err == nil {
					t.Fatal("harusnya error, tapi nil")
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("pesan error salah, dapat: %s", err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("harusnya sukses, dapat: %v", err)
			}
			if cat.Name != tt.inputName {
				t.Errorf("nama salah, dapat: %s", cat.Name)
			}
		})
	}
}

func TestCategoryService_GetAll(t *testing.T) {
	ctx := context.Background()
	svc := setupCategoryService()

	_, _ = svc.Create(ctx, "Makanan", "Kategori makanan")
	_, _ = svc.Create(ctx, "Minuman", "Kategori minuman")

	list, err := svc.GetAll(ctx)
	if err != nil {
		t.Fatalf("harusnya sukses, dapat: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("harusnya 2, dapat: %d", len(list))
	}
}

func TestCategoryService_GetByID(t *testing.T) {
	ctx := context.Background()
	svc := setupCategoryService()

	cat, _ := svc.Create(ctx, "Makanan", "Kategori makanan")

	got, err := svc.GetByID(ctx, cat.ID)
	if err != nil {
		t.Fatalf("harusnya sukses, dapat: %v", err)
	}
	if got.ID != cat.ID {
		t.Errorf("ID salah, dapat: %d", got.ID)
	}
}

func TestCategoryService_Update(t *testing.T) {
	ctx := context.Background()
	svc := setupCategoryService()

	cat, _ := svc.Create(ctx, "Makanan", "Kategori makanan")

	updated, err := svc.Update(ctx, cat.ID, "Minuman", "Kategori minuman")
	if err != nil {
		t.Fatalf("harusnya sukses, dapat: %v", err)
	}
	if updated.Name != "Minuman" {
		t.Errorf("nama salah: %s", updated.Name)
	}
}

func TestCategoryService_Delete(t *testing.T) {
	ctx := context.Background()
	svc := setupCategoryService()

	cat, _ := svc.Create(ctx, "Makanan", "Kategori makanan")

	if err := svc.Delete(ctx, cat.ID); err != nil {
		t.Fatalf("harusnya sukses, dapat: %v", err)
	}
	if _, err := svc.GetByID(ctx, cat.ID); err == nil {
		t.Error("setelah delete, GetByID harusnya error")
	}
}
