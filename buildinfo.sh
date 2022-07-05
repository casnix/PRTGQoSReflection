#!/bin/bash

source ./ansi.sh

outputDirectory="./dynamicgeneration/buildinfo"
outputFile="$outputDirectory/buildinfo.go"

export AppName="PRTGQoSReflection"
export ShortAppName="PRTGQoSReflection"
export ShortAppPath="\$PRTGQOSPATH"
export ReleaseVersion="v1.0.0"
export BuildVersion="$(date +%Y%m%d).$(echo $(date +%H%M).$(whoami) | base64 | sed 's/==$//')"
export CopyrightNotice="Copyright 2022 Matthew Rienzo for Southwestern Healthcare Inc."

IFS= read -r -d '' GeneratedHeader << EndOfHeader
// $outputFile -- Dynamically generated build info for $ShortAppName
/* Written by Matt Rienzo for Southwestern Healthcare, Inc.'s IT */
/* department.                                                   */
/*---------------------------------------------------------------*/
// $CopyrightNotice
/* Licensed under the Apache License, Version 2.0 (the "License" */
/* ); you may not use this file except in compliance with the    */
/* license.  You may obtain a copy of the license at             */
/*     http://www.apache.org/licenses/LICENSE-2.0                */
/* Unless required by applicable law or agreed to in writing,    */
/* software distributed under the Licenses is deistributed on an  */
/* "AS IS" BASIS, WITHOUT WARRENTIES OR CONDITIONS OF ANY KIND,  */
/* either express or implied.  See the License for specific      */
/* language governing permissions and limitations under the      */
/* License.                                                      */
/*****************************************************************/

package buildinfo

const AppName string = "$AppName"
const ShortAppName string = "$ShortAppName"
const ReleaseVersion string = "$ReleaseVersion"
const BuildVersion string = "$BuildVersion"
const CopyrightNotice string = "$CopyrightNotice"
const LogName = "[$ShortAppName]"
EndOfHeader

Generate(){
    CWD=$(pwd)
    echo -e "${Cyan}[buildinfo.sh]${Yellow} Checking for $outputDirectory...${NC}"
    if [[ ! -d "$outputDirectory" ]]; then
        echo -e "${Cyan}[buildinfo.sh]${Red} Not found. Making directory...${NC}"
        mkdir $outputDirectory
    fi

    echo -e "${Cyan}[buildinfo.sh]${Yellow} Checking for $outputDirectory/go.mod...${NC}"
    if [[ ! -f "$outputDirectory/go.mod" ]]; then
        echo -e "${Cyan}[buildinfo.sh]${Red} Not found. Initializing module...${NC}"
        cd $outputDirectory
        go mod init github.com/casnix/PRTGQoSReflection/buildinfo
        cd $CWD
    fi

    echo "$GeneratedHeader" > $outputFile
    echo -e "${Cyan}[buildinfo.sh]${Yellow} Placed generated header into $outputFile.${NC}"
}
