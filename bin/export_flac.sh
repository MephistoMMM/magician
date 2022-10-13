#!/bin/bash
#
# The script to export audio from bilibili video.
#
# Usage: ./export_flac.sh <TASK_LIST_FILE> <EXPORT_DIR>
#
# Format Of TASK_LIST_FILE:
#   name_of_flac_file     av_code_of_video
#
# OS: macOS
#
# Author: Mephis Pheies <mephistommm@gmail.com>
# Date: 2020-03-08

ROOT=`pwd`

if [[ $# -lt 2 ]]; then
    echo "Usage: ./export_flac.sh <TASK_LIST_FILE> <EXPORT_DIR>"
    exit 1
fi

# workspace
TODO_LIST=$1
EXPORT_DIR=$2
TMP_DIR=$EXPORT_DIR/.tmp
FINISH_LIST=$TMP_DIR/finish_list

# constants
BILIBILI_VIDEO_URL_PREFIX=https://www.bilibili.com/video


# check you-get
if type -a you-get > /dev/null 2>&1; then
    echo "you-get has been installed."
else
    brew install you-get
fi

# check ffmpeg
if type -a ffmpeg > /dev/null 2>&1; then
    echo "ffmpeg has been installed."
else
    brew install ffmpeg
fi

if [[ ! -f $TODO_LIST ]];then
    echo "None task."
    exit 1
fi

# $1 stream file
function highest_flv_format() {
    flvlist=(flv flv1080 flv720 flv480 flv360)
    for flv in ${flvlist[@]};
    do
        cat $1 | grep " format:" | grep -v "dash" | grep "$flv" > /dev/null 2>&1
        if [[ $? -eq 0 ]]; then 
            echo "$flv"
            return 0
        fi
    done
}


# $1 task name
# $2 task code
function step_download_video() {
    if ls $TMP_DIR/${2}/${2}.flv > /dev/null 2>&1; then
        return 0
    fi

    mkdir -p $TMP_DIR/$2
    download_url=$BILIBILI_VIDEO_URL_PREFIX/$2
    you-get -i $download_url > $TMP_DIR/streams
    you-get --format=$(highest_flv_format $TMP_DIR/streams) -o $TMP_DIR/$2 -O $2 $download_url
    return $?
}

# $1 task name
# $2 task code
function step_export_audio() {
    if ls $EXPORT_DIR/${1}.m4a > /dev/null 2>&1; then
        return 0
    fi

    ffmpeg -i $TMP_DIR/${2}/${2}.flv -acodec alac -vn $EXPORT_DIR/${1}.m4a
    return $?
}

# $1 task name
# $2 task code
function step_finish_and_clean() {
    echo "$1" >> $FINISH_LIST && rm -rf $TMP_DIR/${2}
    return $?
}

IFS=$'\n'
set -f
for task in `cat $TODO_LIST`
do
    echo "Start Task: $task"

    task_name=`echo $task | awk '{ print $1 }'`
    task_code=`echo $task | awk '{ print $2 }'`

    if grep -F "$task_name" $FINISH_LIST > /dev/null 2>&1; then
        continue
    fi

    step_download_video $task_name $task_code
    if [[ ! $? -eq 0 ]]; then
        echo "Task $task_name is failed in step DOWNLOADING_VIDEO."
        continue
    fi

    step_export_audio $task_name $task_code
    if [[ ! $? -eq 0 ]]; then
        echo "Task $task_name is failed in step EXPORT_AUDIO."
        continue
    fi

    step_finish_and_clean $task_name $task_code
    if [[ ! $? -eq 0 ]]; then
        echo "Task $task_name is failed in step FINISH."
        continue
    fi

    echo "Finish Task: $task"
done
