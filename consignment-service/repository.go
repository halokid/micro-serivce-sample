/**
 负责数据库的操作，ORM层
 */
package  main

import (
  pb "shippy/consignment-service/proto/consignment"
  "gopkg.in/mgo.v2"
)

const (
  DB_NAME           =       "shippy"
  CON_COLLECTION    =       "consignment"
)

type Repository interface {
  Create(*pb.Consignment) error
  GetAll()  ([]*pb.Consignment, error)
  Close()
}

type ConsignmentRepository struct {
  session *mgo.Session
}

func (repo *ConsignmentRepository) Create(c *pb.Consignment) error {
  return repo.collection().Insert()
}

func (repo *ConsignmentRepository) GetAll() ([]*pb.Consignment, error)  {
  var cons []*pb.Consignment
  err := repo.collection().Find(nil).All(&cons)
  return cons, err
}

func (repo *ConsignmentRepository) Close() {
  repo.session.Close()
}

func (repo *ConsignmentRepository) collection() *mgo.Collection {
  return repo.session.DB(DB_NAME).C(CON_COLLECTION)
}