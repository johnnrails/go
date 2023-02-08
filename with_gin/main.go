package main

import (
	// "github.com/gin-gonic/gin"
	"github.com/johnnrails/ddd_go/with_gin/common"
)

func AutoMigrate() {

}

func main() {
	db := common.Init()
	defer db.Close()

	// r := gin.Default()
}
