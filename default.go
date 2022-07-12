package fastlylogging

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/grokify/mogo/type/stringsutil"
)

const (
	FieldTimestamp                     = "timestamp"
	FieldTimestampDefault              = `%{strftime(\{"%Y-%m-%dT%H:%M:%S%z"\}, time.start)}V`
	FieldServiceID                     = "service_id"
	FieldServiceIDDefault              = `%{req.service_id}V`
	FieldTimeElapsed                   = "time_elapsed"
	FieldTimeElapsedDefault            = `%{time.elapsed.usec}V`
	FieldClientIP                      = "client_ip"
	FieldClientIPDefault               = `%{req.http.Fastly-Client-IP}V`
	FieldGeoCountry                    = "geo_country"
	FieldGeoCountryDefault             = `%{client.geo.country_name}V`
	FieldGeoCity                       = "geo_city"
	FieldGeoCityDefault                = `%{client.geo.city}V`
	FieldHost                          = "host"
	FieldHostDefault                   = `%{if(req.http.Fastly-Orig-Host, req.http.Fastly-Orig-Host, req.http.Host)}V`
	FieldURL                           = "url"
	FieldURLDefault                    = `%{json.escape(req.url)}V`
	FieldRequest                       = "request"
	FieldRequestDefault                = `%{req.request}V`
	FieldRequestMethod                 = "request_method"
	FieldRequestMethodDefault          = `%{json.escape(req.method)}V`
	FieldRequestProtocol               = "request_protocol"
	FieldRequestProtocolDefault        = `%{json.escape(req.proto)}V`
	FieldRequestReferer                = "request_referer"
	FieldRequestRefererDefault         = `%{json.escape(req.http.referer)}V`
	FieldRequestUserAgent              = "request_user_agent"
	FieldRequestUserAgentDefault       = `%{json.escape(req.http.User-Agent)}V`
	FieldResponse                      = "response"
	FieldResponseDefault               = `%{resp.response}V`
	FieldResponseState                 = "response_state"
	FieldResponseStateDefault          = `%{json.escape(fastly_info.state)}V`
	FieldResponseStatus                = "response_status"
	FieldResponseStatusDefault         = `%{resp.status}V`
	FieldResponseReason                = "response_reason"
	FieldResponseReasonDefault         = `%{if(resp.response, "%22"+json.escape(resp.response)+"%22", "null")}V`
	FieldResponseBodySize              = "response_body_size"
	FieldResponseBodySizeDefault       = `%{resp.body_bytes_written}V`
	FieldFastlyServer                  = "fastly_server"
	FieldFastlyServerDefault           = `%{json.escape(server.identity)}V`
	FieldFastlyServerDatacenter        = "server_datacenter"
	FieldFastlyServerDatacenterDefault = `%{server.datacenter}V`
	FieldFastlyIsEdge                  = "fastly_is_edge"
	FieldFastlyIsEdgeDefault           = `%{if(fastly.ff.visits_this_service == 0, "true", "false")}V`
	FieldFastlyCacheState              = "fastly_cache_state"
	FieldFastlyCacheStateDefault       = `%{fastly_info.state}V`
	FieldWAFBlock                      = "waf_block"
	FieldWAFBlockDefault               = `%{req.http.X-Waf-Block}V`
	FieldWAFBlockID                    = "waf_block_id"
	FieldWAFBlockIDDefault             = `%{req.http.X-Waf-Block-Id}V`
)

type Entry struct {
	Timestamp        time.Time `json:"timestamp"`
	ClientIP         string
	GeoCountry       string
	GeoCity          string
	Host             string
	RequestMethod    string
	RequestProtocol  string
	RequestReferer   string
	RequestUserAgent string
	RequestState     string
	RequestStatus    string
	RequestReason    string
	RquestBodySize   string
	FastlyServer     string
	FastlyIsEdge     bool
}

func KnownFields() []string {
	fieldsMap := KnownFieldsDefaultMap()
	fields := []string{}
	for k := range fieldsMap {
		fields = append(fields, k)
	}
	sort.Strings(fields)
	return fields
}

