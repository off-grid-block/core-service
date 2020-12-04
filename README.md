
***This repository is a work in progress.***

# DEON core-service
The Core Service of the DEON platform. This repository enables the administrative configuration and management of the DEON Hyperledger Fabric network. The service deploys the DEON platform REST API that exposes the configuration functions and the DEON suite of applications. The API is deployed locally on ```http://localhost:8001```. See documentation of the API at https://app.swaggerhub.com/apis/haniavis/deon-core/0.1.0.

## Setup

### Prerequisites

1. Docker Desktop (2.2.0.0)

To be added...

## Demo
For running a demo of the full DEON stack follow the instructions below.

### DEON Fabric network

1. Clone the DEON repository for the Fabric network `off-grid-block/off-grid-net`:

```git clone https://github.com/off-grid-block/off-grid-net.git```

2. Launch the Fabric network:

```cd off-grid-ne/t && ./cyfn.sh up -c off-channel -s couchdb```

### VON Network (Indy ledger)
The DEON services rely on VON Network, an implementation of a development level Indy Node network, developed by BCGov. For more information on the project and for additional instructions, see their [github repository](https://github.com/bcgov/von-network).

1. Clone the repository: ```git clone https://github.com/bcgov/von-network.git```
2. In the von-network repository: ```cd von-network/```
3. Generate the Docker images: ```./manage build```
4. Start up the network: ```./manage start```

### Aries Agents
The DEON services rely on the Aries Cloud Agent-Python (aca-py), developed by BCGov. For more information on the project and for additional instructions, see the official [github repository](https://github.com/hyperledger/aries-cloudagent-python). You need to clone the modified aca-py agents provided by DEON to enable Decentralized Identity Management.

1. clone the repository: ```git clone https://github.com/off-grid-block/aca-py-controllers.git```

### Launch using Docker
After launching the Fabric and VON networks, start up the DEON service API.
1. Clone this repository:
```git clone https://github.com/off-grid-block/core-service.git```
2. ```cd core-service```
3. ```export DOCKERHOST=`docker run --rm --net=host eclipse/che-ip````
4. ```docker-compose up```
5. access the API at ```localhost:8000/api/v1/```

To stop the networks and DEON service:
1. ```./manage down``` inside ```von-network``` directory
2. ```./cyfn.sh down``` inside ```off-grid-net``` directory
4. ```docker volume prune```
5. ```docker-compose down```
6. ```docker-compose rm -f```

### Register the DEON service on Indy

To test the demo, the first step is establishing a connection between the client and CI/MSP Aries Cloud Agents and creating a verifiable credential.
1. Access the client agent hosted at http://localhost:4201
2. Click the button labeled "Get invitation from Issuer agent"
3. Navigate to the CI/MSP agent at http://localhost:4200
4. On the sidebar, select "Schema and Credential definition" and create a schema with attributes "app_name, app_id" (name the schema whatever you like)
5. On the credential tab, issue a credential to the client agent.

Next, we will register the DEON vote application with the identity management agents. Send a POST request to http://localhost:8000/api/v1/register with the following body: `{
"Name": "Voting",
"Secret": "kerapwd",
"Type": "user"
}`

### Launch DEON apps (example: Vote Service)
1. clone the repository at ```github.com/off-grid-block/vote``` into ```deon```
2. ```docker-compose up```
3. See instructions on the full demo in the ```off-grid-block/vote``` repository

