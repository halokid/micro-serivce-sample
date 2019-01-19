package  main

import (
  "context"
  "log"
  "net"
  pb "shippy/consignment-service/proto/consignment"
  "google.golang.org/grpc"
)

const (
  PORT = ":50051"
)

type IRepository interface {
  Create(consignment *pb.Consignment) (*pb.Consignment, error)

  GetAll() []*pb.Consignment
}


type Repositry struct {
  consignments []*pb.Consignment
}

func (repo *Repositry) Create(consignment *pb.Consignment) (*pb.Consignment, error)  {
  repo.consignments = append(repo.consignments, consignment)
  return consignment, nil
}

func (repo *Repositry) GetAll() []*pb.Consignment  {
  return repo.consignments
}

type service struct {
  repo Repositry
}


func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
  consignment, err := s.repo.Create(req)
  if err != nil {
    return nil, err
  }
  resp := &pb.Response{Created: true, Consignment: consignment}
  return resp, nil
}


func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
  allCOnsignments := s.repo.GetAll()
  resp := &pb.Response{Consignments: allCOnsignments}
  return resp, nil
}

func main() {
  listener, err := net.Listen("tcp", PORT)
  if err != nil {
    log.Fatal("failed to listen: %v", err)
  }
  log.Printf("listen on: %s\n", PORT)

  server := grpc.NewServer()
  repo := Repositry{}
  pb.RegisterShippingServiceServer(server, &service{repo})

  if err := server.Serve(listener); err != nil {
    log.Fatal("failed to server: %v", err)
  }
}





