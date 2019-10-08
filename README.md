#Zeta Storage

#Language
Golang

#Database
LevelDB

#How to run
Go to zeta-st, execute the following command:
    go run .

#Node initial connection to DISCOVERY_NODE
1. Connects to discovery node in congif/<environment>.yml
2. Send this node information
3. Receive dn message. The message is a history of nodes connected to DISCOVERY_NODe
4. close connection

