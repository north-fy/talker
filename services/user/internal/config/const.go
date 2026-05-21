package config

import (
	"os"

	"golang.org/x/crypto/bcrypt"
)

const (
	CostHash = bcrypt.DefaultCost
)

var JwtSecret = os.Getenv("JWT_SECRET")
