package dtos

type TransactionDashboardResponseDTO struct {
	IncomeAmount  float64 `bson:"incomeAmount" json:"incomeAmount"`
	ExpenseAmount float64 `bson:"expenseAmount" json:"expenseAmount"`
	TotalAmount   float64 `bson:"totalAmount" json:"totalAmount"`
}
