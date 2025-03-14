Loan Billing System

Overview

This is a Loan Billing System written in Go that provides functionality to:

Generate a billing schedule for a loan (weekly payments over 52 weeks)

Track the outstanding balance of a loan

Detect if a customer is delinquent (missed 2 consecutive payments)

Allow customers to make payments on their loans

The system ensures fast retrieval of loans using a **map-based singleton struct**, allowing efficient management of multiple loans.