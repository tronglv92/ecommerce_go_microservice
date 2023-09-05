package customerrepo

import (
	"context"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
	customermodel "github.com/tronglv92/accounts/module/customer/model"
	"github.com/tronglv92/ecommerce_go_common/logger"
)

type GetCustomerByIdStorage interface {
	GetCustomer(
		context context.Context,
		cond map[string]interface{},
		moreKeys ...string) (*customermodel.Customer, error)
	GetAccountsFromCustomerId(
		context context.Context,
		customerId int,
		moreKeys ...string,
	) ([]customermodel.Account, error)
}

type GetCardsRemoteStorage interface {
	GetCardsFromCustomerId(
		context context.Context,

		id int,
	) ([]customermodel.Card, error)
}

type GetLoansRemoteStorage interface {
	GetLoansFromCustomerId(
		c *gin.Context,
		id int,
	) ([]customermodel.Loan, error)
}

type getCustomerRepo struct {
	store     GetCustomerByIdStorage
	cardStore GetCardsRemoteStorage
	loanStore GetLoansRemoteStorage
}

func NewCustomerByIdRepo(
	store GetCustomerByIdStorage,
	cardStore GetCardsRemoteStorage,
	loanStore GetLoansRemoteStorage,

) *getCustomerRepo {
	return &getCustomerRepo{
		store:     store,
		cardStore: cardStore,
		loanStore: loanStore,
	}
}
func (repo *getCustomerRepo) GetCustomerById(
	c *gin.Context,
	id int,
) (*customermodel.FullCustomer, error) {
	ctx := c.Request.Context()
	logger := logger.GetCurrent().GetLogger("account.repo.list_account")
	commandCard := "list_card_customer_by_id"
	commandLoan := "list_loan_customer_by_id"
	commandConfig := hystrix.CommandConfig{
		Timeout:                1000, // Timeout in milliseconds
		MaxConcurrentRequests:  10,   // Max concurrent requests allowed
		ErrorPercentThreshold:  50,   // Error percentage at which the circuit should open
		RequestVolumeThreshold: 3,    // Minimum number of requests needed for statistics
		SleepWindow:            5000, // Time to wait before allowing further requests after circuit opens

	}
	hystrix.ConfigureCommand(commandCard, commandConfig)
	hystrix.ConfigureCommand(commandLoan, commandConfig)
	customer, err := repo.store.GetCustomer(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return nil, err
	}

	accounts, err := repo.store.GetAccountsFromCustomerId(ctx, id)
	if err != nil {
		return nil, err
	}

	// With hystrix do

	var cards []customermodel.Card
	err = hystrix.Do(commandCard, func() error {
		// talk to dependency services
		cards, err = repo.cardStore.GetCardsFromCustomerId(ctx, id)
		if err != nil {
			return err
		}
		// cardsChan <- cards
		return nil

	}, nil)

	if err != nil {
		return nil, err
	}

	// with hystrix go
	loansChan := make(chan []customermodel.Loan)
	// errChan := make(chan error)
	var loans []customermodel.Loan
	errChan := hystrix.Go(commandLoan, func() error {
		// talk to dependency services
		loans, err := repo.loanStore.GetLoansFromCustomerId(c, id)
		if err != nil {
			return err
		}
		loansChan <- loans
		return nil

	}, nil)
	select {
	case loan := <-loansChan:
		loans = loan
	case err := <-errChan:
		if err != nil {
			logger.Errorf("err %v", err)
			return nil, err
		}

	}

	fullCustomer := customer.ConvertToFullCustomer(cards, accounts, loans)
	return fullCustomer, nil
}
