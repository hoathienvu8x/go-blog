package model

import (
    "database/sql"
    "errors"
    "gopkg.in/mgo.v2"
)

var (
    ErrCode         = errors.New("Case statement in code is not correct.")
    ErrNoResult     = errors.New("Result not found.")
    ErrUnavailable  = errors.New("Database is unavailable.")
    ErrUnauthorized = errors.New("User does not have permission to perform this operation.") 
)

func standardizeError(err error) error {
    if err == sql.ErrNoRows || err == mgo.ErrNotFound {
        return ErrNoResult
    }
    return err
}
