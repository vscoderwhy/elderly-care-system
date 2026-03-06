package service

import (
	"elderly-care-system/internal/model"
	"elderly-care-system/internal/repository"
	"errors"
	"time"
)

type BillService struct {
	billRepo *repository.BillRepository
}

func NewBillService(billRepo *repository.BillRepository) *BillService {
	return &BillService{billRepo: billRepo}
}

func (s *BillService) List(elderlyID uint, page, pageSize int) ([]model.Bill, int64, error) {
	offset := (page - 1) * pageSize
	return s.billRepo.FindByElderlyID(elderlyID, offset, pageSize)
}

func (s *BillService) ListAll(page, pageSize int) ([]model.Bill, int64, error) {
	offset := (page - 1) * pageSize
	return s.billRepo.List(offset, pageSize)
}

func (s *BillService) Get(id uint) (*model.Bill, error) {
	return s.billRepo.FindByID(id)
}

func (s *BillService) Pay(id uint, amount float64, method string, transactionNo string) error {
	bill, err := s.billRepo.FindByID(id)
	if err != nil {
		return err
	}

	if bill.Status == "paid" {
		return errors.New("bill already paid")
	}

	// Create payment
	payment := &model.Payment{
		BillID:        id,
		Amount:        amount,
		Method:        method,
		TransactionNo: transactionNo,
		PaidAt:        time.Now(),
	}

	if err := s.billRepo.CreatePayment(payment); err != nil {
		return err
	}

	// Update bill status
	return s.billRepo.UpdateStatus(id, "paid")
}
