package services_test

import (
	"context"
	"testing"

	"POS/internal/models"
	"POS/internal/repositories"
	"POS/internal/services"
)

func setupTxnService() (*services.TransactionService, *repositories.MockProductRepository, *repositories.MockCategoryRepository, *repositories.MockSettingRepository) {
	catRepo := repositories.NewMockCategoryRepository()
	prodRepo := repositories.NewMockProductRepository(catRepo)
	stockRepo := repositories.NewMockStockRepository(prodRepo)
	txnRepo := repositories.NewMockTransactionRepository(prodRepo, stockRepo)
	settingRepo := repositories.NewMockSettingRepository()

	_ = settingRepo.Upsert(context.Background(), "tax_rate", "11")
	_ = settingRepo.Upsert(context.Background(), "invoice_prefix", "INV")

	settingSvc := services.NewSettingService(settingRepo)
	svc := services.NewTransactionService(txnRepo, prodRepo, settingSvc)

	return svc, prodRepo, catRepo, settingRepo
}

func seedTxnProduct(t *testing.T, prodRepo *repositories.MockProductRepository, catRepo *repositories.MockCategoryRepository, sku, name string, price float64, stock int) uint {
	t.Helper()
	ctx := context.Background()
	cat := &models.Category{Name: "Test"}
	_ = catRepo.Create(ctx, cat)
	p := &models.Product{
		SKU: sku, Name: name, Price: price, Stock: stock,
		CategoryID: cat.ID, IsActive: true,
	}
	_ = prodRepo.Create(ctx, p)
	return p.ID
}

func TestTransactionService_Create_Pending(t *testing.T) {
	ctx := context.Background()
	svc, prodRepo, catRepo, _ := setupTxnService()
	pid := seedTxnProduct(t, prodRepo, catRepo, "MKN-001", "Indomie", 3500, 100)

	txn, err := svc.Create(ctx, services.TransactionInput{
		Items:     []services.TransactionItemInput{{ProductID: pid, Quantity: 2}},
		CashierID: 1,
	})
	if err != nil {
		t.Fatalf("harusnya sukses: %v", err)
	}
	if txn.Status != models.TransactionStatusPending {
		t.Errorf("status harusnya pending, dapat: %s", txn.Status)
	}
	if txn.Subtotal != 7000 {
		t.Errorf("subtotal harusnya 7000, dapat: %f", txn.Subtotal)
	}
	if txn.TaxAmount != 770 {
		t.Errorf("tax harusnya 770, dapat: %f", txn.TaxAmount)
	}
	if txn.Total != 7770 {
		t.Errorf("total harusnya 7770, dapat: %f", txn.Total)
	}
	if txn.InvoiceNumber == "" {
		t.Error("invoice harus terisi")
	}

	// Stok berkurang
	p, _ := prodRepo.FindByID(ctx, pid)
	if p.Stock != 98 {
		t.Errorf("stok harusnya 98, dapat: %d", p.Stock)
	}
}

func TestTransactionService_Create_WithPayment(t *testing.T) {
	ctx := context.Background()
	svc, prodRepo, catRepo, _ := setupTxnService()
	pid := seedTxnProduct(t, prodRepo, catRepo, "MKN-001", "Indomie", 3500, 100)

	txn, err := svc.Create(ctx, services.TransactionInput{
		Items:         []services.TransactionItemInput{{ProductID: pid, Quantity: 2}},
		PaymentMethod: models.PaymentMethodCash,
		PaymentAmount: 10000,
		CashierID:     1,
	})
	if err != nil {
		t.Fatalf("harusnya sukses: %v", err)
	}
	if txn.Status != models.TransactionStatusPaid {
		t.Errorf("status harusnya paid, dapat: %s", txn.Status)
	}
	if txn.ChangeAmount != 10000-7770 {
		t.Errorf("kembalian salah, dapat: %f", txn.ChangeAmount)
	}
	if txn.PaidAt == nil {
		t.Error("paid_at harus terisi")
	}
}

func TestTransactionService_Create_WithDiscount(t *testing.T) {
	ctx := context.Background()
	svc, prodRepo, catRepo, _ := setupTxnService()
	pid := seedTxnProduct(t, prodRepo, catRepo, "MKN-001", "Indomie", 10000, 100)

	txn, err := svc.Create(ctx, services.TransactionInput{
		Items:         []services.TransactionItemInput{{ProductID: pid, Quantity: 1}},
		DiscountType:  models.DiscountTypePercentage,
		DiscountValue: 10,
		CashierID:     1,
	})
	if err != nil {
		t.Fatalf("harusnya sukses: %v", err)
	}
	if txn.DiscountAmount != 1000 {
		t.Errorf("diskon harusnya 1000, dapat: %f", txn.DiscountAmount)
	}
	// Total = 10000 - 1000 + (9000 * 11%) = 9000 + 990 = 9990
	if txn.Total != 9990 {
		t.Errorf("total harusnya 9990, dapat: %f", txn.Total)
	}
}

func TestTransactionService_Create_InsufficientStock(t *testing.T) {
	ctx := context.Background()
	svc, prodRepo, catRepo, _ := setupTxnService()
	pid := seedTxnProduct(t, prodRepo, catRepo, "MKN-001", "Indomie", 3500, 5)

	_, err := svc.Create(ctx, services.TransactionInput{
		Items:     []services.TransactionItemInput{{ProductID: pid, Quantity: 10}},
		CashierID: 1,
	})
	if err == nil {
		t.Fatal("harusnya error stok tidak mencukupi")
	}
}

