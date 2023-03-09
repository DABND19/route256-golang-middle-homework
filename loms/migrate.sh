#!/bin/sh

goose -dir ./migrations postgres "postgres://user:hackme123@localhost:5431/loms?sslmode=disable" status

goose -dir ./migrations postgres "postgres://user:hackme123@localhost:5431/loms?sslmode=disable" up
