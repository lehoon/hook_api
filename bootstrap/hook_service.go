package main

import (
	"encoding/json"
	"fmt"

	"github.com/lehoon/hook_api/v2/api"
	"github.com/lehoon/hook_api/v2/library/net/http"
	"github.com/lehoon/hook_api/v2/message"
)

type BusinessResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {

	for i := 0; i < 100; i++ {
		device1 := message.DeviceInfo{
			StreamId: "111",
			Username: "admin",
			Password: "Effort@2022",
			Hostname: "172.17.18.233",
			//VHostName: "__defaultVhost__",
			AppName: "jinan",
		}

		rsp, err := http.PostWithBody("http://localhost:9527/notify", device1, "application/json")

		if err != nil {
			fmt.Printf("%v", err)
		}

		if rsp != "" {
			fmt.Printf("%s", rsp)
		}

		fmt.Printf("\n")
	}
}

func test1() {
	json_data := "{\"MediaServiceId\":\"xxxxxx\",\"Data\":{\"Buffer\":0,\"BufferLikeString\":0,\"BufferList\":0,\"BufferRaw\":0,\"Frame\":0,\"FrameImp\":0,\"MediaSource\":120,\"MultiMediaSourceMuxer\":90,\"RtmpPacket\":80,\"RtpPacket\":70,\"Socket\":60,\"TcpClient\":50,\"TcpServer\":40,\"TcpSession\":30,\"UdpServer\":20,\"UdpSession\":10}}"
	var request message.KeepAliveReportRequest
	//buf, _ := json.Marshal(request)
	//fmt.Printf("%s\n", string(buf))

	json.Unmarshal([]byte(json_data), &request)
	fmt.Printf("%v\n", request)

	device := message.DeviceInfo{
		StreamId: "1000000004",
		Username: "admin",
		Password: "******",
		Hostname: "172.17.18.233",
	}

	deviceList := []message.DeviceInfo{}
	deviceList = append(deviceList, device)

	bizResult := BusinessResult{
		Code:    0,
		Message: "success",
		Data:    deviceList,
	}
	buf, _ := json.Marshal(bizResult)
	fmt.Printf("%s\n", string(buf))

	index := 0
	for {
		if index > 100 {
			break
		}

		id := 100000 + index
		streamid := fmt.Sprintf("%d", id)
		device := message.DeviceInfo{
			StreamId: streamid,
			Username: "admin",
			Password: "Effort@2022",
			Hostname: "172.17.18.233",
			//VHostName: "__defaultVhost__",
			AppName: "jinan",
		}

		//fmt.Printf("current device id is %d\n", id)
		fmt.Printf("%v\n", device)
		//url := fmt.Sprintf("http://localhost:8080/api/v1/device/%s", streamid)
		//http.DeleteUrl(url)
		rsp, err := http.PostWithBody("http://localhost:9000/api/v1/device/", device, "application/json")
		fmt.Printf("%s,%v\n", rsp, err)
		index++
	}

	//http.PostWithBody("http://localhost:8080/api/v1/device/", device, "application/json")
	//fmt.Printf("%s\n", string(buf))

	fmt.Printf("%v\n", api.OperateCodeMessage())
}
