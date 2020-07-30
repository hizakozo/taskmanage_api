#!/bin/bash

for file_name in `ls *.jpg`
do
aws s3 cp $file_name s3://taskmanage/user/ --endpoint-url=http://localhost:4572
done