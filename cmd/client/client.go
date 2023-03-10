package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/EdisonLai/ddns/cmd/client/config"
	dnscommon "github.com/EdisonLai/ddns/pkg/dns-provider/dns-common"
	"github.com/EdisonLai/ddns/pkg/dns-provider/tencent"
	"github.com/EdisonLai/ddns/pkg/logger"
)

var record dnscommon.DNSRecord

func main() {
	var err error

	if err = config.InitConfig(); err != nil {
		panic(err)
	}

	ctx := logger.InitLogger(context.Background(), "")
	logCtx := logger.GetEntry(ctx)

	client := tencent.NewDNSClient(config.GConf.Provider.SecretId, config.GConf.Provider.SecretKey)

	var records []dnscommon.DNSRecord
	if records, err = tencent.DescribeDomain(ctx, client, config.GConf.Domain.Domain, []string{config.GConf.Domain.SubDomain}); err != nil {
		logCtx.WithError(err).Debugf("describe domain fail")
		os.Exit(-2)
		return
	}
	if len(records) == 1 {
		record = records[0]
		logCtx.Infof("record found[%+v]", record)
	} else {
		os.Exit(-2)
		return
	}

	newIPCh := make(chan string)
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				newIP, err := GetPublicIP(ctx, config.GConf.EIPMethod.Server)
				if err != nil {
					logCtx.WithError(err).Errorf("get public ip fail")
					close(newIPCh)
					return
				}

				newIPCh <- newIP
				time.Sleep(time.Duration(config.GConf.Domain.CheckTime) * time.Minute)
			}
		}
	}(ctx)

	for {
		select {
		case newIP, ok := <-newIPCh:
			if !ok {
				logCtx.Error("backend error happend")
				os.Exit(-2)
				return
			}

			if newIP != record.Value {
				if err = tencent.ModifyDynamicDNS(ctx, client, record.RecordId, config.GConf.Domain.Domain, config.GConf.Domain.SubDomain, record.Line, newIP); err != nil {
					logCtx.WithError(err).Error("backend error happend")
					os.Exit(0)
					return
				}
			}

		case <-ctx.Done():
			logCtx.Info("exit safely")
			os.Exit(0)
			return
		}
	}
}

func GetPublicIP(ctx context.Context, server string) (ip string, err error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/client_ip", server), nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

	return string(body), nil
}
