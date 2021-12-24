package usecase

import "finance/adaptor"

var categoryAdaptor adaptor.CategoryAdaptorInterface
var transactionAdaptor adaptor.TransactionAdaptorInterface
var userAdaptor adaptor.UserAdaptorInterface
var walletAdaptor adaptor.WalletAdaptorInterface

func InitUsecase() {
	categoryAdaptor = adaptor.NewCategoryMySQLAdaptor()
	transactionAdaptor = adaptor.NewTransactionMySQLAdaptor()
	userAdaptor = adaptor.NewUserMySQLAdaptor()
	walletAdaptor = adaptor.NewWalletMySQLAdaptor()
}
