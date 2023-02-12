package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/johnnrails/ddd_go/swaggo/web"
)

// @Summary		Add a new pet to the store
// @Description	get string by ID
// @Accept			json
// @Produce		json
// @Param			some_id	path		int				true	"Some ID"
// @Success		200		{string}	string			"ok"
// @Failure		400		{object}	web.APIError	"We need ID!!"
// @Failure		404		{object}	web.APIError	"Can not find ID"
// @Router			/api/a/{some_id} [get]
func GetStringByInt(c *gin.Context) {
	err := web.APIError{}
	fmt.Println(err)
}

// @Description	get struct array by ID
// @Accept			json
// @Produce		json
// @Param			some_id	path		string			true	"Some ID"
// @Param			offset	query		int				true	"Offset"
// @Param			limit	query		int				true	"Limit"
// @Success		200		{string}	string			"ok"
// @Failure		400		{object}	web.APIError	"We need ID!!"
// @Failure		404		{object}	web.APIError	"Can not find ID"
// @Router			/api/b/{some_id} [get]
func GetStructArrayByString(c *gin.Context) {

}

type Pet3 struct {
	ID int `json:"id"`
}
