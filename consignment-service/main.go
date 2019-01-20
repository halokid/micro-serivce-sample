package main

import (
  "context"
  "github.com/micro/go-micro"
  "log"
  pb "shippy/consignment-service/proto/consignment"
  vesselPb "shippy/vessel-service/proto/vessel"
)


// ------------------ 方法的接口, 主要是流程 -----------------------------
type IRepository interface {
  Create(consignment *pb.Consignment, err error)
  GetAll() []*pb.Consignment
}


// ------------------- 数据结构体, 主要为数据源 ---------------------------
type Repository struct {
  consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
  repo.consignments = append(repo.consignments, consignment)
  return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment  {
  return repo.consignments
}


// -------------- 定义服务， 实现方法 ---------------------
type service struct {
  repo Repository
  vesselClient  vesselPb.VesselServiceClient
}

//func (s *service) CreateConsignment (ctx context.Context, req *pb.Consignment) (*pb.Consignment, error) {
  func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, resp *pb.Response) error {
  vReq := &vesselPb.Specification{
    Capacity:   int32(len(req.Containers)),
    MaxWeight:  req.Weight,
  }
  vResp, err := s.vesselClient.FindAvailable(context.Background(), vReq)
  if err != nil {
    return err
  }

  log.Fatalf("found vessel: %s\n", vResp.Vessel.Name)
  req.VesselId = vResp.Vessel.Id
  consignment, err := s.repo.Create(req)
  if err != nil {
    return err
  }

  resp.Created = true
  resp.Consignmet = consignment
  return nil
}

//func (s *service) GetConsignments (ctx context.Context, req *pb.GetRequest, resp *pb.Response) (*pb.Response, error)  {
  func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, resp *pb.Response) error {
  allConsignments := s.repo.GetAll()
  resp = &pb.Response{Consignments: allConsignments}
  return nil
}

func main() {
  server := micro.NewService(
    micro.Name("go.micro.srv.consignment"),
    micro.Version("latest"),
  )

  server.Init()
  repo := Repository{}
  vClient := vesselPb.NewVesselServiceClient("go.micro.srv.vessel", server.Client())
  pb.RegisterShippingServiceHandler(server.Server(), &service{repo, vClient})

  if err := server.Run(); err != nil {
    log.Fatalf("failed to server: %v", err)
  }
}














