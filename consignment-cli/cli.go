package main

import (
  pb "shippy/consignment-service/proto/consignment"
  "io/ioutil"
  "encoding/json"
  "errors"
  "github.com/micro/go-micro"
  "os"
  "log"
  "context"
)

const (
  ADDRESS               =   "localhost:50051"
  DEFAULT_INFO_FILE     =    "consignment.json"
)

func parseFile(fileName string) (*pb.Consignment, error)  {
  data, err := ioutil.ReadFile(fileName)
  if err != nil {
    return nil, err
  }

  var consignment *pb.Consignment
  err = json.Unmarshal(data, &consignment)
  if err != nil {
    return nil, errors.New("consignment.json file content error")
  }
  return consignment, nil
}

func main()  {
  service := micro.NewService(micro.Name("go.micro.srv.consignment"))
  service.Init()

  client := pb.NewShippingServiceClient("go.micro.srv.consignment", service.Client())

  infoFile := DEFAULT_INFO_FILE
  if len(os.Args) > 1 {
    infoFile = os.Args[1]
  }

  consignment, err := parseFile(infoFile)
  if err != nil {
    log.Fatalf("create consignment error: %v", err)
  }

  resp, err := client.CreateConsignment(context.Background(), consignment)
  if err != nil {
    log.Fatalf("create consignments: %v ", err)
  }

  for _, c := range resp.Consignments {
    log.Printf("%+v", c)
  }
}








