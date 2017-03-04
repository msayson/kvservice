# kvservice
A simple key-value service, with variations demonstrating different failure recovery strategies.

Client command line interface:

<table>
  <td>Command</td><td>Description</td>
  <tr><td>get(id)</td><td>returns value for id</td></tr>
  <tr><td>set(id,val)</td><td>sets value for id</td></tr>
  <tr><td>testset(id,testVal,newVal)</td><td>if id has testVal as its value, set to newVal</td></tr>
  <tr><td>exit</td><td>shuts down client</td></tr>
</table>

### Variation 1 - single server
A simple client/server system in which clients can send requests to read/write key-values.

![alt-text](https://github.com/msayson/kvservice/wiki/design_mockups/images/variation1singleserver.png "Diagram of single-server key-value service")

### Variation 2 - dynamic chain of servers (in progress)
A simple key-value service with data replication across a chain of N back-end servers.

- Client interacts with the front-end server exactly as in Variation 1
- A chain of N back-end servers store identical copies of all key-values
- Key-value write operations are performed on each back-end server, passed from one to the next until all are updated
- Key-value read operations are performed on the first back-end server in the chain
- Back-end nodes may join or leave the network at any time

Failure recovery strategy:

- Each server is aware of the next two nodes in the chain
- If one back-end server fails, its predecessor links to the next known node to reconnect the chain

Design properties:

- Robust to one back-end failing at a time, up to a maximum of N-1 back-end failures
- Not robust to the front-end server failing
- If two adjacent back-end nodes fail simultaneously, then subsequent back-end nodes will be stranded.  However, if at least one of the first two back-end nodes survives the failure, we can continue operating without a break in service.

### Disclaimer

This project was developed for educational purposes, and comes without warrantee or support.  However, feel free to copy and modify its code and ideas as you wish.
