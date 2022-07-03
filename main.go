/*****************************************************************/
/* main.go -- Entry point for Golang PRTGQoSReflection.          */
/*                                                               */
/* Written by Matt Rienzo for Southwestern Healthcare, Inc.'s IT */
/* department.                                                   */
/*---------------------------------------------------------------*/
/* Copyright 2022 Matt Rienzo                                    */
/*                                                               */
/* Licensed under the Apache License, Version 2.0 (the           */
/* "License"); you may not use this file except in compliance    */
/* with the license.  You may obtain a copy of the license at    */
/*    http://www.apache.org/licenses/LICENSE-2.0                 */
/* Unless required by applicable law or agreed to in writing,    */
/* software distributed under the License is distributed on an   */
/* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,  */
/* either express or implied.  See the License for specific      */
/* language governing permissions and limitations under the      */
/* License.                                                      */
/*****************************************************************/

package main

import (
	"log"
	"net"
	"os"

	"github.com/casnix/PRTGQoSReflection/buildinfo"
	"github.com/jessevdk/go-flags"
	"github.com/BurntSushi/toml"
	"github.com/TwiN/go-color"
)


// var cli_opts struct -- The potential command line arguments.
var cli_opts struct {
	Port 		string 	`short:"p" long:"port" description:"Provide port defined in PRTG"`
	Configuration 	string 	`short:"c" long:"conf" description:"Path of config file, if not provided default qosreflect.conf will be used"`
	Host 		string 	`short:"o" long:"host" description:"Provide the IP address if the interface the script should bind to.\nUse ''All'' to bind to all available interfaces (recommended)"`
	ReplyIP		string 	`short:"r" long:"replyip" description:"NOT SUPPORTED YET Provide the IP address of the PRTG Probe which sends the packets.\nThe reflector will then only reply to this IP"`
	ReplyPort 	string 	`short:"t" long:"replyport" description:"Provide the port the packets should be bounced to"`
	NATMode 	bool 	`short:"n" long:"nat" description:"NOT SUPPORTED YET Option enables the NAT mode so packets are reflected exactly to the port they are received from"`
	Debug 		bool 	`short:"d" long:"debug" description:"NOT SUPPORTED YET Set to turn on detailed output"`
}


// IsFlagPassed(string) bool -- Checks if CLI flag was passed
// Input: - name string		-- Flag name
// Output: - true		-- Flag found
// 	   - false		-- Flag not found
/*func IsFlagPassed(name string) bool {
	found := false
	flag.Visit(fun(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	
	return found
}
*/

// ParseCLI(void) error -- Parses command line arguments
// Input: none
// Output: None, modifies cli_opts.
func ParseCLI() error  {
	// go-flags is broken for error handling.  May switch to https://github.com/geomyidia/flagswrap in the future.
	var parser = flags.NewParser(&cli_opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if w, ok := err.(*flags.Error); ok {
			if w.Type == flags.ErrHelp {
				os.Exit(0)
			}
		}
		os.Exit(1)
	}

	return nil
}

// ReadConfig(void) error -- Parses configuration file, if there is one.
// Input: none
// Output: none, modifies cli_opts.
func ReadConfig() error {
	var configFile = cli_opts.Configuration
	_, err := os.Stat(configFile)
	if err != nil {
		log.Fatalf("%s Configuration file is missing: %s", color.Ize(color.Red, buildinfo.LogName), color.Ize(color.Red, configFile))
	}

	if _, err := toml.DecodeFile(configFile, &cli_opts); err != nil {
		log.Fatal(err)
	}

	return nil
}

// main(void) -- Entry point for program.
func main() {
	log.Printf("%s %s release %s, build %s is starting...\n", color.Ize(color.Cyan, buildinfo.LogName), buildinfo.AppName, color.Ize(color.Purple, buildinfo.ReleaseVersion), color.Ize(color.Purple, buildinfo.BuildVersion))

	// Local variable declaration space
	var hOST = ""
	var pORT = "0"
	var rESTRICT = true

	ParseCLI()
	
	// FIX THIS, NEED TO CHECK GO-FLAGS FOR IF FLAG IS PASSED
	// Validate configuration options
	if false {
		ReadConfig()
	}
	/*if !IsFlagPassed("Configuration") { // No configuration file specified
		log.Printf("%s debug cli_opts.Configuration == %s", color.Ize(color.Cyan, buildinfo.LogName), color.Ize(color.Yellow, cli_opts.Configuration))
		// Validate that required options are set in this case
		if !IsFlagPassed("ReplyIP") || !IsFlagPassed("Host") || !IsFlagPassed("Port") { // We must fall back to the configuration file
			log.Printf("%s debug cli_opts.(ReplyIP, Host, Port) == (%s, %s, %s)", color.Ize(color.Cyan, buildinfo.LogName), color.Ize(color.Yellow, cli_opts.ReplyIP), color.Ize(color.Yellow, cli_opts.Host), color.Ize(color.Yellow, cli_opts.Port))
			cli_opts.Configuration = "qosreflect.conf"
			ReadConfig()
		}
	} else {
		ReadConfig()
	}*/

	if cli_opts.Host != "All" {
		hOST = cli_opts.Host
	}

	if cli_opts.ReplyPort != "" {
		pORT = cli_opts.ReplyPort
	} else {
		pORT = cli_opts.Port
	}

	if cli_opts.ReplyIP == "None" || cli_opts.ReplyIP == "" {
		rESTRICT = false
	}

	

	server, err := net.ListenPacket("udp", hOST+":"+pORT)
	if err != nil {
		log.Fatalf("%s%s Fatal error: %s", color.Ize(color.Red, buildinfo.LogName), color.Ize(color.Yellow, "{0}"), color.Ize(color.Red, err.Error()))
	}
	defer server.Close()
	
	log.Printf("%s%s UDP reflection server is running, listening on %s:%s.\n", color.Ize(color.Cyan, buildinfo.LogName), color.Ize(color.Yellow, "{0}"), hOST, pORT)
	// Add check for replyIP.
	rESTRICT = false
	for {
		udpBuffer := make([]byte, 1024)
		index, address, err := server.ReadFrom(udpBuffer)
		if err != nil && !rESTRICT {
			continue
		}
		go Reflect(server, address, udpBuffer[:index])
	}
}


// Reflect(net.PacketConn, net.Addr, []byte) void -- Reflects/serves UPD packet
// Input: - net.PacketConn		- PacketConn object with connect data
//	  - net.Addr			- Addr object with network address data
//	  - []byte			- Byte buffer with UDP data to serve
// Outpu: None.
func Reflect(server net.PacketConn, address net.Addr, udpBuffer []byte) {
	server.WriteTo(udpBuffer, address)
}
