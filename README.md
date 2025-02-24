# bank
## Overview

This project is a simple bank application written in Go. It allows users to create accounts, deposit and withdraw money, and check their balance.

## Features

- Create a new account
- Deposit money into an account
- Withdraw money from an account
- Check account balance

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/yogesh-mp/bank.git
    ```
2. Navigate to the project directory:
    ```sh
    cd bank
    ```
3. Install dependencies:
    ```sh
    go mod tidy
    ```

## Usage

1. Run the application:
    ```sh
    docker-compose up --build
    ```

2. From postman or Curl hit the URLS
- Create a new account
    ```sh
    POST /accounts/create
    Content-Type: application/json

    {
      "name": "John Doe",
      "initial_deposit": 1000
    }
    ```
- Deposit money into an account
    ```sh
    POST /transactions/deposit
    Content-Type: application/json

    {
      "account_id": 1,
      "amount": 500
    }
    ```
- Withdraw money from an account
    ```sh
    POST /transactions/withdraw
    Content-Type: application/json

    {
      "account_id": 1,
      "amount": 200
    }
    ```
- Check account balance
    ```sh
    GET /transactions/balance?id=3
    ```