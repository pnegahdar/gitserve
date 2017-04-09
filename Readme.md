## gitserve: webserver to serve git archives of any git path  

Creates a webserver that serves on top of git-archive, fetches the repo in the background every 30 seconds. 

#### Installation:

Grab the right precompiled bin from github releases and put it in your path. Don't forget to `chmod +x` the bin.

OSX:

```bash
curl -SL https://github.com/pnegahdar/gitserve/releases/download/0.2.0/darwin_amd64 > /usr/local/bin/gitserve && \
    chmod +x /usr/local/bin/gitserve 
```

Nix:

```bash
curl -SL https://github.com/pnegahdar/gitserve/releases/download/0.2.0/linux_amd64 > /usr/local/bin/gitserve && \
    chmod +x /usr/local/bin/gitserve 
```

#### Usage:

##### Run the server

```bash
# Minimal 
gitserver 

# All args  
gitserver -root=/Users/myuser/git/project  -prefix=subdira/ -listen=":9000"

# Pypi Server
# All packages are in the ${root}/{prefix}/ where package_name = <org>_<name>. Git tags with are ${package_name}V${version} to declare versions
gitserver -root=/Users/myuser/git/project  -prefix=subdira/ -listen=":9000" -pypi-tag-prefix org_ -pypi-tag-delimiter V
```
    
     
#### Get the archives:

    curl localhost:9000/.zip 
    curl localhost:9000/.tar 
    curl localhost:9000/projecta.tar?tree=HEAD
    curl localhost:9000/projecta/bar.tar?tree=HEAD
                                               
    pip install localhost:9000/projecta.tar?tree=origin/master
    
    # If pypi server is enabled
    pip install --extra-index-url http://127.0.0.1:8020/_pypi-simple package_name==<version>
    
     
