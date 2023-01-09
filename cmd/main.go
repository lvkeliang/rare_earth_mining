package main

import (
	"rare_earth_mining_BE/api"
	"rare_earth_mining_BE/dao"
)

func main() {
	dao.InitDB()
	api.InitRouter()
}
