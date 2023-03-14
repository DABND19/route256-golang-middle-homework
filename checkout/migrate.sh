#!/bin/sh

goose -dir ./migrations postgres "postgres://user:hackme123@localhost:5430/checkout?sslmode=disable" status

goose -dir ./migrations postgres "postgres://user:hackme123@localhost:5430/checkout?sslmode=disable" up
