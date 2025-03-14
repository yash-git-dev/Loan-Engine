package main

import (
	"fmt"
	loan "loan-engine/Loans"
)

func main() {
	loanManager := loan.GetInstance()

	loanID1 := loanManager.CreateLoan(5000000, 0.10, 50)
	loanID2 := loanManager.CreateLoan(3000000, 0.10, 30)

	fmt.Printf("Loan %d Created. Outstanding: %.2f\n", loanID1, loanManager.GetOutstanding(loanID1))
	fmt.Printf("Loan %d Created. Outstanding: %.2f\n", loanID2, loanManager.GetOutstanding(loanID2))

	for week := 1; week <= 50; week++ {
		loanManager.MakePayment(loanID1)

		if loanManager.IsDelinquent(loanID1) {
			fmt.Println("Loan 1 is Delinquent!")
		}
	}

	fmt.Println("Final Outstanding for Loan 1:", loanManager.GetOutstanding(loanID1))
}
