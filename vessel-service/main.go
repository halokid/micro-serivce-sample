package main

import (
  pb "shippy/vessel-service/proto/vessel"
  "github.com/micro/go-micro"
  "log"
  "os"
)

const (
  DEFAULT_HOST = "localhost:27017"
)

func main() {

  host := os.Getenv("DB_HOST")
  if host == "" {
    host = DEFAULT_HOST
  }
  session, err := CreateSession(host)
  defer session.Close()
  if err != nil {
    log.Fatalf("create session error: %v\n", err)
  }


  // 停留在港口的货船，先写死
  repo := &VesselRepository{session.Copy()}
  CreateDummyData(repo)
  server := micro.NewService(
    micro.Name("go.micro.srv.vessel"),
    micro.Version("latest"),
  )
  server.Init()

  // 将实现服务端的 API 注册到服务端
  pb.RegisterVesselServiceHandler(server.Server(), &handler{session})

  if err := server.Run(); err != nil {
    log.Fatalf("failed to serve: %v", err)
  }
}

func CreateDummyData(repo Repository)  {
  defer repo.Close()
  vessels := []*pb.Vessel{
    {Id: "vessel001", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500},
  }
  for _, v := range vessels {
    repo.Create(v)
  }
}