func KnownFieldsDefaultMap() map[string]string {
	return map[string]string{
		FieldTimestamp:              FieldTimestampDefault,
		FieldServiceID:              FieldServiceIDDefault,
		FieldTimeElapsed:            FieldTimeElapsedDefault,
		FieldClientIP:               FieldClientIPDefault,
		FieldGeoCountry:             FieldGeoCountryDefault,
		FieldGeoCity:                FieldGeoCityDefault,
		FieldHost:                   FieldHostDefault,
		FieldURL:                    FieldURLDefault,
		FieldRequest:                FieldRequestDefault,
		FieldRequestMethod:          FieldRequestMethodDefault,
		FieldRequestProtocol:        FieldRequestProtocolDefault,
		FieldRequestReferer:         FieldRequestRefererDefault,
		FieldRequestUserAgent:       FieldRequestUserAgentDefault,
		FieldResponse:               FieldResponseDefault,
		FieldResponseState:          FieldResponseStateDefault,
		FieldResponseStatus:         FieldResponseStatusDefault,
		FieldResponseReason:         FieldResponseReasonDefault,
		FieldResponseBodySize:       FieldResponseBodySizeDefault,
		FieldFastlyServer:           FieldFastlyServerDefault,
		FieldFastlyIsEdge:           FieldFastlyIsEdgeDefault,
		FieldFastlyServerDatacenter: FieldFastlyServerDatacenterDefault,
		FieldFastlyCacheState:       FieldFastlyCacheStateDefault,
		FieldWAFBlock:               FieldWAFBlockDefault,
		FieldWAFBlockID:             FieldWAFBlockIDDefault,
	}
}

func EntryFormatDefault(fields ...string) (map[string]string, error) {
	fields = stringsutil.SliceCondenseSpace(fields, true, false)
	if len(fields) == 0 {
		fields = KnownFields()
	}
	format := map[string]string{}
	unknownFields := []string{}
	knownFieldsMap := KnownFieldsDefaultMap()
	for _, field := range fields {
		fieldCanonical := strings.ToLower(strings.TrimSpace(field))
		if len(fieldCanonical) == 0 {
			continue
		}
		if def, ok := knownFieldsMap[fieldCanonical]; ok {
			format[fieldCanonical] = def
		} else {
			unknownFields = append(unknownFields, field)
		}
	}
	unknownFields = stringsutil.SliceCondenseSpace(unknownFields, true, true)
	if len(unknownFields) > 0 {
		return format, fmt.Errorf("known fields [%s]", strings.Join(unknownFields, ","))
	}
	return format, nil
}

/*

{
 // "timestamp": "%{strftime(\{"%Y-%m-%dT%H:%M:%S%z"\}, time.start)}V",
 // "client_ip": "%{req.http.Fastly-Client-IP}V",
 // "geo_country": "%{client.geo.country_name}V",
 // "geo_city": "%{client.geo.city}V",
 // "host": "%{if(req.http.Fastly-Orig-Host, req.http.Fastly-Orig-Host, req.http.Host)}V",
 // "url": "%{json.escape(req.url)}V",
 // "request_method": "%{json.escape(req.method)}V",
 // "request_protocol": "%{json.escape(req.proto)}V",
 // "request_referer": "%{json.escape(req.http.referer)}V",
 // "request_user_agent": "%{json.escape(req.http.User-Agent)}V",
 // "response_state": "%{json.escape(fastly_info.state)}V",
 // "response_status": %{resp.status}V,
 // "response_reason": %{if(resp.response, "%22"+json.escape(resp.response)+"%22", "null")}V,
 // "response_body_size": %{resp.body_bytes_written}V,
 // "fastly_server": "%{json.escape(server.identity)}V",
 // "fastly_is_edge": %{if(fastly.ff.visits_this_service == 0, "true", "false")}V
}



{ "timestamp":%{time.start.sec}V,
"service_id":"%{req.service_id}V",
 "time_elapsed":%{time.elapsed.usec}V,
"is_edge":%{if(fastly.ff.visits_this_service == 0, "true", "false")}V,
"client_ip":"%{req.http.Fastly-Client-IP}V",
"client_as":"%{if(client.as.number!=0, client.as.number, "null")}V",
"geo_city":"%{client.geo.city}V", "geo_country_code":"%{client.geo.country_code}V",
 "request":"%{req.request}V",
 "host":"%{json.escape(req.http.Host)}V",

 "request_referer":"%{json.escape(req.http.Referer)}V",
"request_user_agent":"%{json.escape(req.http.User-Agent)}V",
"resp_status":"%s",
 "resp_message":"%{resp.response}V",
"resp_body_size":%{resp.body_bytes_written}V,

"server_datacenter":"%{server.datacenter}V",
"cache_status":"%{fastly_info.state}V",
"waf_block":"%{req.http.X-Waf-Block}V",
"waf_block_id":"%{req.http.X-Waf-Block-Id}V" }

*/
