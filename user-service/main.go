package main

import (
  "fmt"
  "github.com/micro/go-micro"
  "github.com/prometheus/common/log"
  pb "shippy/user-service/proto/user"
)

func main() {
  db, err := CreateConnection()

  fmt.Println("%+v\n", db)
  fmt.Println("err: %v\n", err)

  defer db.Close()

  if err != nil {
    log.Fatalf("db conn error: %v\n", err)
  }

  repo := &UserRepository{}

  db.AutoMigrate(&pb.User{})

  s := micro.NewService(
      micro.Name("go.micro.srv.user"),
      micro.Version("latest"),
    )

  s.Init()

  pb.RegisterUserServiceHandler(s.Server(), &handler{repo})

  if err := s.Run(); err != nil {
    log.Fatalf("user service error: %v\n", err)
  }
}
