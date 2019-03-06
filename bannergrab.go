// Copyright © 2019 S. van der Baan <steven@vdbaan.net>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"sort"
)

const programVersion = "4.0\nCopyright (C) 2019 S. van der Baan\n"
const programBanner =
`  ________      __________       	                             
 /  _____/  ____\______   \_____    ____   ____   ___________ 
/   \  ___ /  _ \|    |  _/\__  \  /    \ /    \_/ __ \_  __ \
\    \_\  (  <_> )    |   \ / __ \|   |  \   |  \  ___/|  | \/
 \______  /\____/|______  /(____  /___|  /___|  /\___  >__|   
        \/              \/      \/     \/     \/     \/       

              Version 4.0
              http://github.com/vdbaan/gobanner
              Copyright (C) 2019 S. van der Baan

`
type triggerConfig struct{
	size int
	trigger string
}

type serviceConfig struct{
	port int
	tcp bool
	show bool
	connectBanner bool
	ssl bool
	timeout int
	trigger []triggerConfig
}

var (

	trig_default  = []triggerConfig{{0, "GET / HTTP/1.0\r\n\r\n"},{0, "HELP\r\n"}}
	trig_null     = []triggerConfig{{0, ""}}
	trig_echo     = []triggerConfig{{0, "Echo\r\n"}}
	trig_ftp      = []triggerConfig{{0, "HELP\r\n"},{0, "USER anonymous\r\n"},{0, "PASS banner@grab.com\r\n"},{0, "QUIT\r\n"}}
	trig_telnet   = []triggerConfig{{0, "\r\n"}, {0, "\r\n",}}
	trig_smtp     = []triggerConfig{{0, "HELO bannergrab.com\r\n"},{0, "HELP\r\n"},{0, "VRFY postmaster\r\n"}, {0, "VRFY bannergrab123\r\n"},{0, "EXPN postmaster\r\n"},{0, "QUIT\r\n"}}
	trig_finger   = []triggerConfig{{0, "root bin lp wheel spool adm mail postmaster news uucp snmp daemon\r\n"}}
	trig_http     = []triggerConfig{{0, "OPTIONS / HTTP/1.0\r\n\r\n"}}
	trig_pop      = []triggerConfig{{0, "QUIT\r\n"}}
	trig_nntp     = []triggerConfig{{0, "HELP\r\n"},{0, "LIST NEWSGROUPS\r\n"},{0, "QUIT\r\n"}}
	trig_ntp      = []triggerConfig{{48, "\xe3\x00\x04\xfa\x00\x01\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xca\x9b\xa3\x35\x2d\x7f\x95\x0b"},{12, "\x16\x02\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00"},{12, "\x16\x01\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00"}}
	trig_nbns     = []triggerConfig{{50, "\xa2\x48\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x20\x43\x4b\x41\x41\x41\x41\x41\x41\x41\x41\x41\x41\x41\x41\x41\x41\x41\x41\x41\x41\x41\x41\x41\x41\x41\x41\x41\x41\x41\x41\x41\x41\x00\x00\x21\x00\x01"}}
	trig_snmp     = []triggerConfig{{43, "\x30\x29\x02\x01\x00\x04\x06\x70\x75\x62\x6c\x69\x63\xa0\x1c\x02\x04\xff\xff\xff\xff\x02\x01\x00\x02\x01\x00\x30\x0e\x30\x0c\x06\x08\x2b\x06\x01\x02\x01\x01\x01\x00\x05\x00"},{44, "\x30\x2a\x02\x01\x00\x04\x07\x70\x72\x69\x76\x61\x74\x65\xa0\x1c\x02\x04\xff\xff\xff\xfe\x02\x01\x00\x02\x01\x00\x30\x0e\x30\x0c\x06\x08\x2b\x06\x01\x02\x01\x01\x01\x00\x05\x00"}}
	trig_fw1admin = []triggerConfig{{0, "???\r\n?\r\n"}}
	trig_isakmp   = []triggerConfig{{336, "\x22\xde\x92\x3f\x69\x61\xcc\xe2\x00\x00\x00\x00\x00\x00\x00\x00\x01\x10\x02\x00\x00\x00\x00\x00\x00\x00\x01\x50\x00\x00\x01\x34\x00\x00\x00\x01\x00\x00\x00\x01\x00\x00\x01\x28\x01\x01\x00\x08\x03\x00\x00\x24\x01\x01\x00\x00\x80\x01\x00\x05\x80\x02\x00\x02\x80\x03\x00\x01\x80\x04\x00\x02\x80\x0b\x00\x01\x00\x0c\x00\x04\x00\x00\x70\x80\x03\x00\x00\x24\x02\x01\x00\x00\x80\x01\x00\x05\x80\x02\x00\x01\x80\x03\x00\x01\x80\x04\x00\x02\x80\x0b\x00\x01\x00\x0c\x00\x04\x00\x00\x70\x80\x03\x00\x00\x24\x03\x01\x00\x00\x80\x01\x00\x01\x80\x02\x00\x02\x80\x03\x00\x01\x80\x04\x00\x02\x80\x0b\x00\x01\x00\x0c\x00\x04\x00\x00\x70\x80\x03\x00\x00\x24\x04\x01\x00\x00\x80\x01\x00\x01\x80\x02\x00\x01\x80\x03\x00\x01\x80\x04\x00\x02\x80\x0b\x00\x01\x00\x0c\x00\x04\x00\x00\x70\x80\x03\x00\x00\x24\x05\x01\x00\x00\x80\x01\x00\x05\x80\x02\x00\x02\x80\x03\x00\x01\x80\x04\x00\x01\x80\x0b\x00\x01\x00\x0c\x00\x04\x00\x00\x70\x80\x03\x00\x00\x24\x06\x01\x00\x00\x80\x01\x00\x05\x80\x02\x00\x01\x80\x03\x00\x01\x80\x04\x00\x01\x80\x0b\x00\x01\x00\x0c\x00\x04\x00\x00\x70\x80\x03\x00\x00\x24\x07\x01\x00\x00\x80\x01\x00\x01\x80\x02\x00\x02\x80\x03\x00\x01\x80\x04\x00\x01\x80\x0b\x00\x01\x00\x0c\x00\x04\x00\x00\x70\x80\x00\x00\x00\x24\x08\x01\x00\x00\x80\x01\x00\x01\x80\x02\x00\x01\x80\x03\x00\x01\x80\x04\x00\x01\x80\x0b\x00\x01\x00\x0c\x00\x04\x00\x00\x70\x80"}}
	trig_ldap     = []triggerConfig{{14, "\x30\x0c\x02\x01\x01\x60\x07\x02\x01\x03\x04\x00\x80\x00"},{55, "\x30\x35\x02\x01\x02\x63\x30\x04\x00\x0a\x01\x00\x0a\x01\x00\x02\x01\x00\x02\x01\x00\x01\x01\x00\x87\x0b\x6f\x62\x6a\x65\x63\x74\x43\x6c\x61\x73\x73\x30\x10\x04\x0e\x6e\x61\x6d\x69\x6e\x67\x43\x6f\x6e\x74\x65\x78\x74\x73"}}
	trig_mssql    = []triggerConfig{{224, "\x10\x01\x00\xe0\x00\x00\x01\x00\xd8\x00\x00\x00\x01\x00\x00\x71\x00\x00\x00\x00\x00\x00\x00\x07\x6c\x04\x00\x00\x00\x00\x00\x00\xe0\x03\x00\x00\x00\x00\x00\x00\x09\x08\x00\x00\x56\x00\x0a\x00\x6a\x00\x0a\x00\x7e\x00\x00\x00\x7e\x00\x20\x00\xbe\x00\x09\x00\x00\x00\x00\x00\xd0\x00\x04\x00\xd8\x00\x00\x00\xd8\x00\x00\x00\x00\x0c\x29\xc6\x63\x42\x00\x00\x00\x00\xc8\x00\x00\x00\x42\x00\x61\x00\x6e\x00\x6e\x00\x65\x00\x72\x00\x47\x00\x72\x00\x61\x00\x62\x00\x42\x00\x61\x00\x6e\x00\x6e\x00\x65\x00\x72\x00\x47\x00\x72\x00\x61\x00\x62\x00\x4d\x00\x69\x00\x63\x00\x72\x00\x6f\x00\x73\x00\x6f\x00\x66\x00\x74\x00\x20\x00\x44\x00\x61\x00\x74\x00\x61\x00\x20\x00\x41\x00\x63\x00\x63\x00\x65\x00\x73\x00\x73\x00\x20\x00\x43\x00\x6f\x00\x6d\x00\x70\x00\x6f\x00\x6e\x00\x65\x00\x6e\x00\x74\x00\x73\x00\x31\x00\x32\x00\x37\x00\x2e\x00\x30\x00\x2e\x00\x30\x00\x2e\x00\x31\x00\x4f\x00\x44\x00\x42\x00\x43\x00"}}

// Add from NMap	
// gpsd-info.nse
// http-cisco-anyconnect.nse
// openlookup-info.nse
// acarsd-info.nse
// irc-info.nse
// ganglia-info.nse
// metasploit-info.nse
// svn
// enip-info.nse

// Oracle
// Postgres


	services  = map[string]serviceConfig {
		"DEFAULT"    : {0,     true,  false, true,  false, 0, trig_default},
		"Echo"       : {7,     true,  true,  false, false, 0, trig_echo},
		"Discard"    : {9,     true,  false, false, false, 0, trig_echo},
		"Daytime"    : {13,    true,  false, true,  false, 0, trig_null},
		"QOTD"       : {17,    true,  false, true,  false, 0, trig_null},
		"Chargen"    : {19,    true,  false, true,  false, 0, trig_null},
		"FTP"        : {21,    true,  true,  true,  false, 0, trig_ftp},
		"SSH"        : {22,    true,  false, true,  false, 0, trig_null},
		"Telnet"     : {23,    true,  true,  true,  false, 0, trig_telnet},
		"SMTP"       : {25,    true,  true,  true,  false, 0, trig_smtp},
		"Finger"     : {79,    true,  false, false, false, 0, trig_finger},
		"HTTP"       : {80,    true,  false, false, false, 0, trig_http},
		"POP2"       : {109,   true,  false, true,  false, 0, trig_pop},
		"POP3"       : {110,   true,  false, true,  false, 0, trig_pop},
		"NNTP"       : {119,   true,  false, true,  false, 0, trig_nntp},
		"NTP"        : {123,   false, false, false, false, 2, trig_ntp},
		"NetBIOS-NS" : {137,   false, false, false, false, 2, trig_nbns},
		"SNMP"       : {161,   false, false, false, false, 2, trig_snmp},
		"FW1Admin"   : {256,   true,  false, true,  false, 0, trig_fw1admin},
		"FW1Auth"    : {259,   true,  true,  true,  false, 0, trig_telnet},
		"LDAP"       : {389,   true,  false, false, false, 0, trig_ldap},
		"HTTPS"      : {443,   true,  false, false, true,  0, trig_http},
		"ISA-KMP"    : {500,   false, false, false, false, 2, trig_isakmp},
		"Submission" : {587,   true,  true,  true,  false, 0, trig_smtp},
		"IPP"        : {631,   true,  false, false, false, 0, trig_http},
		"LDAPS"      : {636,   true,  false, false, true,  0, trig_ldap},
		"VMWare"     : {902,   true,  false, true,  false, 0, trig_null},
		"MSSQL"      : {1433,  true,  false, false, false, 0, trig_mssql},
		"MySQL"      : {3306,  true,  false, true,  false, 6, trig_null},
		"Printer"    : {9100,  true,  false, true,  false, 0, trig_null},
	}
)


func ShowTriggers() {
	// Go does not keep the order of maps, so to have it in the order of the ports we need to
	var m = make(map[int]string)
	var keys []int
	// capture all ports and services
	for s,i := range services {
		keys = append(keys, i.port)
		m[i.port] = s
	}
	// sort the ports
	sort.Ints(keys)
	
	log.Debug("Show Triggers")
	log.Infof("DEFAULT         (Port:N/A) ")
	// per port show the service info
	for _,port := range keys {		
		service := m[port]
		if service == "DEFAULT" {continue}
		tcp := services[service].tcp
		if tcp {
			log.Infof("%-15s (TCP Port:%d)",service, port)
		}else {
			log.Infof("%-15s (UDP Port:%d)",service, port)
		}
	}
}