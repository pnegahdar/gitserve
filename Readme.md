## gitserve: webserver to serve git archives of any git path  

#### Installation:

Grab the right precompiled bin from github releases and put it in your path. Don't forget to `chmod +x` the bin.

OSX:

    curl -SL https://github.com/pnegahdar/gitserve/releases/download/0.1.0/gitserve_0.1.0_darwin_amd64.tar.gz \
        | tar -xzC /usr/local/bin --strip 1 && chmod +x /usr/local/bin/gitserve

Nix:

    curl -SL https://github.com/pnegahdar/gitserve/releases/download/0.1.0/gitserve_0.1.0_linux_amd64.tar.gz \
        | tar -xzC /usr/local/bin --strip 1 && chmod +x /usr/local/bin/gitserve

#### Usage:

##### Run the server

     # Minimal 
     gitserver 

     # All args  
     gitserver -root=/Users/myuser/git/project  -prefix=subdira/ -listen=":9000"
    
     
#### Get the archives:

    curl localhost:9000/.zip 
    curl localhost:9000/.tar 
    curl localhost:9000/projecta.tar?tree=HEAD
    curl localhost:9000/projecta/bar.tar?tree=HEAD
                                               
    pip install localhost:9000/projecta.tar?tree=origin/master
    
     
