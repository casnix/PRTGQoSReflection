#!/bin/bash

# Build script for PRTGQoSReflection.  This script may become my 'autobild' project.

#####################
## Dependencies:   ##
##  ./buildinfo.sh ##
##  ./ansi.sh      ##
#####################
source ./ansi.sh
source ./buildinfo.sh


##############################################################################
####
##
## Validate build arguments, print help.
if [[ "$1" != "linux" ]] && [[ "$1" != "windows" ]] && [[ "$1" != "macos" ]]; then
	echo "Usage: $0 <linux|windows|macos>"
	echo "    linux        - Will build with linux/unix naming style."
	echo "    windows      - Will build a .EXE file."
	echo "    macos        - Alias of linux."
	exit
fi

##############################################################################
####
##
## Print information, set up build environment.
echo -e "${Yellow}Building ${Red}$AppName${Yellow} release version ${Purple}$ReleaseVersion.${NC}"
echo -e "${Yellow}This build version will be ${Purple}$BuildVersion.${NC}"
echo -e "${Red}$AppName${Yellow}, also known as ${Red}$ShortAppName${Orange}, $CopyrightNotice.${NC}"

echo -e "${Yellow}Checking for ./dynamicgeneration/ directory...${NC}"
if [[ ! -d "./dynamicgeneration" ]]; then
  echo -e "${LightRed}Not found.  Making directory...${NC}"
  mkdir ./dynamicgeneration
fi

##############################################################################
#### NOT SUPPORTED FOR THIS PROGRAM YET
##
## Build ./dynamicgeneration/startservices.sh
## Depends on:
##   TEMPLATE ./startservice.sh
#echo -e "${Yellow}Building ./dynamicgeneration/startservices.sh based on ./startservice.sh:${NC}"
#echo "  {{ShortAppName}} = $ShortAppName"
#echo "  {{ReleaseVersion}} = $ReleaseVersion"
#echo "  {{BuildVersion}} = $BuildVersion"
#echo "  {{ShortAppPath}} = $ShortAppPath"

#sed "s/{{ShortAppName}}/$ShortAppName/" ./startservice.sh \
#| sed "s/{{ReleaseVersion}}/$ReleaseVersion/" \
#| sed "s/{{BuildVersion}}/$BuildVersion/" \
#| sed "s/{{ShortAppPath}}/$ShortAppPath/" \
#> ./dynamicgeneration/startservices.sh

##############################################################################
####
##
## Build ./README.md
## Depends on:
##   TEMPLATE ./readme-template.md
echo -e "${Yellow}Cleaning ./README.md...${NC}"
echo -e "${Yellow}Building ./README.md based on ./readme-template.md...${NC}"
if [[ -f "./README.md" ]]; then
	rm ./README.md
fi

sed "s/{{ShortAppName}}/$ShortAppName/" ./readme-template.md \
	| sed "s/{{ReleaseVersion}}/$ReleaseVersion/" \
	| sed "s/{{BuildVersion}}/$BuildVersion/" \
	| sed "s/{{ShortAppPath}}/$ShortAppPath/" \
	| sed "s/{{AppName}}/$AppName/" \
	> ./README.md

##############################################################################
####
##
## Run dynamic build info script.
## Depends on:
##   SCRIPT ./buildinfo.sh
echo -e "${Yellow}Triggering dynamic buildinfo script...${NC}"
Generate
echo -e "${LightRed}No QA has been built to verify this.  De facto QA will happen during compile.${NC}"

##############################################################################
####
##
## Build Go program!
## Depends on:
##   GOLANG ./dynamicgeneration/buildinfo.go
echo -e "${Yellow}Building Go program...${NC}"
if [[ "$1" == "linux" ]] || [[ "$1" == "macos" ]]; then
	go build -o PRTGQoSReflection
fi
if [[ "$1" == "windows" ]]; then
	go build -o PRTGQoSReflection.exe
fi

