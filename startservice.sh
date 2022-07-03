#!/bin/bash

# {{ShortAppName}} release version: {{ReleaseVersion}}
# Build version: {{BuildVersion}}

echo "Moving into {{ShortAppName}} runtime directory..."
cd {{ShortAppPath}}/runtime
echo "Checking for stop flag..."
if [[ -f "./.stopapp" ]]; then
    echo "Stop flag found."
    echo "    Removing stop flag..."
    rm ./.stopapp
    echo "    Stopping..."
    exit
fi
if [[ -f "./.updatepls" ]]; then
    echo "Update flag found."
    echo "    Moving into {{ShortAppName}} updates directory..."
    cd {{ShortAppPath}}/updates
    echo "Resetting repo to last good merge..."
    git reset --hard
    echo "    Pulling updates from git remote..."
    git pull
    echo "    Building..."
    bash ./build.sh
    echo "    Organizing files..."
    bash ./organdize-me.sh
    echo "    Returning to {{ShortAppName}} runtime directory..."
    cd {{ShortAppPath}}/runtime
    echo "    Removing update flag..."
    rm ./.updatepls
    echo "    Restarting..."
    bash {{ShortAppPath}}/bin/startservice.sh
    exit
fi
echo "Starting {{ShortAppName}} app..."
{{ShortAppPath}}/bin/{{ShortAppName}}.exe
echo "App stopped. Restarting..."
bash {{ShortAppPath}}/bin/startservice.sh
