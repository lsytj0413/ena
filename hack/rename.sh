#!/bin/bash -x
# Copyright (c) 2023 The Songlin Yang Authors
#
# Permission is hereby granted, free of charge, to any person obtaining a copy of
# this software and associated documentation files (the "Software"), to deal in
# the Software without restriction, including without limitation the rights to
# use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
# the Software, and to permit persons to whom the Software is furnished to do so,
# subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
# FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
# COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
# IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
# CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.


usage() {
    echo "$PROGNAME: usage: $PROGNAME new_project_name"
    return
}

PROGNAME=$(basename $0)
BASEPATH=$(cd `dirname $0` && pwd)
dir=$(dirname ${BASEPATH})
nname=$1

if [[ -z $nname ]]; then
    usage >&2
    echo "MUST specify new project name"
    exit 1
fi

shift
while [[ -n $1 ]]; do
    usage >&2
    exit 1
done

echo "Rename the project to $nname, under directory: $dir"

find $dir -type f -not -path '*/\.git/*' -not -path '*rename\.sh*' -exec sed -i '' "s/golang-project-template/${nname}/g" {} \;