#!/bin/bash

# Install path: /usr/local/bin/git-hook.sh

# This file is used to keep track of when the git-hook was last run
TS_FILE="/var/run/git-hook.timestamp"

# The minimum number of seconds between runs of the git-hook script
MIN_INTERVAL=60

# The log file for the git-hook service
LOG_FILE=/home/git-hook/git-hook-service.log

# The nginx log file for the git-hook URL
TAIL_FILE=/var/log/nginx/git-hook.log

# Get the currently checked-out branch of a git repository
function getCurrentBranch {
    git branch | awk '/^\*/{print $2}'
}

# Function that updates a public_html folder, given a path as first argument
function updatePublicHtml {
    pushd .
    cd frontend
    cp *.html *.css *.js *.php "$1/"
    cp -r assets "$1/"
    <index.html sed "s/GIT_VERSION_TAG/$CID/" >"$1/index.html"
    popd
}

# Function that updates backend and restarts service, given file/service name as first argument
function updateBackend {
    # Delete previous build of backend
    rm backend/backend
    pushd .
    cd backend
    go build backend
    if [[ -x ./backend ]]; then
        cp backend /usr/local/bin/$1
        systemctl restart $1
    else
        >$LOG_FILE echo "[$1]: Backend did not build successfully."
    fi
    popd
}

function runGitHook {
    touch "$TS_FILE"
    git pull
    git checkout master

    updatePublicHtml "/var/www/live/public_html"
    updateBackend "backend"

    git checkout dev

    updatePublicHtml "/var/www/dev/public_html"
    updateBackend "backend-dev"
}

# Check if it has been at least 
function shouldRunGitHook {
    LAST_RUN_TIMESTAMP=$(stat -c '%Y' "$TS_FILE") || LAST_RUN_TIMESTAMP=0
    CURRENT_TIMESTAMP=$(date +'%s')

    if [[ $(( $LAST_RUN_TIMESTAMP - $CURRENT_TIMESTAMP )) -gt $MIN_INTERVAL ]]; then
        echo "Running git-hook, over $MIN_INTERVAL seconds since last run."
        runGitHook
    fi
}

tail --follow=name --retry $TAIL_FILE | while read; do
    shouldRunGitHook && runGitHook
done