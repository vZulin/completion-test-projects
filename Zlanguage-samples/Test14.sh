#!/bin/bash

# This script creates and pushes a new branch, which commits the history of HEAD movements and the changes that were made

REPORT_ISSUE_TEMPLATE_URL="https://kotl.in/report-jps-ic-bug"
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

checkStatus() {
    if [ $? -eq 0 ]; then
        isSuccessful=true
    else
        log "!!!Command failed!!!"
        isSuccessful=false
    fi
}

# Log with bold formatting
log() {
    echo -e "\033[1m$1\033[0m"
}

log "Current branch is $CURRENT_BRANCH"

# In case $User is empty, use host name
if [ -z "$USER" ]; then
    USER=$(hostname)
fi

# Replace all illegal symbols to "_"
USER_NAME=$(echo "$USER" | sed 's/[^a-zA-Z0-9._-]/_/g')

NEW_BRANCH=$USER_NAME-ic-bug-$(date +"%Y-%m-%d-%H-%M-%S")
log "Creating branch $NEW_BRANCH..."
git checkout -B $NEW_BRANCH
checkStatus

log "Create reflog..."
git reflog --date=iso --all > git-ref.log
checkStatus
ls "$(realpath git-ref.log)"

log "Adding changes..."
git add -A .
checkStatus
git status

log "Commiting changes..."
git commit -m"IC bug changes"
checkStatus

log "Push $NEW_BRANCH..."
git push -f origin $NEW_BRANCH
checkStatus

# Return to initial state
log "Cleaning..."
git reset --soft HEAD~
checkStatus
rm git-ref.log
checkStatus

#
git checkout $CURRENT_BRANCH
checkStatus

# Remove temporary branch
git branch -d $NEW_BRANCH
checkStatus

log "Push report to $NEW_BRANCH branch finished"

if [ "$isSuccessful" = true ]; then
    log "Script executed successfully"
else
    log "Script encountered an error. Please check output for '!!!Command failed!!!' message"
fi

log "[ACTION NEEDED] Please, report an issue with branch name: $REPORT_ISSUE_TEMPLATE_URL"

