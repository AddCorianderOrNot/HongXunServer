package utils

import (
	"HongXunServer/auth"
	"fmt"
	cms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cms/v20190321"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"log"
	"strings"
	"unicode/utf8"
)

func Moderation(text *string) {

	credential := common.NewCredential(
		auth.SecretId,
		auth.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cms.tencentcloudapi.com"
	client, _ := cms.NewClient(credential, "ap-guangzhou", cpf)

	request := cms.NewTextModerationRequest()

	request.Content = common.StringPtr(Base64([]byte(*text)))

	response, err := client.TextModeration(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%s", response.ToJsonString())
	if *response.Response.Data.Suggestion == "Block" {
		for _, keyword := range response.Response.Data.Keywords {
			*text = strings.Replace(*text, *keyword, strings.Repeat("*", utf8.RuneCountInString(*keyword)), -1)
		}
	}
	log.Println(*text)
}
