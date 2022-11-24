package account

import "ctigroupjsc.com/phuocnn/gps-management/database/repository"

type DeleteAccountHandler struct {
	AccountRepository     repository.AccountRepository
	UserRepository        repository.UserRepository
	ActivityLogRepository repository.ActivityLogRepository
	ReportRepository      repository.ReportRepository
}

func (h *DeleteAccountHandler) Handle(id string) error {
	account, err := h.AccountRepository.FindByID(id)
	if err != nil {
		return err
	}

	err = h.ReportRepository.RemoveByAccountID(account.ID.Hex())
	if err != nil {
		return err
	}

	err = h.ActivityLogRepository.RemoveByAccountID(account.ID.Hex())
	if err != nil {
		return err
	}

	err = h.UserRepository.RemoveByID(account.UserID)
	if err != nil {
		return err
	}

	return h.AccountRepository.RemoveByID(account.ID.Hex())
}
