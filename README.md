# go-simple-lending
An ovesimplified lending service app

### Features

1. Create Account

    Create account with particular amount of lending limit.

    ```
    POST /accounts

    {
        "name": string,
        "limit": number
    }
    ```
2. Get Account Detail

    Get account detail, including available limit, total loan amount, paid amount, and outstanding amount.

    ```
    GET /accounts/{id}
    ```

2. Create Loan

    Create loan for account with specified loan amount.

    ```
    POST /accounts/{id}/loans

    {
        "amount": number
    }
    ```

3. List of Loans

    Get list of loans of account

    ```
    GET /accounts/{id}/loans
    ```

4. List of Installments

    Get list of all installments of account

    ```
    GET /accounts/{id}/installments
    ```

5. Create Payment

    Crate payment for account to adjusting loan amount.

    ```
    POST /accounts/{id}/payments

    {
        "amount": number
    }
    ```

6. List of Payments

    Get list of all payments of account.

    ```
    GET /accounts/{id}/payments
    ```