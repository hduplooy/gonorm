# github/hduplooy/gonorm

## A basic no-orm for golang

This is currently really just an straight forward way to take a query and populate a structure, or a slice based on a structure or generate a json representation.

More functionality will be added over time for inserting, updating and deleting, etc.

### An example

Let's say you have a postgres database (can be anything as long as you have the driver for it) called *myerp* with a table defined as:

    create table public.clients (
        clientid serial,
        clientname varchar(50),
        clientsurname varchar(50),
        clientphone varchar(20),
        clientaddress varchar(100)
    );

Now you would like to generate a json representation of the data in that table. You would do it with the following code:

    package main

    import (
        "fmt"
        "log"

        "github.com/hduplooy/gonorm"
    )

    type Client struct {
        Clientid int    `json:"id"`
        Name     string `fldnm:"clientname"`
        Surname  string `fldnm:"clientsurname"`
        Phone    string `fldnm:"clientphone"`
    }

    func main() {
        db, err := gonorm.NewNorm("postgres", "host=localhost user=postgres password=passwd dbname=myerp sslmode=disable")
        if err != nil {
            log.Printf("Could not connect: %s", err.Error())
            return
        }
        result, err := db.GetRowsJson("select clientid, clientname, clientsurname,clientphone from clients", Client{})
        if err != nil {
            log.Printf("Could not retrieve data: %s", err.Error())
            return
        }
        fmt.Println(result)
    }

This will then print out the json as something like:

    [{"id":1,"Name":"Peter","Surname":"Pan","Phone":"+9(999)1234567"},
     {"id":2,"Name":"Cap","Surname":"Hook","Phone":"+6(666)9999999"},
     {"id":3,"Name":"Lily","Surname":"Tiger","Phone":"+9(999)1212121"}]


## API

### func NewNorm(driver, connect string) (*Norm, error)
NewNorm just connects to a database with the necessary details and return a *Norm structure and an error if any.  
The driver and connect strings are kept in the returned structure.

**driver** is the database type that one connects to (for example "postgres")  
**connect** is the connection string for the type of database with your database name, user, password etc.

### func (ent *Norm) GetRows(sql string, val interface{}) (interface{}, error)
GetRows will execute the sql query and then based on the type of val generate a slice with entries populated from each row returned.

**sql** is the sql query (select query normally)  
**val** is the val (an empty structure) which is used as template to generate the slice

### func (ent *Norm) GetRow(sql string, val interface{}) (interface{}, error)
GetRow similar to GetRows except it is for only the first row returned.

**sql** is the sql query (select query normally)  
**val** is the val (an empty structure) which is used as template to generate the slice

### func (ent *Norm) GetRowsJson(sql string, val interface{}) (string, error)
GetRowsJson get the rows from the sql query map them to the struct type of val and convert the slice to json.

**sql** is the sql query (select query normally)  
**val** is the val (an empty structure) which is used as template to generate the slice

### func (ent *Norm) GetRowJson(sql string, val interface{}) (string, error)
GetRowJson get the first row from the sql query map it to the struct type of val and convert to json.

**sql** is the sql query (select query normally)  
**val** is the val (an empty structure) which is used as template
