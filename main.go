package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	goku_plugin "github.com/eolinker/goku-plugin"
)

const pluginName = "goku-ip_restriction"

var _ goku_plugin.PluginBeforeMatch = (*gokuIp)(nil)
var _ goku_plugin.PluginAccess = (*gokuIp)(nil)

var builder = new(gokuIpPluginFactory)

func Builder() goku_plugin.PluginFactory {
	return builder
}

type IPList struct {
	IpListType  string   `json:"ipListType"`
	IpWhiteList []string `json:"ipWhiteList"`
	IpBlackList []string `json:"ipBlackList"`
}
type gokuIpPluginFactory struct {
}

func ip2binary(ip string) string {
	str := strings.Split(ip, ".")
	var ipstr string
	for _, s := range str {
		i, _ := strconv.ParseUint(s, 10, 8)

		ipstr = ipstr + fmt.Sprintf("%08b", i)
	}
	return ipstr
}

func match(ip, iprange string) bool {
	ipb := ip2binary(ip)
	ipr := strings.Split(iprange, "/")
	if len(ipr) < 2 {
		return ip == ipr[0]
	}
	masklen, err := strconv.ParseUint(ipr[1], 10, 32)
	if err != nil {

		return false
	}
	iprb := ip2binary(ipr[0])
	return strings.EqualFold(ipb[0:masklen], iprb[0:masklen])
}

func (f *gokuIpPluginFactory) Create(config string, clusterName string, updateTag string, strategyId string, apiId int) (*goku_plugin.PluginObj, error) {
	if config == "" {
		return nil, errors.New("config is empty")
	}
	var ipList IPList

	err := json.Unmarshal([]byte(config), &ipList)
	if err != nil {

		return nil, err
	}

	p := &gokuIp{
		ipList: &ipList,
	}

	return &goku_plugin.PluginObj{
		BeforeMatch: p,
		Access:      p,
		Proxy:       nil,
	}, nil

}

type gokuIp struct {
	ipList *IPList
}

// 匹配URI前执行
func (this *gokuIp) BeforeMatch(ctx goku_plugin.ContextBeforeMatch) (bool, error) {

	ipList := this.ipList
	remoteAddr := ctx.Request().RemoteAddr()
	if realIP, ok := ctx.Request().Headers()["X-Real-Ip"]; ok {
		remoteAddr = strings.Join(realIP, ",")
	}

	if ipList.IpListType == "white" {
		flag := false
		for _, v := range ipList.IpWhiteList {
			if v != "" {
				if match(remoteAddr, v) {
					flag = true
					break
				}
			}
		}
		if !flag {
			ctx.SetStatus(403, "403")
			ctx.SetBody([]byte("[ip_restriction] Illegal IP!"))
			return false, errors.New("[ip_restriction] Illegal IP!")
		}
	} else if ipList.IpListType == "black" {
		for _, v := range ipList.IpBlackList {
			if v != "" {
				if match(remoteAddr, v) {

					ctx.SetStatus(403, "403")
					ctx.SetBody([]byte("[ip_restriction] Illegal IP!"))
					return false, errors.New("[ip_restriction] Illegal IP!")
				}
			}
		}
	}

	return true, nil
}

func convertIP(ip string) (error, string) {
	ipr := strings.Split(ip, "/")
	errInfo := "[ip_restriction] Illegal ip:" + ip
	if len(ipr) > 0 {
		ips := strings.Split(ipr[0], ".")
		ipLen := len(ips)
		if firstIndex := strings.Index(ipr[0], "*"); firstIndex > 0 {
			if lastIndex := strings.LastIndex(ipr[0], "*"); firstIndex == lastIndex && ips[ipLen-1] == "*" {
				v := ""
				for i := 0; i < 4; i++ {
					if i < ipLen-1 {
						v += ips[i] + "."
					} else {
						v += "0"
						if i != 3 {
							v += "."
						}
					}
				}
				v += "/" + strconv.Itoa((ipLen-1)*8)
				return nil, v
			} else {
				return errors.New(errInfo), ""
			}
		} else {
			if ipLen < 4 {
				return errors.New(errInfo), ""
			}
			return nil, ip
		}
	} else {
		return errors.New(errInfo), ""
	}
}

// 转发前执行
func (this *gokuIp) Access(ctx goku_plugin.ContextAccess) (bool, error) {

	ipList := this.ipList
	var err error
	remoteAddr := ctx.Request().RemoteAddr()
	if realIP, ok := ctx.Request().Headers()["X-Real-Ip"]; ok {
		remoteAddr = strings.Join(realIP, ",")
	}
	//
	ips := strings.Split(remoteAddr, ":")
	ip := ""
	if len(ips) > 0 {
		ip = ips[0]
	}
	if ipList.IpListType == "white" {
		flag := false
		for _, v := range ipList.IpWhiteList {
			if v == "*" {
				flag = true
				break
			}
			err, v = convertIP(v)
			if err != nil {
				ctx.SetStatus(403, "403")
				ctx.SetBody([]byte(err.Error()))
				return false, err
			}
			if v != "" {
				if match(ip, v) {
					flag = true
					break
				}
			}
		}
		if !flag {
			ctx.SetStatus(403, "403")
			ctx.SetBody([]byte("[ip_restriction] Illegal IP!"))
			return false, errors.New("[ip_restriction] Illegal IP!")
		}
	} else if ipList.IpListType == "black" {
		for _, v := range ipList.IpBlackList {
			if v == "*" {
				ctx.SetStatus(403, "403")
				ctx.SetBody([]byte("[ip_restriction] Illegal IP!"))
				return false, errors.New("[ip_restriction] Illegal IP!")
			}
			err, v = convertIP(v)
			if err != nil {
				ctx.SetStatus(403, "403")
				ctx.SetBody([]byte(err.Error()))
				return false, err
			}
			if v != "" {
				if match(ip, v) {
					ctx.SetStatus(403, "403")
					ctx.SetBody([]byte("[ip_restriction] Illegal IP!"))
					return false, errors.New("[ip_restriction] Illegal IP!")
				}
			}
		}
	}
	return true, nil
}
