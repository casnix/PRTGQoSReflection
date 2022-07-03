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


import (
	"log"
	"net"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/BurntSushi/toml"
)


// var cli_opts struct -- The potential command line arguments.
var cli_opts struct {
	Port 		string 	`short:"p" long:"port" description:"Provide port defined in PRTG"`
	Configuration 	string 	`short:"c" long:"conf" description:"Path of config file, if not provided default qosreflect.conf will be used"`
	Host 		string 	`short:"o" long:"host" description:"Provide the IP address if the interface the script should bind to.\nUse ''All'' to bind to all available interfaces (recommended)"`
	ReplyIP		string 	`short:"r" long:"replyip" description:"Provide the IP address of the PRTG Probe which sends the packets.\nThe reflector will then only reply to this IP"`
	ReplyPort 	string 	`short:"t" long:"replyport" description:"Provide the port the packets should be bounced to"`
	NATMode 	bool 	`short:"n" long:"nat" description:"Option enables the NAT mode so packets are reflected exactly to the port they are received from"`
	Debug 		bool 	`short:"d" long:"debug" description:"Set to turn on detailed output"`
}

// ParseCLI(void) error -- Parses command line arguments
// Input: none
// Output: None, modifies cli_opts.
func ParseCLI() error  {
	var parser = flags.NewParser(&cli_opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}

		return error
	}

	return nil
}

// ReadConfig(void) error -- Parses configuration file, if there is one.
// Input: none
// Output: none, modifies cli_opts.
func ReadConfig() error {
	var configFile = cli_opts.Configuration
	_, err: os.Stat(configFile)
	if err != nil {
		log.Fatal("Configuration file is missing: ", configFile)
	}

	if _, err := toml.DecodeFile(configFile, &cli_opts); err != nil {
		log.Fatal(err)
	}

	return nil
}

// main(void) -- Entry point for program.
func main() {
	// Local variable declaration space
	var hOST = ""
	var pORT = "0"
	var rESTRICT = True

	ParseCLI()
	
	// Validate configuration options
	if cli_opts.Configuration == "" { // No configuration file specified
		// Validate that required options are set in this case
		if cli_opts.ReplyIP == "" || cli_opts.Host == "" || cli_opts.Port == "" { // We must fall back to the configuration file
			cli_opts.Configuration = "qosreflect.conf"
			ReadConfig()
		}
	} else {
		ReadConfig()
	}

	if cli_opts.Host != "All" {
		hOST = cli_opts.Host
	}

	if cli_opts.ReplyPort != "" {
		pORT = cli_opts.ReplyPort
	} else {
		pORT = cli_opts.Port
	}

	if cli_opts.ReplyIP == "None" || cli_opts.ReplyIP == "" {
		rESTRICT = False
	}

	

	server, err := net.ListenPacket("udp", ":"+pORT)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	for {
		udpBuffer := make([]byte, 1024)
		index, address, err := server.ReadFrom(udpBuffer)
		if err != nil {
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
	// 0 - 1: ID
	// 2: QR(1): Opcode(4)
	udpBuffer[2] |= 0x80 // set qr bit
	server.WriteTo(udpBuffer, address)
}
