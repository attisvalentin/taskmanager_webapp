package dbhandler

import (
	"database/sql"
	"fmt"

	//"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type Database struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

type NewUser struct {
	Username string
	Password string
	Access   string
	ID       string
	Name     string
}

type Task struct {
	Name     string
	Game     string
	Task     string
	Priority string
	Status   string
	ID       int
	Comment  string
	Date     string
	Company  string
}

// var taskID int
var tasks []Task
var db *sql.DB

// -------------- Database start --------------

func Firstinit() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	data := Database{
		Username: viper.Get("username").(string),
		Password: viper.Get("password").(string),
		Host:     viper.Get("host").(string),
		Port:     viper.Get("port").(string),
		Database: viper.Get("database").(string),
	}
	fmt.Println(data.Username)
	dbConn(data)
	// taskID = ReturnTaskIndex()
}

func dbConn(data Database) {
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", data.Username, data.Password, data.Host, data.Port, data.Database))
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

func ReturnTaskIndex() int {
	var id int
	err := db.QueryRow("SELECT taskindex FROM vars").Scan(&id)
	if err != nil {
		panic(err)
	}
	return id
}

// -------------- Login methods --------------

func LoginHandler(user, password, company string) (string, string) {
	var userpassword string
	var userID string
	dbuser, err := db.Query("SELECT password, ID FROM users WHERE username=? AND company=?", user, company)
	if err != nil {
		panic(err)
	}

	for dbuser.Next() {
		err := dbuser.Scan(&userpassword, &userID)
		if err != nil {
			panic(err)
		}
	}

	return userpassword, userID
}

// -------------- Registration  --------------

func RegisterHandler(newuser []string) {
	username := newuser[0]
	password := newuser[1]
	company := newuser[2]
	name := newuser[3]
	CreateCompanyTable(company)
	AddNewUser([]string{username, password, "admin", name, company})
}

// -------------- Create the company's tasklist table, and insert the company taskindex to the vars --------------

func CreateCompanyTable(company string) {
	dbcreatecomp, err := db.Prepare(fmt.Sprintf("CREATE TABLE %s (name VARCHAR(255), game VARCHAR(255), task VARCHAR(255), priority VARCHAR(255), status VARCHAR(255), ID INT, comment VARCHAR(255), start VARCHAR(255))", company))
	if err != nil {
		panic(err)
	}
	dbcreatecomp.Exec()
	insertNewtaskID(company)
}

func insertNewtaskID(company string) {
	addtaskid, err := db.Prepare(fmt.Sprintf("ALTER TABLE vars ADD COLUMN %s INT NOT NULL DEFAULT 0", company))
	if err != nil {
		panic(err)
	}
	addtaskid.Exec()
}
// -------------- User handling --------------

func GetAllUser(company string) []string {
	dbNameQuery, err := db.Query(fmt.Sprintf("SELECT name FROM users WHERE company=\"%s\" ORDER BY name ASC", company))
	if err != nil {
		panic(err)
	}
	var userList []string
	for dbNameQuery.Next() {
		var name string
		err := dbNameQuery.Scan(&name)
		if err != nil {
			panic(err)
		}
		userList = append(userList, name)
	}

	return userList
}

func GetUser(id string) (string, string, string) {
	dbVarQuery, err := db.Query("SELECT name, access, company FROM users WHERE ID=?", id)
	if err != nil {
		panic(err)
	}
	var name string
	var access string
	var company string
	for dbVarQuery.Next() {
		err = dbVarQuery.Scan(&name, &access, &company)
		if err != nil {
			panic(err.Error())
		}
	}
	return name, access, company
}

func AddNewUser(user []string) {
	dbinsertuser, err := db.Prepare("INSERT INTO users(username, password, access, name, ID, company) VALUES( ?, ?, ?, ?, ?, ? )")
	if err != nil {
		panic(err)
	}
	dbinsertuser.Exec(user[0], user[1], user[2], user[3], GetNewID(), user[4])
}

func GetNewID() string {
	var id string
	for {
		id = uuid.New().String()
		_, err := db.Query("SELECT * FROM users WHERE ID=\"%s\"", id)
		if err != nil {
			break
		} else {
			continue
		}
	}
	return id
}

// -------------- Returning data about tasks ----------------

func ReturnData(name, acces, company string) []Task {
	var sql string
	if acces == "admin" {
		sql = fmt.Sprintf("SELECT * FROM %s ORDER BY ID DESC", company)
	} else {
		sql = fmt.Sprintf("SELECT * FROM %s WHERE name=\"%s\" ORDER BY ID DESC", company, name)
	}
	dbselect, err := db.Query(sql)
	if err != nil {
		panic(err)
	}

	tasks = []Task{}

	for dbselect.Next() {

		var name, game, task, priority, status, comment, date string
		var id int
		err = dbselect.Scan(&name, &game, &task, &priority, &status, &id, &comment, &date)
		if err != nil {
			panic(err.Error())
		}
		taskdata := Task{
			Name:     name,
			Game:     game,
			Task:     task,
			Priority: priority,
			Status:   status,
			ID:       id,
			Comment:  comment,
			Date:     date,
			Company:  company,
		}

		tasks = append(tasks, taskdata)
	}

	return tasks
}

// -------------- Update the table ----------------

func UpdateTable(persontask []string) {
	dbinsert, err := db.Prepare(fmt.Sprintf("INSERT INTO %s(name, game, task, priority, status, ID, comment, start) VALUES( ?, ?, ?, ?, ?, ?, ?, ?)", persontask[6]))
	if err != nil {
		panic(err)
	}
	taskID := gettaskID(persontask[6])
	dbinsert.Exec(persontask[0], persontask[1], persontask[2], persontask[3], persontask[4], taskID, persontask[5], returnTime())

	updatevariables(persontask[6], taskID+1)
}

func gettaskID(company string) int {
	var id int
	dbselect, err := db.Query(fmt.Sprintf("SELECT %s FROM vars", company))
	if err != nil {
		panic(err)
	}
	for dbselect.Next() {
		err = dbselect.Scan(&id)
		if err != nil {
			panic(err.Error())
		}
	}
	return id
}

func updatevariables(company string, taskID int) {
	dbupdate, err := db.Prepare(fmt.Sprintf("UPDATE vars SET %s=?", company))
	if err != nil {
		panic(err)
	}
	dbupdate.Exec(taskID)
}

// -------------- Status changing ----------------

func ChangeStatus(status, comment, company string, id int) {
	var sql string
	if status != "deleted" {
		sql = fmt.Sprintf("UPDATE %s SET status=\"%s\", comment=\"%s\", start=\"%s\" WHERE ID=%d", company, status, comment, returnTime(), id)
	} else {
		sql = fmt.Sprintf("DELETE FROM %s WHERE ID=%d", company, id)
	}
	dbupdate, err := db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	dbupdate.Exec()
}

func returnTime() string {
	now := time.Now()
	cTime := now.Format("2006-01-02 15:04:05")
	return cTime
}
