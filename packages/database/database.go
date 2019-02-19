package database

import (
    "encoding/json"
    "fmt"
    "log"
    "time"
    "github.com/boltdb/bolt"
    _ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
    "gopkg.in/mgo.v2"
)

var (
    BoltDB      *bolt.DB
    Mongo       *mgo.Session
    SQL         *sqlx.DB
    databases   Info
)

type Type string

const (
    TypeBolt    Type = "Bolt"
    TypeMongoDB Type = "MongoDB"
    TypeMySQL   Type = "MySQL"
)

type Info struct {
    Type        Type
    MySQL       MySQLInfo
    Bolt        BoltInfo
    MongoDB     MongoDBInfo
}

type MySQLInfo struct {
    UserName    string
    Password    string
    Name        string
    HostName    string
    Port        int
    Parameter   string
}

type BoltInfo struct {
    Path    string
}

type MongoDBInfo struct {
    URL         string
    Database    string
}

func DSN(ci MySQLInfo) string {
    return ci.UserName + ":" + ci.Password + "@tcp(" + ci.HostName + ":" + fmt.Sprintf("%d", ci.Port) + ")/" + ci.Name + ci.Parameter
}

func Connect(d Info) {
    var err error
    databases = d
    switch d.Type {
        case TypeMySQL:
            if SQL, err = sqlx.Connect("mysql", DSN(d.MySQL)); err != nil {
                log.Println("SQL Driver Error", err)
            }
            if err = SQL.Ping(); err != nil {
                log.Println("Database Error", err)
            }
        case TypeBolt:
            if BoltDB, err = bolt.Open(d.Bolt.Path, 0600, nil); err != nil {
                log.Println("Bolt Driver Error", err)
            }
        case TypeMongoDB:
            if Mongo, err = mgo.DialWithTimeout(d.MongoDB.URL, 5 * time.Second); err != nil {
                log.Println("MongoDB Driver Error", err)
                return
            }
            Mongo.SetSocketTimeout(1 * time.Second)
            if err = Mongo.Ping(); err != nil {
                log.Println("Database Error", err)
            }
        default:
            log.Println("No registered database in config")
    }
}

func Update(bucketName string, key string, dataStruct interface{}) error {
    err := BoltDB.Update(func(tx *bolt.Tx) error {
        bucket, e := tx.CreateBucketIfNotExists([]byte(bucketName))
        if e != nil {
            return e
        }
        encodedRecord, e := json.Marshal(dataStruct)
        if e != nil {
            return e
        }
        if e = bucket.Put([]byte(key), encodedRecord); e != nil {
            return e
        }
        return nil
    })
    return err
}

func View(bucketName string, key string, dataStruct interface{}) error {
    err := BoltDB.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(bucketName))
        if b == nil {
            return bolt.ErrBucketNotFound
        }
        v := b.Get([]byte(key))
        if len(v) < 1 {
            return bolt.ErrInvalid
        }
        e := json.Unmarshal(v, &dataStruct)
        if e != nil {
            return e
        }
        return nil
    })
    return err
}

func Delete(bucketName string, key string) error {
    err := BoltDB.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(bucketName))
        if b == nil {
            return bolt.ErrBucketNotFound
        }
        return b.Delete([]byte(key))
    })
    return err
}

func CheckConnection() bool {
    if Mongo == nil {
        Connect(databases)
    }
    if Mongo != nil {
        return true
    }
    return false
}

func LoadConfig() Info {
    return databases
}
