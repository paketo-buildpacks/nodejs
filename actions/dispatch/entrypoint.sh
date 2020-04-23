#!/bin/sh -l

echo "==="
env
echo "==="
cat "${GITHUB_EVENT_PATH}"
echo "==="
