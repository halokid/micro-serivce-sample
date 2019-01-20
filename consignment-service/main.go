package main

import (
  "context"
  "github.com/micro/go-micro"
  "log"
  pb "shippy/consignment-service/proto/consignment"
  vesselPb "shippy/vessel-service/proto/vessel"
)


// ----------------------------- 定义接口和数据结构 ----------------
// 定义接口定义， 负责方法梳理
type IRepository interface {
  Create(consignment *pb.Consignment) (*pb.Consignment, error)
  GetAll() []*pb.Consignment
}


// 定义结构体，负责数据梳理， 实现了接口
type Repository struct {
  consignments []*pb.Consignment
}

// 以数据为入口去实现方法， 编程清晰逻辑，结构化的套路
func (repo *Repository) Create (consignment *pb.Consignment) (*pb.Consignment, error) {
  repo.consignments = append(repo.consignments, consignment)
  return consignment, nil
}

func (repo *Repository) GetAll() ([]*pb.Consignment)  {
  return repo.consignments
}



// ---------------------------- 定义微服务 ---------------------
type service struct {
  repo Repository
  vesselClient vesselPb.VesselServiceClient
}


//
// 实现 consignment.pb.go 中的 ShippingServiceHandler 接口
// 使 service 作为 gRPC 的服务端
//
// 托运新的货物
// func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, resp *pb.Response) error {

  // 检查是否有适合的货轮
  vReq := &vesselPb.Specification{
    Capacity:  int32(len(req.Containers)),
    MaxWeight: req.Weight,
  }
  vResp, err := s.vesselClient.FindAvailable(context.Background(), vReq)
  if err != nil {
    return err
  }

  // 货物被承运
  log.Printf("found vessel: %s\n", vResp.Vessel.Name)
  req.VesselId = vResp.Vessel.Id
  consignment, err := s.repo.Create(req)
  if err != nil {
    return err
  }
  resp.Created = true
  resp.Consignment = consignment
  return nil
}

// 获取目前所有托运的货物
// func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, resp *pb.Response) error {
  allConsignments := s.repo.GetAll()
  resp = &pb.Response{Consignments: allConsignments}
  return nil
}

func main() {
  server := micro.NewService(
    // 必须和 consignment.proto 中的 package 一致
    micro.Name("go.micro.srv.consignment"),
    micro.Version("latest"),
  )

  // 解析命令行参数
  server.Init()
  repo := Repository{}
  // 作为 vessel-service 的客户端
  vClient := vesselPb.NewVesselServiceClient("go.micro.srv.vessel", server.Client())
  pb.RegisterShippingServiceHandler(server.Server(), &service{repo, vClient})

  if err := server.Run(); err != nil {
    log.Fatalf("failed to serve: %v", err)
  }
}

