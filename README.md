***This repository is a work in progress.***

# DEON core-service
The Core Service of the DEON platform. This repository allows the adminstrative configuration and management of the DEON Hyperledger Fabric network. The service deploys the DEON platform REST API that exposes the configuration functions and the DEON suite of applications. The API is deployed locally on localhost:8001.

## Setup

### Hyperledger Fabric

1. Clone the Hyperledger `fabric-samples` repository:
```git clone https://github.com/hyperledger/fabric-samples.git```
2. Inside `fabric-samples`, checkout to version 1.4.2:
```git checkout 1.4.2```
3. Replace the `first-network/byfn.sh` script with the script found [here](https://github.com/off-grid-block/off-grid-net/blob/master/cyfn.sh).
4. Download the Hyperledger Fabric v1.4.2 docker images:
```curl -sSL https://bit.ly/2ysbOFE | bash -s -- 1.4.2```

### Launch using Docker

Start up the Fabric network:
1. ```cd fabric-samples/first-network``` (inside your fabric-samples repository)
2. ```./byfn.sh up -s couchdb```

Start up the DEON service API:
1. ```mkdir deon```
2. clone this repository into the ```deon``` directory
3. mv the ```docker-compose.yaml``` file included in this repository into the parent ```deon``` directory
4. clone the repository at ```github.com/off-grid-block/vote``` into ```deon```
5. ```docker-compose up``` at ```deon``` directory level
6. access the API at ```localhost:8001/api/v1/```

To stop the network and DEON service:
1. ```./byfn.sh down``` inside ```fabric-samples/first-network```
2. ```docker-compose down```