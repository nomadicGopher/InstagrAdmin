#!/bin/bash
#rm ./*.log
go run main.go -username="$(cat username)" -access_token="$(cat access_token)"