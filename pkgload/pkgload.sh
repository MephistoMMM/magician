#!/usr/bin/env bash
#
# Author: Mephis Pheies <mephistommm@gmail.com>
# Date  : 2020-03-13

ROOT=$PWD

if [[ $# -lt 2 ]];then
    echo "Usage: pkgload /path/to/task_file /path/to/zip_file"
    exit 1
fi

TASK_FILE=$1
ZIP_FILE=$2
MAVEN_REPO=~/.m2/repository
TMP_DIR=$ROOT/.tmp
TMP_REPOS=$TMP_DIR/repository

if [[ ! -f $TASK_FILE ]]; then
    echo "Task file $TASK_FILE doesn't exist."
    exit 1
fi

# 1. generate pom.xml
dependencies=""

IFS=$'\n'       # make newlines the only separator
set -f          # disable globbing
for task in `cat $TASK_FILE`
do
    groupId=`echo $task | awk '{ print $1 }'`
    artifactId=`echo $task | awk '{ print $2 }'`
    version=`echo $task | awk '{ print $3 }'`
    typ=`echo $task | awk '{ print $4 }'`
    if [[ -z "$typ" ]]; then 
        typ="jar"
    fi
    dependency="\
        <dependency>\\
            <groupId>$groupId</groupId>\\
            <artifactId>$artifactId</artifactId>\\
            <version>$version</version>\\
            <type>$typ</type>\\
        </dependency>"
    dependencies="${dependencies}\\
${dependency}"
done

mkdir -p $TMP_DIR
cat $ROOT/pom.template \
    | sed "s~{{DEPENDENCIES}}~$dependencies~g" \
    > $TMP_DIR/pom.xml

# 2. download package and get dependencies tree
cd $TMP_DIR
export MAVEN_OPTS="-Dorg.slf4j.simpleLogger.defaultLogLevel=INFO"
mvn dependency:list -Dclassifier=sources -DincludeParents=true \
    && mvn -B dependency:list -DincludeParents=true > deplog
if [[ $? -ne 0 ]]; then
    echo "Download dependencies failed."
    exit 1
fi
cat deplog | grep -E "(:compile$)|(:pom:)|(:jar:)" | awk '{ print $2 }' \
    | sed 's/:jar:/:/g' \
    | sed 's/:no_aop:/:/g' \
    | sed 's/:compile$//g' \
    | sed 's/:pom:/:/g' \
    > deplist

rm depPath
IFS=$'\n'       # make newlines the only separator
set -f          # disable globbing
for dep in `cat deplist`
do
    groupId=`echo $dep | awk -F ':' '{ print $1 }' | sed 's/\./\//g'`
    artifactId=`echo $dep | awk -F ':' '{ print $2 }'`
    version=`echo $dep | awk -F ':' '{ print $3 }'`
    jarpath="$groupId/$artifactId/$version"
    echo "$jarpath" >> depPath
done

# 3. package downloaded dependencies
rm -rf $TMP_REPOS && mkdir -p $TMP_REPOS
IFS=$'\n'       # make newlines the only separator
set -f          # disable globbing
for pkg in `cat $TMP_DIR/depPath`
do
    orgDir=$MAVEN_REPO/$pkg
    if [[ ! -d $orgDir ]]; then
        echo "Error: not exists: $orgDir. Ignore It."
        continue
    fi

    tmpDir=$TMP_REPOS/$pkg
    mkdir -p $tmpDir
    cp -a $orgDir/. $tmpDir
done

zip -r $ZIP_FILE ./repository
if [[ $? -ne 0 ]]; then
    echo "Error: package failed."
    exit 1
fi

echo "Done. Packaged to $ZIP_FILE"
