package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	pb "github.com/synerex/synerex_api"
	sxutil "github.com/synerex/synerex_sxutil"
)

var (
	srcSrv             = flag.String("srcsrv", "127.0.0.1:9990", "Source Synerex Node ID Server")
	dstSrv             = flag.String("destsrv", "127.0.0.1:9990", "Destination Synerex Node ID Server")
	channel            = flag.Int("channel", 3, "Forwarding channel type")
	mu                 sync.Mutex
	sxSrcServerAddress string
	sxDstServerAddress string
	sxDstClient        *sxutil.SXServiceClient
	msgCount           int64
)

func init() {
	msgCount = 0
}

// callback for each Supply
func supplyCallback(clt *sxutil.SXServiceClient, sm *pb.Supply) {

	msgCount++
	cont := pb.Content{Entity: sm.Cdata.Entity}
	smo := sxutil.SupplyOpts{
		Name:  sm.SupplyName,
		Cdata: &cont,
	}
	//			fmt.Printf("Res: %v",smo)
	_, nerr := sxDstClient.NotifySupply(&smo)
	if nerr != nil {
		log.Printf("Error on sent %v", nerr)
	}

}

func subscribeSupply(client *sxutil.SXServiceClient) {
	// goroutine!
	ctx := context.Background() //
	client.SubscribeSupply(ctx, supplyCallback)
	// comes here if channel closed
	log.Printf("Server closed... on Forward provider")
}

// just for stat
func monitorStatus() {
	for {
		sxutil.SetNodeStatus(int32(msgCount), fmt.Sprintf("dt:%d", msgCount))
		time.Sleep(time.Second * 3)
	}
}

func main() {
	flag.Parse()
	if *srcSrv == *dstSrv {
		log.Fatal("Input servers should not be same address")
	}

	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	channelTypes := []uint32{uint32(*channel)}
	// obtain synerex server address from nodeserv
	srcSSrv, err := sxutil.RegisterNode(*srcSrv, fmt.Sprintf("FowardSink[%d]", *channel), channelTypes, nil)
	if err != nil {
		log.Fatal("Can't register to source node...")
	}
	log.Printf("Connecting Source Server [%s]\n", srcSSrv)
	sxSrcServerAddress = srcSSrv

	dstSSrv, derr := sxutil.RegisterNode(*dstSrv, fmt.Sprintf("FowardSource[%d]", *channel), channelTypes, nil)
	if derr != nil {
		log.Fatal("Can't register to destination node...")
	}
	log.Printf("Connecting Destination Server [%s]\n", dstSSrv)
	sxDstServerAddress = dstSSrv

	wg := sync.WaitGroup{} // for syncing other goroutines
	srcClient := sxutil.GrpcConnectServer(sxSrcServerAddress)
	argJson := fmt.Sprintf("{ForwardSink[%d]}", *channel)
	sxSrclient := sxutil.NewSXServiceClient(srcClient, uint32(*channel), argJson)

	dstClient := sxutil.GrpcConnectServer(sxDstServerAddress)
	argDstJson := fmt.Sprintf("{ForwardSrc[%d]}", *channel)
	sxDstClient = sxutil.NewSXServiceClient(dstClient, uint32(*channel), argDstJson)

	wg.Add(1)

	go subscribeSupply(sxSrclient)

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
