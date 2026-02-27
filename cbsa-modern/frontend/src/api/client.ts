import axios, { AxiosInstance } from 'axios';
import {
  Customer,
  CreateCustomerRequest,
  UpdateCustomerRequest,
  Account,
  CreateAccountRequest,
  UpdateAccountRequest,
  DebitCreditRequest,
  TransferRequest,
  CustomerListResponse,
  AccountListResponse,
  TransactionListResponse,
} from './types';

// Centralized, typed API client replacing the scattered axios.get/post calls
// found throughout the original JavaScript components. The original frontend
// had no API abstraction — each component imported axios and constructed URLs
// manually with string concatenation.

const api: AxiosInstance = axios.create({
  baseURL: process.env.REACT_APP_API_BASE_URL || 'http://localhost:8080/api/v1',
  headers: { 'Content-Type': 'application/json' },
});

export const customerApi = {
  list: (limit = 20, offset = 0) =>
    api.get<CustomerListResponse>(`/customers?limit=${limit}&offset=${offset}`).then(r => r.data),

  get: (customerNumber: string) =>
    api.get<Customer>(`/customers/${customerNumber}`).then(r => r.data),

  create: (data: CreateCustomerRequest) =>
    api.post<Customer>('/customers', data).then(r => r.data),

  update: (customerNumber: string, data: UpdateCustomerRequest) =>
    api.put<Customer>(`/customers/${customerNumber}`, data).then(r => r.data),

  delete: (customerNumber: string) =>
    api.delete(`/customers/${customerNumber}`).then(r => r.data),
};

export const accountApi = {
  list: (limit = 20, offset = 0) =>
    api.get<AccountListResponse>(`/accounts?limit=${limit}&offset=${offset}`).then(r => r.data),

  get: (accountNumber: string) =>
    api.get<Account>(`/accounts/${accountNumber}`).then(r => r.data),

  listByCustomer: (customerNumber: string) =>
    api.get<{ accounts: Account[]; total: number }>(`/accounts/customer/${customerNumber}`).then(r => r.data),

  create: (data: CreateAccountRequest) =>
    api.post<Account>('/accounts', data).then(r => r.data),

  update: (accountNumber: string, data: UpdateAccountRequest) =>
    api.put<Account>(`/accounts/${accountNumber}`, data).then(r => r.data),

  delete: (accountNumber: string) =>
    api.delete(`/accounts/${accountNumber}`).then(r => r.data),
};

export const transactionApi = {
  list: (limit = 20, offset = 0) =>
    api.get<TransactionListResponse>(`/transactions?limit=${limit}&offset=${offset}`).then(r => r.data),

  debit: (data: DebitCreditRequest) =>
    api.put<Account>('/transactions/debit', data).then(r => r.data),

  credit: (data: DebitCreditRequest) =>
    api.put<Account>('/transactions/credit', data).then(r => r.data),

  transfer: (data: TransferRequest) =>
    api.put('/transactions/transfer', data).then(r => r.data),
};

export const infoApi = {
  companyName: () =>
    api.get<{ name: string }>('/company').then(r => r.data.name),

  sortCode: () =>
    api.get<{ sortCode: string }>('/sortcode').then(r => r.data.sortCode),
};
