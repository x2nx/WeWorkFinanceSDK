package main

import (
	"bytes"
	"fmt"
	"github.com/x2nx/WeWorkFinanceSDK"
	"io/ioutil"
	"os"
	"path"
)

func main() {
	corpID := "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	corpSecret := "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	privateKeys := map[string]string{
		"版本号": "版本私钥",
	}

	//初始化客户端
	client, err := WeWorkFinanceSDK.NewClient(corpID, corpSecret, privateKeys)
	if err != nil {
		fmt.Printf("SDK 初始化失败：%v \n", err)
		return
	}

	//同步消息
	chatDataList, err := client.GetChatData(0, 100, "", "", 3)
	if err != nil {
		fmt.Printf("消息同步失败：%v \n", err)
		return
	}
	for _, chatData := range chatDataList {
		//消息解密
		chatInfo, err := client.DecryptData(chatData.EncryptRandomKey, chatData.EncryptChatMsg, fmt.Sprintf("%d", chatData.PublickeyVer))

		if err != nil {
			fmt.Printf("消息解密失败：%v \n", err)
			return
		}

		if chatInfo.Type == "image" {
			image := chatInfo.GetImageMessage()
			sdkfileid := image.Image.SdkFileID

			isFinish := false
			buffer := bytes.Buffer{}
			index_buf := ""
			for !isFinish {
				//获取媒体数据
				mediaData, err := client.GetMediaData(index_buf, sdkfileid, "", "", 5)
				if err != nil {
					fmt.Printf("媒体数据拉取失败：%v \n", err)
					return
				}
				buffer.Write(mediaData.Data)
				if mediaData.IsFinish {
					isFinish = mediaData.IsFinish
				}
				index_buf = mediaData.OutIndexBuf
			}
			filePath, _ := os.Getwd()
			filePath = path.Join(filePath, "test.png")
			err := ioutil.WriteFile(filePath, buffer.Bytes(), 0666)
			if err != nil {
				fmt.Printf("文件存储失败：%v \n", err)
				return
			}
			break
		}
	}
}
