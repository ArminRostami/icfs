# ICFS
**ICFS** (intracampus file-sharing) is a distributed file-sharing solution for private networks with a focus on speed and ease of use.

This repository inlcudes the client node that comes with a standalone [ipfs node](https://ipfs.io/) and a shell interface to interact with the network.

## Architecture
There are two types of node in each icfs network:
1. **Client** nodes: each client runs its own ipfs daemon and connects to peers in order to share files

2. **Bootstrap** or discovery nodes: client nodes query these nodes to find other client nodes
## Build and Run
(bootstrap node should be configured first)

1. cd into the repo dirtectory
2. build the docker image: `docker build -t icfs .`
3. start the container: `docker run --rm -it --name c1 icfs`
4. follow the on-screen instructions

now you can connect to container directly via bash to see changes made after running commands: 

`docker exec -it c1 /bin/bash`