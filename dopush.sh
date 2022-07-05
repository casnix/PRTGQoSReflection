#!/bin/sh

if [[ -f "./PRTGQoSReflection.exe" ]]; then
	rm ./PRTGQoSReflection.exe
fi
if [[ -f "./PRTGQoSReflection" ]]; then
	rm ./PRTGQoSReflection
fi
if [[ -d "./dynamicgeneration" ]]; then
	rm -rf ./dynamicgeneration
fi

git add .
git commit -m "$1"
git push
