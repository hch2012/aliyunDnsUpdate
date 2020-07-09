package main

import(
	"github.com/denverdino/aliyungo/dns"
	"encoding/json"
	"net/http"
	"log"
	"io/ioutil"
	"time"
)
const(
	accessKeyId     = "your accessKeyId"
	accessKeySecret = "your accessKeySecret"
	domainName      = "your domainName"
)
var client *dns.Client
func init() {
	client=dns.NewClient(accessKeyId, accessKeySecret)
}
func catch(fun func()){
    defer func() {
        if r := recover(); r != nil {
        	log.Println(r)
        	time.Sleep(5*time.Second)
        	catch(fun)
        }
    }()
    fun()
}
func main() {
	catch(run)
}
func run(){
	describeArgs:= dns.DescribeDomainRecordsArgs{
		DomainName: domainName,
	}
	var ip string;
	for{
		ip=getLocalIp()
		log.Println(ip)
		descResponse,err:=client.DescribeDomainRecords(&describeArgs)
		checkErr(err)
		for _,descRecord:=range descResponse.DomainRecords.Record{
			if(ip!=descRecord.Value && descRecord.Type=="A"){
				update(descRecord.RecordId,descRecord.RR,ip)
			}
		}
		time.Sleep(5*time.Minute)
	}
}


type IpInfo struct{
	Cip string
	Cid string
	Cname string
}

func getLocalIp() string{
	resp, err := http.Get("http://pv.sohu.com/cityjson")
	checkErr(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	body=body[19:len(body)-1];
	var ipInfo IpInfo
	json.Unmarshal(body,&ipInfo)
	return ipInfo.Cip
}

func update(recordId string,rr string,ip string) {
	updateArgs := dns.UpdateDomainRecordArgs{
		RecordId: recordId,
		RR:       rr,
		Value:    ip,
		Type:     "A",
	}
	log.Println(updateArgs)
	_,err:=client.UpdateDomainRecord(&updateArgs);
	checkErr(err)
}

func checkErr(err error){
	if(err!=nil){
		panic(err)
	}
}
