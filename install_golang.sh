#!/bin/bash
echo "install golang"
mkdir -p $HOME/tools
cd $HOME/tools
curl -o go.tar -L https://go.dev/dl/go1.22.1.linux-amd64.tar.gz
tar -xf go.tar
mkdir -p $HOME/gopath
echo 'export GOPATH=$HOME/gopath' >> $HOME/.bashrc
echo 'export PATH=$PATH:$HOME/tools/go/bin:$GOPATH/bin' >> $HOME/.bashrc
export GOPATH=$HOME/gopath
export PATH=$PATH:$HOME/tools/go/bin:$GOPATH/bin
go install golang.org/x/tools/gopls@latest
