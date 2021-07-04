package main

/* Forwarding Provider
   (Currently only supports NotifySupply)
*/

import (
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
	srcSxAddr          = flag.String("srcsxsrv", "", "Source Synerex Server Addr")
	dstSrv             = flag.String("dstsrv", "127.0.0.1:9990", "Destination Synerex Node ID Server")
	dstSxAddr          = flag.String("dstsxsrv", "", "Destination Synerex Server Addr")
	channel            = flag.Int("channel", 3, "Forwarding channel type")
	mu                 sync.Mutex
	sxSrcServerAddress string
	sxDstServerAddress string
	sxDstClient        *sxutil.SXServiceClient
	msgCount           int64
)

const SEND_RETRY_COUNT = 5

func init() {
	msgCount = 0
}

// callback for each Supply
func supplyCallback(clt *sxutil.SXServiceClient, sm *pb.Supply) {

	msgCount++
	var smo sxutil.SupplyOpts
	if sm.Cdata != nil {
		cont := pb.Content{Entity: sm.Cdata.Entity}
		smo = sxutil.SupplyOpts{
			Name:  sm.SupplyName,
			Cdata: &cont,
			JSON:  sm.ArgJson,
		}
		//			fmt.Printf("Res: %v",smo)
	}else{
		smo = sxutil.SupplyOpts{
			Name:  sm.SupplyName,
			Cdata: nil,
			JSON:  sm.ArgJson,
		}
	}
	_, nerr := sxDstClient.NotifySupply(&smo)
	if nerr != nil {
		ecount := 1
		for ecount < SEND_RETRY_COUNT && nerr != nil {
			log.Printf("Error on sent %v", nerr)
			
			time.Sleep(5 * time.Second * time.Duration(ecount))
			// we need to reconecct dst server
			log.Printf("Reconnect Dst server [%s]", sxDstServerAddress)
		/* sxutil may reconnect !? */
		//		dstClient := sxutil.GrpcConnectServer(sxDstServerAddress)
		//		argDstJson := fmt.Sprintf("{ForwardSrc[%d]}", *channel)
		//		sxDstClient = dstNI.NewSXServiceClient(dstClient, uint32(*channel), argDstJson)
			_, nerr = sxDstClient.NotifySupply(&smo)
			ecount++
		}
		if ecount == SEND_RETRY_COUNT {
			log.Printf("Error count exceeded! reason: %v", nerr)			
		}
	}

}
/*
func subscribeSupply(client *sxutil.SXServiceClient) {
	// goroutine for Src Server.
	for {
		
		ctx := context.Background() //
		log.Printf("SubscirbeSupply with %v", client)
		serr := client.SubscribeSupply(ctx, supplyCallback)
		// comes here if channel closed
		log.Print("Server closed... on Src Forward provider from:", sxSrcServerAddress, ",error:", serr)

		time.Sleep(5 * time.Second)
		//TODO: should check nodeserver.
		newClt := sxutil.GrpcConnectServer(sxSrcServerAddress)
		if newClt != nil {
			log.Printf("Reconnect Src server [%s]", sxSrcServerAddress)
			client.SXClient = newClt
		} else {
			log.Printf("Connection Error!! on Src Server")
		}
	}
 }
*/

// just for stat
func monitorStatus() {
	for {
		sxutil.SetNodeStatus(int32(msgCount), fmt.Sprintf("recv:%d", msgCount))
		time.Sleep(time.Second * 3)
	}
}

func monitorStatusDst(dstNI *sxutil.NodeServInfo) {
	for {
		dstNI.SetNodeStatus(int32(msgCount), fmt.Sprintf("sent:%d", msgCount))
		time.Sleep(time.Second * 3)
	}
}

func main() {
	log.Printf("FowardProvider(%s) built %s sha1 %s", sxutil.GitVer, sxutil.BuildTime, sxutil.Sha1Ver)

	flag.Parse()
	if *srcSrv == *dstSrv {
		log.Fatal("Input servers should not be same address")
	}

	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	dstNI := sxutil.NewNodeServInfo() // for dst node server connection
	sxutil.RegisterDeferFunction(dstNI.UnRegisterNode)

	channelTypes := []uint32{uint32(*channel)}
	// obtain synerex server address from nodeserv
	srcSSrv, err := sxutil.RegisterNode(*srcSrv, fmt.Sprintf("FowardSrc[%d]", *channel), channelTypes, nil)
	if err != nil {
		log.Fatal("Can't register to source node...")
	}
	if *srcSxAddr != "" {
		srcSSrv = *srcSxAddr
	}

	log.Printf("Connecting Source Server [%s]\n", srcSSrv)
	sxSrcServerAddress = srcSSrv

	dstSSrv, derr := dstNI.RegisterNodeWithCmd(*dstSrv, fmt.Sprintf("FowardDst[%d]", *channel), channelTypes, nil, nil)
	if derr != nil {
		log.Fatal("Can't register to destination node...")
	}
	if *dstSxAddr != "" {
		dstSSrv = *dstSxAddr
	}

	log.Printf("Connecting Destination Server [%s]\n", dstSSrv)
	sxDstServerAddress = dstSSrv

	wg := sync.WaitGroup{} // for syncing other goroutines
	srcClient := sxutil.GrpcConnectServer(sxSrcServerAddress)
	if srcClient == nil {
		log.Fatal("Can't connect source Synerex server")
	}
	argJson := fmt.Sprintf("{ForwardSink[%d]}", *channel)
	sxSrcClient := sxutil.NewSXServiceClient(srcClient, uint32(*channel), argJson)

	dstClient := sxutil.GrpcConnectServer(sxDstServerAddress)
	if dstClient == nil {
		log.Fatal("Can't connect destination Synerex server")
	}
	argDstJson := fmt.Sprintf("{ForwardSrc[%d]}", *channel)
	sxDstClient = dstNI.NewSXServiceClient(dstClient, uint32(*channel), argDstJson)

	wg.Add(1)

	sxutil.SimpleSubscribeSupply(sxSrcClient,supplyCallback)
	go monitorStatus()
	go monitorStatusDst(dstNI)

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
