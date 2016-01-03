# Ambition  

Be a companion for Web API developers by mocking, filtering, replaying, "diff'ing" HTTP req/responses

It may also help Web API hosting via diff'ing between API versions and taking actions when unexpected errors occur. 
Note : by unexpected, we intend that traffic differs from prevailing traffic.


# Architecture

SmartProxy acts as a reverse proxy that
 
- maintains an history of requests
- allows to inspect them
- allows to modify them : YOU take action


# Roadmap

[X] reverse proxy basics 

   - custom path, custom port, healthcheck endpoint
   - tested on Windows, 6,9 MB windows exe
   - tested on Linux / Docker, 6,9MB image (from scratch)
   - release tag : 0.1 
      
[X] capture traffic, writes to stdout
   - MIT license
   - dump flag 
   - captures requests/responses, 
   - prints to stdout only (no in-memory, no disk writes)
   - release tag v0.2
   
[X] store traffic on-disk
   - tag each flow (req/resp) with a uuid
   - create go struct to store traffic : ingress, egress, synthesis
   - store traffic via bolt
   - release tag v0.3
    
[X] inspect traffic
   - added route to inspect captured traffic, /traffic  
   - display list of captured traffic in HTML page (no json, 20 latest traces) 
   - display detailed trace in json
   - release tag v0.4
   
[ ] extract API model from traffic analysis

[ ] enrich model

[ ] Realtime traffic capture experience

[ ] admin api
   - dynamically modify configuration (change proxy address, capture on/off, erase capture database)
   

# Feel like giving it a try

## Easy does it : start from an executable

Go to github / releases,

Pick the executables that suits your platform : smartproxy and dummy test service

Run it

```
# on linux
> smart-proxy -route proxy -port 9090 -serve 127.0.0.1:8080 -capture
> tests/dummy
```

```
# on windows
> ./smart-proxy.exe -route proxy -port 9090 -serve 127.0.0.1:8080 -capture
> tests/dummy.exe
```

Then inspect traffic via your web brower, 

- http://localhost:9090/traffic : display the latest 10 traces
- http://localhost:9090/traffic/<trace_id> : display trace details in json

 
## You gopher (have golang installed locally)

1. start the service you want to proxy, or launch the provides test service
   - start with "> go run tests/dummy.go" 
   - check "> curl -X GET http://localhost:8080/" and also /json /txt
2. start smart-proxy : "> go run *.go -route proxy"
   - 2015/12/20 19:25:49 [INFO] Starting SmartProxy, version: draft
   - 2015/12/20 19:25:49 [INFO] Listening on http://localhost:9090
   - 2015/12/20 19:25:49 [INFO] Serving http://127.0.0.1:8080 via path /proxy/
3. try a few URLs
   - curl -X GET http://localhost:9090/
   - curl -X GET http://localhost:9090/proxy/
   - curl -X GET http://localhost:9090/proxy/json
   - curl -X GET http://localhost:9090/ping

## Via Docker

### On linux (pure assumption, need to be checked, please contribute :-))

- Clone this repo
- Pick the linux executable from releases or build it with command : "> make linux"
- Build the docker image with the provided Dockerfile, "> docker build  -t objectisadvantag/smartproxy ."
- Run it "> docker run -d -p 9090:9090 objectisadvantag/smartproxy:latest -route proxy -port 9090 -serve 127.0.0.1:8080 -capture"

```
> make linux
> docker build -t <image-name> .
> docker run -d -p 9090:9090 <image-name> -route proxy -port 9090 -serve 127.0.0.1:8080 -capture
```


### Docker Toolbox 

Note for Windows gophers : I personally use an home-made set of bash commands : 
check https://github.com/ObjectIsAdvantag/my-docker-toolbox

```
> git clone https://github.com/ObjectIsAdvantag/smart-proxy
> make linux
> dmcreate docker-smartproxy    // creates a new box
> dminit                        // initializes the toolbox env
> dmip                          // display your box ip address, simply replace suffix with 1 to get the address where your containers can reach your dev machine
> dibuild smart-proxy .         // builds the image
> dimg 1                        // selects the image that has just been created as current
> drun                          // creates and launches a new container
command ?   [command|(default)]: -route proxy -serve 192.168.99.1:8080
detach or interactive ? [d/(i)]: i
expose ports ? HOST:CONTAINER : 9090:9090
docker run -it -p 9090:9090 smartproxy:latest -route proxy -serve 192.168.99.1:8080 ? [(y)|n] :
Launching...
```
   
# License

MIT, see license file.

Feel free to use, reuse, extend, and contribute







