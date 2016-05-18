package main

import (
    "html/template"
    "bytes"
    "os"
    "bufio"
    "gopkg.in/mgo.v2"
    "github.com/ghodss/yaml"
    "io/ioutil"
    "time"
    "github.com/gin-gonic/gin"
)

type Server struct {
    Hostname string `bson:"_id"`
    Updated time.Time `bson:"Updated"`
    Platops_Support string `bson:"Platops_Support"`
    Applications []string `bson:"Applications"`
    Pager_Playbooks []string `bson:"Pager_Playbooks"`
    Puppet_Modules []string `bson:"Puppet_Modules"`
    Purpose []string `bson:"Purpose"`
}


func check(e error) {
    if e != nil {
        panic(e)
    }
}

func Tofile(file, data string) {
    // Truncate file
    f, err := os.Create(file)
    check(err)
    defer f.Close()
    // Write lines to file
    w := bufio.NewWriter(f)
    w.WriteString(data)
    w.Flush()
    f.Sync()
}

type Config struct {
    Mongo_db string `json:"mongo_db"`
    Mongo_passwd string `json:"mongo_passwd"`
    Mongo_user string `json:"mongo_user"`
    Mongo_authdb string `json:"mongo_authdb"`
    Mongo_addr string `json:"mongo_addr"`
    Jsonstats string `json:"jsonstats"`
}

func config()  (mongo_db,mongo_passwd,mongo_user,mongo_authdb,mongo_addr string){
    var v Config
    config_file, err := ioutil.ReadFile("/etc/relevy/config.yaml")
    check(err)
    yaml.Unmarshal(config_file, &v)
    mongo_db = v.Mongo_db
    mongo_passwd = v.Mongo_passwd
    mongo_user = v.Mongo_user
    mongo_authdb = v.Mongo_authdb
    mongo_addr = v.Mongo_addr
    return
}


func main() {

    go func () {
        for {
            // Read Config, load values
            mongo_db,mongo_passwd,mongo_user,mongo_authdb,mongo_addr := config()

            // We need this object to establish a session to our MongoDB.
            mongoDBDialInfo := &mgo.DialInfo{
              Addrs:    []string{mongo_addr},
              Timeout:  60 * time.Second,
              Database: mongo_authdb,
              Username: mongo_user,
              Password: mongo_passwd,
            }

            // Create a session which maintains a pool of socket connections
            // to our MongoDB.
            mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
            check(err)
            server := []Server{}

            // Request a socket connection
            sessionCopy := mongoSession.Copy()
            coll := sessionCopy.DB(mongo_db).C("relevy")
            // Ensure documents will die after 1 hour of not checking in
            index := mgo.Index{
                Key: []string{"Updated"},
                ExpireAfter: (time.Hour * 1),
            }
            err = coll.EnsureIndex(index)
            if err != nil {
                err = coll.DropIndex("Updated")
                check(err)
                err = coll.EnsureIndex(index)
                check(err)
            }
            // Gather all results from mongo
            err = coll.Find(nil).All(&server)
            check(err)

            // Close session
            mongoSession.Close()
            sessionCopy.Close()

            // Generate servers.html from template
            var doc bytes.Buffer
            t, _ := template.ParseFiles("template.html")
            t.Execute(&doc, server)
            s := doc.String()
            Tofile("servers.html", s)
            time.Sleep(5 * time.Second)
        }
    }()
    // Http stuff
    router := gin.Default()
    router.StaticFile("/sorttable.js", "./sorttable.js")
    router.StaticFile("/", "./servers.html")

    // Listen and server on 0.0.0.0:8080
    router.Run(":8080")

}
