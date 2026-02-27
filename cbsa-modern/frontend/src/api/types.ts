// TypeScript type definitions matching the Go backend models.
// These replace the untyped JSON objects used in the original JavaScript frontend.

export interface Customer {
  id: number;
  customerNumber: string;
  sortCode: string;
  title: string;
  firstName: string;
  lastName: string;
  addressLine1: string;
  addressLine2: string;
  addressLine3: string;
  dateOfBirth: string;
  creditScore: number;
  reviewDate?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateCustomerRequest {
  title: string;
  firstName: string;
  lastName: string;
  addressLine1: string;
  addressLine2: string;
  addressLine3: string;
  dateOfBirth: string;
}

export interface UpdateCustomerRequest {
  title?: string;
  firstName?: string;
  lastName?: string;
  addressLine1?: string;
  addressLine2?: string;
  addressLine3?: string;
}

export interface Account {
  id: number;
  accountNumber: string;
  sortCode: string;
  customerNumber: string;
  accountType: string;
  interestRate: string;
  openedDate: string;
  overdraftLimit: number;
  lastStatement: string;
  nextStatement: string;
  availableBalance: string;
  actualBalance: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateAccountRequest {
  customerNumber: string;
  accountType: string;
  interestRate: string;
  overdraftLimit: number;
}

export interface UpdateAccountRequest {
  accountType?: string;
  interestRate?: string;
  overdraftLimit?: number;
}

export interface ProcessedTransaction {
  id: number;
  sortCode: string;
  accountNumber: string;
  txnDate: string;
  txnTime: string;
  reference: string;
  txnType: string;
  description: string;
  amount: string;
  createdAt: string;
}

export interface DebitCreditRequest {
  accountNumber: string;
  amount: string;
  description: string;
}

export interface TransferRequest {
  fromAccountNumber: string;
  toAccountNumber: string;
  amount: string;
  description: string;
}

export interface PaginatedResponse<T> {
  total: number;
  limit: number;
  offset: number;
  [key: string]: T[] | number;
}

export interface CustomerListResponse {
  customers: Customer[];
  total: number;
  limit: number;
  offset: number;
}

export interface AccountListResponse {
  accounts: Account[];
  total: number;
  limit: number;
  offset: number;
}

export interface TransactionListResponse {
  transactions: ProcessedTransaction[];
  total: number;
  limit: number;
  offset: number;
}
