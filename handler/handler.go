package handler

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strconv"

	"github.com/attisvalentin/tasklist/dbhandler"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Username string
	Access   string
	ID       string
	Company  string
}

func Firstinit() {
	dbhandler.Firstinit()
}

func LoadLogin(c *gin.Context) {
	_, err := c.Cookie("tasklistcookie")
	if err != nil {
		c.HTML(http.StatusOK, "login.html", nil)
		return
	}

	c.Redirect(http.StatusSeeOther, "/index")
}

func Login(c *gin.Context) {
	user := c.PostForm("username")
	password := c.PostForm("password")
	company := c.PostForm("company")
	userpassword, newUserID := dbhandler.LoginHandler(user, password, company)
	hashedPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	if hashedPassword == userpassword {
		c.SetCookie("tasklistcookie", newUserID, 86400, "/", "", false, true)
		c.Redirect(http.StatusSeeOther, "/index")
	} else {
		c.Redirect(http.StatusSeeOther, "/")
	}
}

func Logout(c *gin.Context) {
	c.SetCookie("tasklistcookie", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

func LoadRegisterForm(c *gin.Context) {
	c.HTML(http.StatusOK, "registration.html", nil)
}

func Register(c *gin.Context) {
	newUser := []string{c.PostForm("username"), fmt.Sprintf("%x", sha256.Sum256([]byte(c.PostForm("password")))), c.PostForm("company"), c.PostForm("name")}
	dbhandler.RegisterHandler(newUser)
	c.Redirect(http.StatusSeeOther, "/")
}

func AddUser(c *gin.Context) {
	user := []string{}
	user = append(user, c.PostForm("username"), fmt.Sprintf("%x", sha256.Sum256([]byte(c.PostForm("password")))), c.PostForm("access"), c.PostForm("name"), c.PostForm("company"))
	dbhandler.AddNewUser(user)
	c.Redirect(http.StatusSeeOther, "/index")
}

func ReturnTaskList(c *gin.Context) {
	user_id, err := c.Cookie("tasklistcookie")
	var allUser []string
	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}
	user, access, company := dbhandler.GetUser(user_id)
	if access == "admin" {
		allUser = dbhandler.GetAllUser(company)
	}

	response := dbhandler.ReturnData(user, access, company)
	addUser := User{Username: user, Access: access, ID: user_id, Company: company}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"response": response,
		"user":     addUser,
		"allusers": allUser,
	})
}
// TODO: change the addtask and the changestatus to unique for the companies.
func AddTask(c *gin.Context) {
	persontask := []string{c.PostForm("name"), c.PostForm("game"), c.PostForm("task"), c.PostForm("priority"), "In progress", c.PostForm("comment"), c.PostForm("company")}
	dbhandler.UpdateTable(persontask)
	c.Redirect(http.StatusSeeOther, "/index")
}

func ChangeStatus(c *gin.Context) {
	status := c.PostForm("status")
	id, _ := strconv.ParseFloat(c.PostForm("taskid"), 64)
	comment := c.PostForm("comment")
	company := c.PostForm("company")
	dbhandler.ChangeStatus(status, comment, company, int(id))
	c.Redirect(http.StatusSeeOther, "/index")
}
