package main

import (
	"flag"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/mesosphere/csilvm/pkg/csilvm"
)

const (
	defaultDefaultFs         = "xfs"
	defaultDefaultVolumeSize = 10 << 30
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	vgnameF := flag.String("volume-group", "", "The name of the volume group to manage")
	pvnamesF := flag.String("devices", "", "A comma-seperated list of devices in the volume group")
	defaultFsF := flag.String("default-fs", defaultDefaultFs, "The default filesystem to format new volumes with")
	defaultVolumeSizeF := flag.Uint64("default-volume-size", defaultDefaultVolumeSize, "The default volume size in bytes")
	socketFileF := flag.String("unix-addr", "", "The path to the listening unix socket file")
	removeF := flag.Bool("remove-volume-group", false, "If set, the volume group will be removed when ProbeNode is called.")
	profileF := flag.String("profile", "", "The volume group profile")
	flag.Parse()
	lis, err := net.Listen("unix", *socketFileF)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	var opts []csilvm.ServerOpt
	opts = append(opts, csilvm.DefaultVolumeSize(*defaultVolumeSizeF))
	if *removeF {
		opts = append(opts, csilvm.RemoveVolumeGroup())
	}
	if *profileF != "" {
		opts = append(opts, csilvm.Profile(*profileF))
	}
	s := csilvm.NewServer(*vgnameF, strings.Split(*pvnamesF, ","), *defaultFsF, opts...)
	csi.RegisterIdentityServer(grpcServer, s)
	csi.RegisterControllerServer(grpcServer, s)
	csi.RegisterNodeServer(grpcServer, s)
	grpcServer.Serve(lis)
}
