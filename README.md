# GoChain

GoChain is a study project aimed at learning decentralised networking principles. It uses a simplified Proof of Work(PoW) consensus mechanism and a gossip protocol to dynamically discover other nodes.


To run the whole chain:
<br>```docker compose up```

To restart the chain:
<br>```docker compose restart```

To shutdown the chain:
<br>```docker compose down```

or if you want to destroy the container and images also:
<br>```docker compose down --rmi all -v```

To run tests:
<br>```go test GoChain/block```
