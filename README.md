# Github-notify
A minimal application for sending GitHub notifications to the user through PushBullet.



### Prerequistes
- [docker](https://www.docker.com/)

### Setting up
Option 1: Pulling the docker image from [Docker hub](https://registry.hub.docker.com/u/erwinvaneyk/github-notify/)
```
docker pull erwinvaneyk/github-notify
```

Option 2: Building the docker image (replace version)
```
docker build --tag=erwinvaneyk/github-notify:latest . 
```

### Usage
Run the following commands (replacing the <key>'s)
```
docker run \
-e GITHUB_API_KEY=<key> \
-e PUSHBULLET_API_KEY=<key> \
-e CHECK_INTERVAL=10 \
--name=github-notify-instance \
-d \
erwinvaneyk/github-notify:latest
```
