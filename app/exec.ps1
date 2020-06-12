docker build --rm -t guard-my-app .
docker run --init -p 3000:3000 -it guard-my-app
