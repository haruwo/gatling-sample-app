#!/bin/sh

cd $(dirname $0)
tag=$(docker build . | tee /dev/stderr | tail -1 | cut -d' ' -f3)
if [ "$tag" == "" ]; then
  exit 1
fi

exec docker run -p 8000:8000 --rm -it $tag

