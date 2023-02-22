package service

import "github.com/lehoon/hook_api/v2/message"

func QueryStreamByStreamId(streamId string) (*message.StreamInfo, error) {
	deviceInfo, err := QueryDeviceByDeviceId(streamId)

	if err != nil {
		return nil, err
	}

	return &message.StreamInfo{
		StreamId:  deviceInfo.DeviceId,
		Username:  deviceInfo.Username,
		Password:  deviceInfo.Password,
		Hostname:  deviceInfo.Hostname,
		AppName:   deviceInfo.AppName,
		VHostName: deviceInfo.VHostName,
	}, nil
}

func QueryStreamList() ([]message.StreamInfo, error) {
	resultList := []message.StreamInfo{}
	deviceList, err := QueryDeviceList()

	if err != nil {
		return resultList, nil
	}

	for _, device := range deviceList {
		resultList = append(resultList,
			message.StreamInfo{
				StreamId:  device.DeviceId,
				Username:  device.Username,
				Password:  device.Password,
				Hostname:  device.Hostname,
				AppName:   device.AppName,
				VHostName: device.VHostName,
			},
		)
	}

	return resultList, nil
}