func TestTransactionService_Create_DuplicateProduct(t *testing.T) {
	ctx := context.Background()
	svc, prodRepo, catRepo, _ := setupTxnService()
	pid := seedTxnProduct(t, prodRepo, catRepo, "MKN-001", "Indomie", 3500, 100)

	_, err := svc.Create(ctx, services.TransactionInput{
		Items: []services.TransactionItemInput{
			{ProductID: pid, Quantity: 1},
			{ProductID: pid, Quantity: 2},
		},
		CashierID: 1,
	})
	if err == nil {
		t.Fatal("harusnya error produk duplikat")
	}
}

func TestTransactionService_Pay(t *testing.T) {
	ctx := context.Background()
	svc, prodRepo, catRepo, _ := setupTxnService()
	pid := seedTxnProduct(t, prodRepo, catRepo, "MKN-001", "Indomie", 3500, 100)

	txn, _ := svc.Create(ctx, services.TransactionInput{
		Items:     []services.TransactionItemInput{{ProductID: pid, Quantity: 2}},
		CashierID: 1,
	})

	paid, err := svc.Pay(ctx, txn.ID, models.PaymentMethodCash, 10000)
	if err != nil {
		t.Fatalf("harusnya sukses: %v", err)
	}
	if paid.Status != models.TransactionStatusPaid {
		t.Errorf("status harusnya paid, dapat: %s", paid.Status)
	}
	if paid.ChangeAmount != 10000-paid.Total {
		t.Errorf("kembalian salah")
	}
}

func TestTransactionService_Pay_InsufficientAmount(t *testing.T) {
	ctx := context.Background()
	svc, prodRepo, catRepo, _ := setupTxnService()
	pid := seedTxnProduct(t, prodRepo, catRepo, "MKN-001", "Indomie", 3500, 100)

	txn, _ := svc.Create(ctx, services.TransactionInput{
		Items:     []services.TransactionItemInput{{ProductID: pid, Quantity: 2}},
		CashierID: 1,
	})

	_, err := svc.Pay(ctx, txn.ID, models.PaymentMethodCash, 100)
	if err == nil {
		t.Fatal("harusnya error jumlah kurang")
	}
}

func TestTransactionService_Void_PendingRestoresStock(t *testing.T) {
	ctx := context.Background()
	svc, prodRepo, catRepo, _ := setupTxnService()
	pid := seedTxnProduct(t, prodRepo, catRepo, "MKN-001", "Indomie", 3500, 100)

	txn, _ := svc.Create(ctx, services.TransactionInput{
		Items:     []services.TransactionItemInput{{ProductID: pid, Quantity: 5}},
		CashierID: 1,
	})

	// Stok setelah create
	p, _ := prodRepo.FindByID(ctx, pid)
	if p.Stock != 95 {
		t.Fatalf("stok harusnya 95, dapat: %d", p.Stock)
	}

	voided, err := svc.Void(ctx, txn.ID, "Salah input", 1)
	if err != nil {
		t.Fatalf("harusnya sukses: %v", err)
	}
	if voided.Status != models.TransactionStatusVoid {
		t.Errorf("status harusnya void, dapat: %s", voided.Status)
	}
	if voided.VoidReason != "Salah input" {
		t.Errorf("reason salah: %s", voided.VoidReason)
	}

	// Stok kembali
	p, _ = prodRepo.FindByID(ctx, pid)
	if p.Stock != 100 {
		t.Errorf("stok harusnya kembali 100, dapat: %d", p.Stock)
	}
}

func TestTransactionService_Void_PaidRestoresStock(t *testing.T) {
	ctx := context.Background()
	svc, prodRepo, catRepo, _ := setupTxnService()
	pid := seedTxnProduct(t, prodRepo, catRepo, "MKN-001", "Indomie", 3500, 100)

	txn, _ := svc.Create(ctx, services.TransactionInput{
		Items:         []services.TransactionItemInput{{ProductID: pid, Quantity: 3}},
		PaymentMethod: models.PaymentMethodCash,
		PaymentAmount: 20000,
		CashierID:     1,
	})

	voided, err := svc.Void(ctx, txn.ID, "Customer batal", 1)
	if err != nil {
		t.Fatalf("harusnya sukses: %v", err)
	}
	if voided.Status != models.TransactionStatusVoid {
		t.Errorf("status harusnya void")
	}

	p, _ := prodRepo.FindByID(ctx, pid)
	if p.Stock != 100 {
		t.Errorf("stok harusnya kembali 100, dapat: %d", p.Stock)
	}
}

func TestTransactionService_Void_AlreadyVoid(t *testing.T) {
	ctx := context.Background()
	svc, prodRepo, catRepo, _ := setupTxnService()
	pid := seedTxnProduct(t, prodRepo, catRepo, "MKN-001", "Indomie", 3500, 100)

	txn, _ := svc.Create(ctx, services.TransactionInput{
		Items:     []services.TransactionItemInput{{ProductID: pid, Quantity: 1}},
		CashierID: 1,
	})

	_, _ = svc.Void(ctx, txn.ID, "Batal", 1)

	_, err := svc.Void(ctx, txn.ID, "Batal lagi", 1)
	if err == nil {
		t.Fatal("harusnya error karena sudah void")
	}
}
