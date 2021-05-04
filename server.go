package main

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	/*e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})*/

	//กำหนด Route ก่อนเลย พร้อมให้ call ไปยัง func ต่างๆ
	h := UserHandler{}
	h.Initialize()

	e.GET("/users", h.GetAllUser)
	e.POST("/users", h.SaveUser)
	e.GET("users/:id", h.GetUser)
	e.PUT("/users/:id", h.UpdateUser)
	e.DELETE("/users/:id", h.DeleteUser)

	e.Logger.Fatal(e.Start(":1323"))
}

type UserHandler struct {
	DB *gorm.DB
}

//ให้เชื่อมต่อฐานข้อมูลเมื่อ Initialize
func (h *UserHandler) Initialize() {
	db, err := gorm.Open("mysql", "root:0Ft5dxza@tcp(127.0.0.1:3306)/mysql")
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&User{})

	h.DB = db
}

type User struct {
	Id        uint   `gorm:"primary_key" json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
}

func (h *UserHandler) GetAllUser(c echo.Context) error {
	users := []User{}

	h.DB.Find(&users)

	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	id := c.Param("id")
	user := User{}

	if err := h.DB.Find(&user, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) SaveUser(c echo.Context) error {
	user := User{}

	if err := c.Bind(&user); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := h.DB.Save(&user); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	user := User{}

	if err := h.DB.Find(&user, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	if err := c.Bind(&user).Error; err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := h.DB.Save(&user).Error; err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	user := User{}

	if err := h.DB.Find(&user, id).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.NoContent(http.StatusNoContent)
}
