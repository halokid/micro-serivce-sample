package main

import (
  "context"
  "gopkg.in/mgo.v2"
  "log"
  vesselPb "shippy/vessel-service/proto/vessel"
  pb "shippy/consignment-service/proto/consignment"
)

// 微服务服务端 struct handler 必须实现 protobuf 中定义的 rpc 方法
// 实现方法的传参等可参考生成的 consignment.pb.go
type handler struct {
  session *mgo.Session
  vesselClient vesselPb.VesselServiceClient
}

// 从主会话中 Clone() 出新会话处理查询
// todo: golang的一个牛逼之处就在于，啥程度逻辑都可以返回 interfaces{} 类型，比如这个
// todo: Repository 的三个方法是用来执行SQL查询的， 这个函数表示在执行sql之前， 先clone了一个mgo的会话，所以是先return的， clone的这个会话是供执行sql用的
func (h *handler) GetRepo() Repository  {
  // 执行了 Repository 的方法之后， 然后就返回下面的逻辑， 也就是关闭数据库
  // 所以目的就是 执行sql之后， 然后就关闭数据库
  return &ConsignmentRepository{h.session.Clone()}
}

func (h *handler) CreateConsignment (ctx context.Context, req *pb.Consignment, resp *pb.Response) error {
  defer h.GetRepo().Close()

  vReq := &vesselPb.Specification{
    Capacity:   int32(len(req.Containers)),
    MaxWeight:  req.Weight,
  }

  vResp, err := h.vesselClient.FindAvailable(context.Background(), vReq)
  if err != nil {
    return err
  }

  log.Printf("found vessel:   %s\n", vResp.Vessel.Name)
  req.VesselId = vResp.Vessel.Id
  err = h.GetRepo().Create(req)
  if err != nil {
    return err
  }
  resp.Created = true
  resp.Consignment = req
  return nil
}

func (h *handler)GetConsignments(ctx context.Context, req *pb.GetRequest, resp *pb.Response) error {
  defer h.GetRepo().Close()
  consignments, err := h.GetRepo().GetAll()
  if err != nil {
    return err
  }
  resp.Consignments = consignments
  return nil
}




