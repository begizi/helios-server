# Helios Server


## Compiling
This is used for building the release but can also be used for development:

    git clone git@github.com:begizi/helios-server.git
    cd helios-server && make build
    ./bin/helios


## Enable Auto Compile (optional)
The strategy I use for rapid development with auto compiling:

Add a default GOPATH and add that GOPATH to the PATH:

    export GOPATH=/Users/<username>/go
    export PATH=$GOPATH/bin/:$PATH

Install the gin watch tool and verify it installed correctly:

    go get github.com/codegangsta/gin
    gin -h

Use the gin build command to auto compile when changes are made:

    make gin
