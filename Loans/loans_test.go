package loan_test

import (
	loan "loan-engine/Loans"
	"testing"
)

func TestCreateLoan(t *testing.T) {
	loanManager := loan.GetInstance()
	loanID := loanManager.CreateLoan(5000000, 0.10, 50)

	loan, exists := loanManager.Loans[loanID]
	if !exists {
		t.Fatalf("Loan ID %d should exist but does not", loanID)
	}

	expectedWeeklyDue := (5000000 * (1 + (0.10 / 52.14 * 50))) / 50
	if loan.WeeklyDue != expectedWeeklyDue {
		t.Errorf("Expected weekly due %.2f but got %.2f", expectedWeeklyDue, loan.WeeklyDue)
	}

	if len(loan.BillingSchedule) != 50 {
		t.Errorf("Expected billing schedule of 52.14 weeks but got %d weeks", len(loan.BillingSchedule))
	}
}

func TestGetOutstanding(t *testing.T) {
	loanManager := loan.GetInstance()
	loanID := loanManager.CreateLoan(5000000, 0.10, 50)

	expectedOutstanding := 5000000 * (1 + (0.10 / 52.14 * 50))
	actualOutstanding := loanManager.GetOutstanding(loanID)

	if actualOutstanding != expectedOutstanding {
		t.Errorf("Expected outstanding balance %.2f but got %.2f", expectedOutstanding, actualOutstanding)
	}
}

func TestMakePaymentSuccess(t *testing.T) {
	loanManager := loan.GetInstance()
	loanID := loanManager.CreateLoan(5000000, 0.10, 50)

	loan := loanManager.Loans[loanID]
	loan.CurrentWeek = 1
	loanManager.Loans[loanID] = loan

	success := loanManager.MakePayment(loanID)
	if !success {
		t.Errorf("Expected payment to be successful but it failed")
	}

	expectedOutstanding := loan.Outstanding - loan.WeeklyDue
	actualOutstanding := loanManager.GetOutstanding(loanID)

	if actualOutstanding != expectedOutstanding {
		t.Errorf("Expected outstanding %.2f but got %.2f", expectedOutstanding, actualOutstanding)
	}
}

func TestMakePaymentIncorrectAmount(t *testing.T) {
	loanManager := loan.GetInstance()
	loanID := loanManager.CreateLoan(5000000, 0.10, 50)

	loan := loanManager.Loans[loanID]
	loan.CurrentWeek = 1
	loan.WeeklyDue = 120000
	loanManager.Loans[loanID] = loan

	success := loanManager.MakePayment(loanID)
	if success {
		t.Errorf("Expected payment failure due to incorrect amount, but it succeeded")
	}
}

func TestIsDelinquent(t *testing.T) {
	loanManager := loan.GetInstance()
	loanID := loanManager.CreateLoan(5000000, 0.10, 50)

	loan := loanManager.Loans[loanID]
	loan.CurrentWeek = 3
	loanManager.Loans[loanID] = loan

	if !loanManager.IsDelinquent(loanID) {
		t.Errorf("Expected loan to be delinquent but it was not detected")
	}

	loanManager.MakePayment(loanID)
	if loanManager.IsDelinquent(loanID) {
		t.Errorf("Expected loan to be non-delinquent after payment but it still is")
	}
}

func TestPartialPaymentsNotAllowed(t *testing.T) {
	loanManager := loan.GetInstance()
	loanID := loanManager.CreateLoan(5000000, 0.10, 50)

	loan := loanManager.Loans[loanID]
	loan.CurrentWeek = 1
	loanManager.Loans[loanID] = loan

	loan.WeeklyDue = loan.WeeklyDue / 2
	loanManager.Loans[loanID] = loan

	success := loanManager.MakePayment(loanID)
	if success {
		t.Errorf("Expected failure because partial payments are not allowed")
	}
}

func TestLoanFullyPaid(t *testing.T) {
	loanManager := loan.GetInstance()
	loanID := loanManager.CreateLoan(5000000, 0.10, 50)

	for week := 1; week <= 50; week++ {
		loan := loanManager.Loans[loanID]
		loan.CurrentWeek = week
		loanManager.Loans[loanID] = loan
		loanManager.MakePayment(loanID)
	}

	if outstanding := loanManager.GetOutstanding(loanID); outstanding > 0.0001 || outstanding < -0.0001 {
		t.Errorf("Expected outstanding balance to be 0 after full repayment, but got %.2f", outstanding)
	}
}

func TestLoanNotFound(t *testing.T) {
	loanManager := loan.GetInstance()
	outstanding := loanManager.GetOutstanding(999)

	if outstanding != 0 {
		t.Errorf("Expected outstanding balance to be 0 for non-existing loan but got %.2f", outstanding)
	}
}

func TestMultipleMissedPayments(t *testing.T) {
	loanManager := loan.GetInstance()
	loanID := loanManager.CreateLoan(5000000, 0.10, 50)

	loan := loanManager.Loans[loanID]
	loan.CurrentWeek = 5
	loanManager.Loans[loanID] = loan

	if !loanManager.IsDelinquent(loanID) {
		t.Errorf("Expected loan to be delinquent after 4 missed payments")
	}

	for i := 0; i < 3; i++ {
		loan.CurrentWeek++
		loanManager.Loans[loanID] = loan
		loanManager.MakePayment(loanID)
	}

	if loanManager.IsDelinquent(loanID) {
		t.Errorf("Loan should no longer be delinquent after consecutive payments")
	}
}
