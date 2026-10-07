package services_test

import (
	"testing"

	"POS/internal/repositories"
	"POS/internal/services"
)

func setupCategoryService() *services.CategoryService {
	mockRepo := repositories.NewMockCategoryRepository()
	return services.NewCategoryService(mockRepo)
}

// ---------- CREATE ----------

func TestCategoryService_Create(t *testing.T) {
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
				_, _ = svc.Create("Makanan", "Sudah ada")
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

			cat, err := svc.Create(tt.inputName, tt.inputDesc)

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
			if cat.ID == 0 {
				t.Error("ID harusnya terisi")
			}
		})
	}
}

// ---------- GET ALL ----------

func TestCategoryService_GetAll_Empty(t *testing.T) {
	svc := setupCategoryService()

	list, err := svc.GetAll()
	if err != nil {
		t.Fatalf("harusnya sukses, dapat: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("harusnya kosong, dapat: %d", len(list))
	}
}

func TestCategoryService_GetAll_WithData(t *testing.T) {
	svc := setupCategoryService()

	_, _ = svc.Create("Makanan", "Kategori makanan")
	_, _ = svc.Create("Minuman", "Kategori minuman")

	list, err := svc.GetAll()
	if err != nil {
		t.Fatalf("harusnya sukses, dapat: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("harusnya 2 kategori, dapat: %d", len(list))
	}
}

// ---------- GET BY ID ----------

func TestCategoryService_GetByID(t *testing.T) {
	tests := []struct {
		name        string
		prep        func(svc *services.CategoryService) uint
		expectError bool
	}{
		{
			name: "sukses get by id",
			prep: func(svc *services.CategoryService) uint {
				cat, _ := svc.Create("Makanan", "Kategori makanan")
				return cat.ID
			},
		},
		{
			name: "gagal karena tidak ditemukan",
			prep: func(svc *services.CategoryService) uint {
				return 999
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupCategoryService()
			id := tt.prep(svc)

			cat, err := svc.GetByID(id)

			if tt.expectError {
				if err == nil {
					t.Fatal("harusnya error, tapi nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("harusnya sukses, dapat: %v", err)
			}
			if cat.ID != id {
				t.Errorf("ID salah, dapat: %d", cat.ID)
			}
		})
	}
}

// ---------- UPDATE ----------

func TestCategoryService_Update(t *testing.T) {
	tests := []struct {
		name        string
		prep        func(svc *services.CategoryService) (uint, string, string)
		expectError bool
		errorMsg    string
	}{
		{
			name: "sukses update",
			prep: func(svc *services.CategoryService) (uint, string, string) {
				cat, _ := svc.Create("Makanan", "Kategori makanan")
				return cat.ID, "Minuman", "Kategori minuman"
			},
		},
		{
			name: "gagal karena tidak ditemukan",
			prep: func(svc *services.CategoryService) (uint, string, string) {
				return 999, "Apapun", "Apapun"
			},
			expectError: true,
			errorMsg:    "kategori tidak ditemukan",
		},
		{
			name: "gagal karena nama duplikat dengan kategori lain",
			prep: func(svc *services.CategoryService) (uint, string, string) {
				_, _ = svc.Create("Minuman", "Kategori minuman")
				cat, _ := svc.Create("Makanan", "Kategori makanan")
				return cat.ID, "Minuman", "Coba duplikat"
			},
			expectError: true,
			errorMsg:    "nama kategori sudah ada",
		},
		{
			name: "sukses update dengan nama sendiri (tidak dianggap duplikat)",
			prep: func(svc *services.CategoryService) (uint, string, string) {
				cat, _ := svc.Create("Makanan", "Kategori makanan")
				return cat.ID, "Makanan", "Deskripsi baru"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupCategoryService()
			id, name, desc := tt.prep(svc)

			cat, err := svc.Update(id, name, desc)

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
			if cat.Name != name {
				t.Errorf("nama salah, dapat: %s", cat.Name)
			}
		})
	}
}

// ---------- DELETE ----------

func TestCategoryService_Delete_Success(t *testing.T) {
	svc := setupCategoryService()
	cat, _ := svc.Create("Makanan", "Kategori makanan")

	err := svc.Delete(cat.ID)
	if err != nil {
		t.Fatalf("harusnya sukses, dapat: %v", err)
	}
}

func TestCategoryService_Delete_NotFound(t *testing.T) {
	svc := setupCategoryService()

	err := svc.Delete(999)
	if err == nil {
		t.Fatal("harusnya error, tapi nil")
	}
}

// ---------- SOFT DELETE BEHAVIOR ----------

func TestCategoryService_Delete_IsSoftDelete(t *testing.T) {
	mockRepo := repositories.NewMockCategoryRepository()
	svc := services.NewCategoryService(mockRepo)

	cat, _ := svc.Create("Makanan", "Kategori makanan")
	_ = svc.Delete(cat.ID)

	// 1. Data masih ada di storage
	stored, ok := mockRepo.Categories[cat.ID]
	if !ok {
		t.Fatal("row harusnya masih ada di storage")
	}
	if !stored.DeletedAt.Valid {
		t.Error("DeletedAt harusnya terisi")
	}

	// 2. GetAll tidak menampilkannya
	list, _ := svc.GetAll()
	if len(list) != 0 {
		t.Errorf("GetAll harusnya kosong, dapat: %d", len(list))
	}

	// 3. GetByID gagal
	_, err := svc.GetByID(cat.ID)
	if err == nil {
		t.Error("GetByID untuk data terhapus harusnya error")
	}

	// 4. Nama boleh dipakai lagi
	_, err = svc.Create("Makanan", "Baru")
	if err != nil {
		t.Errorf("harusnya bisa buat nama sama setelah soft delete, dapat: %v", err)
	}
}
