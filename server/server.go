package server

import (
	"github.com/astaxie/beego"
	"net"
	"strings"
	"github.com/jxufeliujj/tcp-server/core"
	"github.com/jxufeliujj/tcp-server/socket"
	"time"
)

func StartTCP() {
	go start_843Port()
	tcpAddr, errAddr := net.ResolveTCPAddr("tcp4", beego.AppConfig.String("bind"))
	if errAddr != nil {
		panic(errAddr.Error())
	}
	server, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err.Error())
	} else {
		core.Writelog(core.LevelInfo, "fun:main.startTCP", "TCP Server on:", tcpAddr.String(), "-", "-")
	}
	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			core.Writelog(core.LevelError, "fun:main.startTCP", "client connect error", err.Error(), "-", "-")
			continue
		}
		core.Writelog(core.LevelInfo, "fun:main.startTCP", "client connect:", conn.RemoteAddr().String(), "-", "-")
		go socket.Start(conn)
	}
}

func start_843Port() {
	addr_port := strings.Split(beego.AppConfig.String("bind"), ":")
	addr843 := addr_port[0] + ":843"

	tcpAddr, errAddr := net.ResolveTCPAddr("tcp4", addr843)
	if errAddr != nil {
		panic(errAddr.Error())
	}
	listen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err.Error())
	} else {
		core.Writelog(core.LevelInfo, "fun:main.start_843Port", "TCP Server on:", tcpAddr.String(), "-", "-")
	}
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			core.Writelog(core.LevelError, "fun:main.start_843Port", "client connect error", err.Error(), "-", "-")
			continue
		}
		core.Writelog(core.LevelInfo, "fun:main.start_843Port", "client connect:", conn.RemoteAddr().String(), "-", "-")
		go sendFirstMsg(conn)
	}
}

func sendFirstMsg(conn *net.TCPConn) {
	xml := `<?xml version="1.0"?>
			<!DOCTYPE cross-domain-policy SYSTEM "/xml/dtds/cross-domain-policy.dtd">
			<cross-domain-policy>
				<site-control permitted-cross-domain-policies="master-only"/>
				<allow-access-from domain="*" to-ports="*" />
			</cross-domain-policy>`
	conn.Write([]byte(xml))
	time.Sleep(time.Second)
	conn.Close()
	core.Writelog(core.LevelInfo, "fun:main.sendFirstMsg", "已经回应策略文件：crossdomain.xml", conn.RemoteAddr().String(), "-", "-")
}
