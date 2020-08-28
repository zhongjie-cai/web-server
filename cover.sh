#!/bin/sh

# Load environment variables from second parameter if given
if [ "$1" = "" ];
then
    echo "Using system environment variables"
else
    echo "Loading environment variables from file: $1"
    export $(grep -v '^#' "$1" | xargs -0)
fi

mkdir -p coverage

allTestCoverResultFile="coverage/all.test.cover.result"
# Clean up all test coverage result file
rm -f $allTestCoverResultFile

for package in $(go list ./... | grep -v '/vendor/')
do
	echo "Running coverage for $package"

	outFileName="coverage/"$(echo $package | sed -e 's/\//-/g' )
	rm -f $outFileName.cover.json

	# Run test coverage
	go test -coverprofile $outFileName.cover.profile -v $package 2>&1 >> $allTestCoverResultFile

	# Generate coverage report HTML
	go tool cover -html=$outFileName.cover.profile -o $outFileName.cover.html
done

# Echo result for visibility
FullReportLines=$(grep -e '/web-server' $allTestCoverResultFile)

echo ""
echo "Coverage result:"
echo "$FullReportLines"

echo ""
echo "Analyzed result:"
echo ""

# all fully covered as default success
ReturnCode=0

WronglySkippedLines=$(echo "$FullReportLines" | grep -e 'no test files' | grep -e '/enum' -e '/constant' -e '/model' -v)
if [ -n "$WronglySkippedLines" ]
then
	# some wrongly skipped packages
	echo "Wrongly skipped packages:"
	echo "$WronglySkippedLines"
	echo ""
	ReturnCode=1
fi

UncoveredLines=$(echo "$FullReportLines" | grep -e 'ok' | grep -e '100.0% of statements' -e 'no test files' -e 'FAIL' -v)
if [ -n "$UncoveredLines" ]
then
	# some not fully covered
	echo "Not fully covered packages:"
	echo "$UncoveredLines"
	echo ""
	ReturnCode=2
fi

FailedLines=$(echo "$FullReportLines" | grep -e 'FAIL')
if [ -n "$FailedLines" ]
then
	# some tests failed during coverage call
	echo "Failed packages:"
	echo "$FailedLines"
	echo ""
	ReturnCode=3
fi

case $ReturnCode in
	0)
		echo "All code lines are covered! You are good to go!"
		;;
	1)
		echo "Some code packages are wrongly skipped for coverage! You should cover them ASAP before anything deployed!"
		;;
	2)
		echo "Some code lines are not covered! You should cover them ASAP before anything deployed!"
		;;
	3)
		echo "Some tests are not even passing! You should fix them NOW before anything committed!"
		;;
esac
exit $ReturnCode
