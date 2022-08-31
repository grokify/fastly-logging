package main

import (
	"log"
	"strings"

	fastlylogging "github.com/grokify/fastly-logging"
	"github.com/grokify/mogo/fmt/fmtutil"
)

const LoggingFieldsS3 = "timestamp,client_ip,geo_country,geo_city,host,url,request_method,request_protocol,request_referer,request_user_agent,response_state,response_status,response_reason,response_body_size,fastly_server,fastly_is_edge"

func main() {
	props := strings.Split(LoggingFieldsS3, ",")
	format, err := fastlylogging.EntryFormatDefault(props...)
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.MustPrintJSON(format)
}

/*

Reference from Fastly Docs:

https://docs.fastly.com/en/guides/log-streaming-amazon-s3

{
	"timestamp": "%{strftime(\{"%Y-%m-%dT%H:%M:%S%z"\}, time.start)}V",
	"client_ip": "%{req.http.Fastly-Client-IP}V",
	"geo_country": "%{client.geo.country_name}V",
	"geo_city": "%{client.geo.city}V",
	"host": "%{if(req.http.Fastly-Orig-Host, req.http.Fastly-Orig-Host, req.http.Host)}V",
	"url": "%{json.escape(req.url)}V",
	"request_method": "%{json.escape(req.method)}V",
	"request_protocol": "%{json.escape(req.proto)}V",
	"request_referer": "%{json.escape(req.http.referer)}V",
	"request_user_agent": "%{json.escape(req.http.User-Agent)}V",
	"response_state": "%{json.escape(fastly_info.state)}V",
	"response_status": %{resp.status}V,
	"response_reason": %{if(resp.response, "%22"+json.escape(resp.response)+"%22", "null")}V,
	"response_body_size": %{resp.body_bytes_written}V,
	"fastly_server": "%{json.escape(server.identity)}V",
	"fastly_is_edge": %{if(fastly.ff.visits_this_service == 0, "true", "false")}V
  }
*/
