#!/bin/sh -l

echo "Good day $1"
time=$(date)
echo "::set-output name=time::$time"
