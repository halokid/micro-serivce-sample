package main

import (
  "context"
  "encoding/json"
  "github.com/pkg/errors"
  "log"
  "google.golang.org/grpc"
  "io/ioutil"
  "os"
  pb "shippy/consignment-service/proto/consignment"
)

const (
  ADDRESS             =   "localhost:50051"
  DEFAULT_INFO_FILE   =   "consignment.json"
)

func parserFile(filename string) (*pb.Consignment, error) {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }

  var consignment *pb.Consignment
  err = json.Unmarshal(data, &consignment)
  if err != nil {
    return nil, errors.New("consignment,json file content error")
  }
  return consignment, nil
}

func main() {
  conn, err := grpc.Dial(ADDRESS, grpc.WithInsecure())
  if err != nil {
    log.Fatal("conect error: v%", err)
  }
  defer conn.Close()

  client := pb.NewShippingServiceClient(conn)

  infoFile := DEFAULT_INFO_FILE
  if len(os.Args) > 1 {
    infoFile = os.Args[1]
  }

  consignment, err := parserFile(infoFile)
  if err != nil {
    log.Fatalf("parser info file error: %v", err)
  }

  resp, err := client.CreateConsignment(context.Background(), consignment)
  if err != nil {
    log.Fatalf("create consignment error: %v", err)
  }

  log.Printf("created: %t", resp.Created)

  resp, err = client.GetConsignments(context.Background(), &pb.GetRequest{})
  if err != nil {
    log.Fatalf("failed to list consignment: %v", err)
  }

  for _, c := range resp.Consignments {
    log.Printf("%+v", c)
  }
}








