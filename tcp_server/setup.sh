#!/bin/bash
for (( c=0; c<20; c++ ))
do
    ## sudo docker run -v $(pwd)/:/go --name golang_tcp_server_$c -d --rm golang:1.12.1 go run simple_tcp_server/server.go && go run client.go -conn=20000 &
    sudo docker run -v $(pwd)/client:/client --name golang_tcp_server_$c -d --rm alpine /client -conn=10000 -ip=172.17.0.1
done