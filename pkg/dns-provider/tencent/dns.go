package tencent

import (
	"context"
	"fmt"

	common "github.com/EdisonLai/ddns/pkg/dns-provider/dns-common"
	"github.com/EdisonLai/ddns/pkg/logger"

	tencommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

func NewDNSClient(secretID, secretKey string) (client *dnspod.Client) {
	// 实例化一个认证对象，入参需要传入腾讯云账户 SecretId 和 SecretKey，此处还需注意密钥对的保密
	// 代码泄露可能会导致 SecretId 和 SecretKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考，建议采用更安全的方式来使用密钥，请参见：https://cloud.tencent.com/document/product/1278/85305
	// 密钥可前往官网控制台 https://console.cloud.tencent.com/cam/capi 进行获取
	credential := tencommon.NewCredential(
		secretID,
		secretKey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ = dnspod.NewClient(credential, "", cpf)
	return
}

func DescribeDomain(ctx context.Context, client *dnspod.Client, domain string, subdomainfilter []string) (records []common.DNSRecord, err error) {
	logCtx := logger.GetEntry(ctx)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := dnspod.NewDescribeRecordListRequest()

	request.Domain = tencommon.StringPtr(domain)
	if len(subdomainfilter) == 1 {
		request.Subdomain = tencommon.StringPtr(subdomainfilter[0])
	}

	// 返回的resp是一个DescribeRecordListResponse的实例，与请求对象对应
	response, err := client.DescribeRecordList(request)
	if _, ok := err.(*errors.TencentCloudSDKError); err != nil || ok {
		fmt.Printf("An API error has returned: %s", err)
		return nil, err
	}

	fmt.Printf("TencentSDKResp: %s", response.ToJsonString())

	recordMap := make(map[string]common.DNSRecord)
	for _, v := range response.Response.RecordList {
		recordMap[*v.Name] = common.DNSRecord{
			Type:      *v.Type,
			RecordId:  *v.RecordId,
			SubDomain: *v.Name,
			Value:     *v.Value,
			Status:    *v.Status,
			Line:      *v.Line,
			LineId:    *v.LineId,

			Weight: v.Weight,
			MX:     v.MX,
			Remark: v.Remark,
			TTL:    v.TTL,
		}
	}

	for _, v := range subdomainfilter {
		records = append(records, recordMap[v])
	}
	if len(subdomainfilter) == 0 {
		for _, v := range recordMap {
			records = append(records, v)
		}
	}
	logCtx.Infof("provider record[%v]", records)
	return
}

func ModifyDynamicDNS(ctx context.Context, client *dnspod.Client, recordId uint64, domain, subdomain, lineZh, value string) (err error) {
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := dnspod.NewModifyDynamicDNSRequest()

	request.Domain = tencommon.StringPtr(domain)
	request.SubDomain = tencommon.StringPtr(subdomain)
	request.RecordId = tencommon.Uint64Ptr(recordId)
	request.RecordLine = tencommon.StringPtr(lineZh)

	request.Value = tencommon.StringPtr(value)

	// 返回的resp是一个ModifyDynamicDNSResponse的实例，与请求对象对应
	response, err := client.ModifyDynamicDNS(request)
	if _, ok := err.(*errors.TencentCloudSDKError); err != nil || ok {
		fmt.Printf("An API error has returned: %s", err)
		return err
	}

	fmt.Printf("TencentSDKResp: %s", response.ToJsonString())

	return
}
