# GoChain

To build the container:
docker build -t node .

To run the container and use terminal to see logs:
docker run -it --rm -p 8001:8001 my-blockchain-node


To run multiple containers:
docker compose up

To restart every container:
docker compose restart

To shutdown containers:
docker compose down

If you want to destroy the volumes:
docker compose down -v
