#!/bin/bash

echo "Building Index Service ..."
go build -o bin/ cmd/index/index.go

echo "Building Auth Service ..."
go build -o bin/ cmd/auth/auth.go

echo "Building Product Service ..."
go build -o bin/ cmd/product/product.go

echo "Building Customer Service  ..."
go build -o bin/ cmd/customer/customer.go

echo "Build completed!"
