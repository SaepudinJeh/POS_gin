package services_test

import (
	"POS/internal/repositories"
	"POS/internal/services"
	"os"
	"testing"
)

func setupAuthService() *services.AuthService {
	mockRepo := repositories.NewMockUserRepository()
	return services.NewAuthService(mockRepo)
}

func TestAuthService_Register_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-min-32-characters-long!!")
	svc := setupAuthService()

	user, err := svc.Register("Budi", "budi@test.com", "rahasia123")
	if err != nil {
		t.Fatalf("harusnya sukses, dapat error: %v", err)
	}

	if user.Email != "budi@test.com" {
		t.Errorf("email salah, dapat: %s", user.Email)
	}

	if user.Password == "rahasia123" {
		t.Errorf("password harusnya sudah di-hash, tapi masih plain")
	}

	if user.Role != "user" {
		t.Errorf("role default harusnya 'user', dapat: %s", user.Role)
	}
}

func TestAuthService_Register_EmailDuplicate(t *testing.T) {
	svc := setupAuthService()

	// Register pertama kali - harus sukses
	_, err := svc.Register("Budi", "budi@test.com", "rahasia123")
	if err != nil {
		t.Fatalf("register pertama harusnya sukses: %v", err)
	}

	// Register dengan email yang sama - harus gagal
	_, err = svc.Register("Andi", "budi@test.com", "password456")
	if err == nil {
		t.Fatal("register dengan email duplikat harusnya gagal")
	}

	if err.Error() != "email sudah terdaftar" {
		t.Errorf("pesan error salah, dapat: %s", err.Error())
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-min-32-characters-long!!")
	svc := setupAuthService()

	// Daftar dulu
	_, err := svc.Register("Budi", "budi@test.com", "rahasia123")
	if err != nil {
		t.Fatalf("register gagal: %v", err)
	}

	// Login
	token, err := svc.Login("budi@test.com", "rahasia123")
	if err != nil {
		t.Fatalf("login harusnya sukses, dapat error: %v", err)
	}

	if token == "" {
		t.Error("token harusnya tidak kosong")
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-min-32-characters-long!!")
	svc := setupAuthService()

	_, _ = svc.Register("Budi", "budi@test.com", "rahasia123")

	_, err := svc.Login("budi@test.com", "password_salah")
	if err == nil {
		t.Fatal("login dengan password salah harusnya gagal")
	}
}

func TestAuthService_Login_EmailNotFound(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-min-32-characters-long!!")
	svc := setupAuthService()

	_, err := svc.Login("tidak-ada@test.com", "apapun")
	if err == nil {
		t.Fatal("login dengan email tidak terdaftar harusnya gagal")
	}
}
