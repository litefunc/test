
### build arm64 docker image

```
docker build -t web:latest .
```

### build amd64 docker image

```
docker build -t web:latest -f Dockerfile_amd64 .
```

### save image to tar

```
docker save -o web.tar web:latest
```


### run web service on docker

```
docker run --name web -p 8888:8888 web:latest
```
open browser at http://localhost:8888/


### run web service on docker in background

```
docker run -d --name web -p 8888:8888 web:latest
```
open browser at http://localhost:8888/