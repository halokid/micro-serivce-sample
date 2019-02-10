package main

import (
  "github.com/jinzhu/gorm"
  pb "shippy/user-service/proto/user"
)

type Repository interface {
  Get(id string) (*pb.User, error)
  GetAll() ([]*pb.User, error)
  Create(*pb.User) error
  GetByEmailAndPassword(*pb.User) (*pb.User, error)
}

type UserRepository struct {
  db *gorm.DB
}

func (repo *UserRepository) Get(id string) (*pb.User, error) {
  var u *pb.User
  u.Id = id
  if err := repo.db.First(&u).Error; err != nil {
    return nil, err
  }
  return u, nil
}

func (repo *UserRepository) GetAll() ([]*pb.User, error) {
  var users []*pb.User
  if err := repo.db.Find(&users).Error; err != nil {
    return nil, err
  }
  return users, nil
}

func (repo *UserRepository) Create(u *pb.User) error {
  if err := repo.db.Create(&u).Error; err != nil {
    return err
  }
  return nil
}

func (repo *UserRepository) GetByEmailAndPassword(u *pb.User) (*pb.User, error) {
  if err := repo.db.Find(&u).Error; err != nil {
    return nil, err
  }
  return u, nil
}



