package main

import "gopkg.in/mgo.v2"

func CreateSession(host string) (*mgo.Session, error) {
  s, err := mgo.Dial(host)
  if err != nil {
    return nil, err
  }
  return s, nil
}

