#!/bin/sh

# -e  causes Exit immediately when a command exits with a non-zero status.
# When a test fails in Go, it returns a non-zero status and causes the 
# build pipeline to fail and exit when unit tests fails.
set -e

# Load environment variables from second parameter if given
if [ "$2" = "" ];
then
    echo "Using system environment variables"
else
    echo "Loading environment variables from file: $2"
    export $(grep -v '^#' "$2" | xargs -0)
fi

# Go Test all packages found within the main directory and subdirectories
# The pipe is required to reverse the order of the comamnd, so that tee does not swallow the exit code

if [ "$1" = "unit" ];
then
    # Delete existing files when running the unit tests
    rm -f testResultsPipe
    rm -f testResults.txt
    rm -f unitTestResults.txt

    mkfifo testResultsPipe
    tee unitTestResults.txt < testResultsPipe &
    go test $(go list ./../... | grep -v '/vendor/') -short -v 2>&1 > testResultsPipe
else
    if [ "$1" = "integration" ];
    then
        # Delete existing files when running the integration tests
        rm -f testResultsPipe
        rm -f testResults.txt
        rm -f integrationTestResults.txt

        mkfifo testResultsPipe
        tee -a integrationTestResults.txt < testResultsPipe &
        go test $(go list ./../... | grep -v '/vendor/') -run Integration -v 2>&1 > testResultsPipe
    else
        # Delete existing files when running the complete tests
        rm -f testResultsPipe
        rm -f testResults.txt
        rm -f unitTestResults.txt
        rm -f integrationTestResults.txt

        mkfifo testResultsPipe
        tee testResults.txt < testResultsPipe &
        go test $(go list ./../... | grep -v '/vendor/') -v 2>&1 > testResultsPipe
    fi
fi
