package loan

import (
	"fmt"
	"sync"
)

type Loan struct {
	ID              int
	Amount          float64
	InterestRate    float64
	Weeks           int
	Outstanding     float64
	LastPaidWeek    int
	CurrentWeek     int
	WeeklyDue       float64
	BillingSchedule map[int]float64
}

type Loans struct {
	Mu     sync.Mutex
	Loans  map[int]Loan
	NextID int
}

var instance *Loans
var once sync.Once

func GetInstance() *Loans {
	once.Do(func() {
		instance = &Loans{
			Loans:  make(map[int]Loan),
			NextID: 1,
		}
	})
	return instance
}

func (l *Loans) CreateLoan(amount float64, annualInterestRate float64, weeks int) int {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	weeklyInterestRate := annualInterestRate / 52.14 // Assuming 52 weeks in a year
	totalRepayable := amount * (1 + weeklyInterestRate*float64(weeks))
	weeklyDue := totalRepayable / float64(weeks)

	billingSchedule := make(map[int]float64)
	for i := 1; i <= weeks; i++ {
		billingSchedule[i] = weeklyDue
	}

	loan := Loan{
		ID:              l.NextID,
		Amount:          amount,
		InterestRate:    annualInterestRate,
		Weeks:           weeks,
		Outstanding:     totalRepayable,
		CurrentWeek:     0,
		LastPaidWeek:    0,
		WeeklyDue:       weeklyDue,
		BillingSchedule: billingSchedule,
	}

	l.Loans[l.NextID] = loan
	l.NextID++

	// Print Billing Schedule
	fmt.Printf("Loan %d Created - Total Repayable: %.2f, Weekly Payment: %.2f\n", loan.ID, totalRepayable, weeklyDue)
	fmt.Println("Billing Schedule:")
	for week, amount := range loan.BillingSchedule {
		fmt.Printf("W%d : %.2f\n", week, amount)
	}

	return loan.ID
}

func (l *Loans) GetOutstanding(loanID int) float64 {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	if loan, exists := l.Loans[loanID]; exists {
		return loan.Outstanding
	}
	return 0
}

func (l *Loans) IsDelinquent(loanID int) bool {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	if loan, exists := l.Loans[loanID]; exists {
		return loan.CurrentWeek-loan.LastPaidWeek >= 2
	}
	return false
}

func (l *Loans) MakePayment(loanID int) bool {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	if loan, exists := l.Loans[loanID]; exists {
		amount := loan.WeeklyDue
		if amount == loan.BillingSchedule[loan.CurrentWeek] {
			loan.Outstanding -= amount
			loan.LastPaidWeek = loan.CurrentWeek
			l.Loans[loanID] = loan

			fmt.Printf("Payment successful for Loan ID %d, Week %d. New outstanding balance: %.2f\n", loanID, loan.CurrentWeek, loan.Outstanding)
			return true
		} else {
			fmt.Println("Invalid payment amount.")
			return false
		}
	} else {
		fmt.Println("Loan not found.")
		return false
	}
}
