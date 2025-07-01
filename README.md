# GoChain

To build single container:
<br>```docker build -t node .```

To run that container and view logs in the same terminal:
<br>```docker run -it --rm -p 8001:8001 node```


To run multiple containers together:
<br>```docker compose up```

To restart multiple containers:
<br>```docker compose restart```

To shutdown multiple containers:
<br>```docker compose down```

or if you want to destroy the volumes also:
<br>```docker compose down -v```

To run tests:
<br>```go test GoChain/block```
