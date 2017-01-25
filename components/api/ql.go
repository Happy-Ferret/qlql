package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/cznic/ql"
	"github.com/gernest/alien"
)

type Result struct {
	Error error       `json:"error"`
	Data  interface{} `json:"data"`
}

type Server struct {
	dbs []*ql.DB
}

func (s *Server) Open(w http.ResponseWriter, r *http.Request) {
}
func (s *Server) AllDB(w http.ResponseWriter, r *http.Request) {
	var a []string
	for _, v := range s.dbs {
		a = append(a, v.Name())
	}
	renderJSON(w, a)
}
func (s *Server) Info(w http.ResponseWriter, r *http.Request) {
	p := alien.GetParams(r)
	dbName := p.Get("dbname")
	if dbName == "" {
		rst := &Result{Error: errors.New("missing database name")}
		d, _ := json.Marshal(rst)
		w.Write(d)
		return
	}
	for _, v := range s.dbs {
		if v.Name() == dbName {
			info, err := v.Info()
			if err != nil {
				renderErr(w, err)
				return
			}
			//renderJSON(w, infoDB(info))
			renderJSON(w, info)
			return
		}
	}
}

//func infoDB(db *ql.DbInfo) common.DBInfo {
//i := common.DBInfo{}
//i.Name = db.Name
//for _, v := range db.Tables {
//i.Tables = append(i.Tables, infoTable(v))
//}
//for _, v := range db.Indices {
//x := common.Index{}
//x.Name = v.Name
//x.Table = v.Table
//x.Column = v.Column
//x.ExpressionList = v.ExpressionList
//x.Unique = v.Unique
//i.Indices = append(i.Indices, x)
//}
//return i
//}

//func infoTable(t ql.TableInfo) common.Table {
//i := common.Table{}
//i.Name = t.Name
//for _, v := range t.Columns {
//c := common.Column{}
//c.Name = v.Name
//c.Constraint = v.Constraint
//c.Default = v.Default
//c.NotNull = v.NotNull
//c.Type = v.Type.String()
//i.Columns = append(i.Columns, c)
//}
//return i
//}

func renderErr(w http.ResponseWriter, err error) {
	rst := &Result{Error: err}
	d, _ := json.Marshal(rst)
	w.Write(d)
}

func renderJSON(w http.ResponseWriter, data interface{}) {
	d, err := json.Marshal(data)
	if err != nil {
		renderErr(w, err)
		return
	}
	w.Write(d)
}

func NewServer() *alien.Mux {
	s := &Server{}
	s.dbs = append(s.dbs, testDB())
	r := alien.New()
	r.Get("/all", s.AllDB)
	r.Get("/info/:dbname", s.Info)
	r.Post("/open/:dbname", s.Open)
	return r
}

func testDB() *ql.DB {
	db, err := ql.OpenMem()
	if err != nil {
		log.Fatal(err)
	}
	schema := `
BEGIN TRANSACTION;
    CREATE TABLE users (id int64,age int64,user_num int64,name string,email string,birthday time,created_at time,updated_at time,billing_address_id int64,shipping_address_id int64,latitude float64,company_id int,role string,password_hash blob,sequence uint ) ;
    CREATE TABLE user_languages (user_id uint,language_id uint ) ;
    CREATE TABLE emails (id int16,user_id int,email string,created_at time,updated_at time ) ;
    CREATE TABLE languages (id uint,created_at time,updated_at time,deleted_at time,name string ) ;
    CREATE INDEX idx_languages_deleted_at ON languages(deleted_at);
    CREATE TABLE companies (id int64,name string ) ;
    CREATE TABLE credit_cards (id int8,number string,user_id int64,created_at time NOT NULL,updated_at time,deleted_at time ) ;
    CREATE TABLE addresses (id int,address1 string,address2 string,post string,created_at time,updated_at time,deleted_at time ) ;
COMMIT;`
	_, _, err = db.Run(ql.NewRWCtx(), schema)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
