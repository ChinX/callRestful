#!/bin/bash
basedir=${pwd}
gitdir=`git rev-parse --git-dir`
prjdir=${gitdir%/.git}

function formatFile()
{
	fname=$1
	if [ -f ${fname} ]; then
		if [ "${fname##*.}"x = "go"x ]; then
			goimports -w ${fname}
		elif [[ "${fname##*.}"x = "proto"x ]]; then
			clang-format -style="{BasedOnStyle: Google, IndentWidth: 4, ColumnLimit: 0, AlignConsecutiveAssignments: true, AlignConsecutiveAssignments: true}" -i ${fname}
		fi
	fi
}

for filepath in `git diff --name-only HEAD~ HEAD`
do
	formatFile "${prjdir}/${filepath}"
done

for filepath in `git diff --name-only HEAD`
do
	formatFile "${prjdir}/${filepath}"
done

for filepath in `git ls-files -o --exclude-standard`
do
	formatFile "${basedir}/${filepath}"
done
