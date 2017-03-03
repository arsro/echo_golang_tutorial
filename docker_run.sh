#!/bin/sh
path=$(cd $(dirname $0) && pwd)
#echo $path
echo 'RUN echo docker!!!'
docker run --name echo -p 3000:3000 --volume $path:/app --link mysql_mysql_1:mysql -i -t --rm echo_image /bin/bash
