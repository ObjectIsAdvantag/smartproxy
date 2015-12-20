# Ambition  

Be a companion for Web API developers by mocking, filtering, replaying, "diff'ing" HTTP req/responses

It may also help Web API hosting via diff'ing between API versions and taking actions when errors.


# Architecture

SmartProxy acts as a reverse proxy that
 
- maintains an history of request
- allows to inspect them
- allows to modify them : YOU take action


# Roadmap

[x] reverse proxy basics 

   - custom path, custom port, healthcheck endpoint
   - tested on Windows, 6Mo exe
   - release tag : 0.1 
      
[ ] capture traffic

[ ] inspect traffic

[ ] extract API model

[ ] enrich model

[ ] admin api

[ ] extensibility


# Feeling like giving it a try

1. start the service you want to proxy, or launch the provides test service
 
   - "> go run tests/dummy.go" to start
   - check http://localhost:8080/ (/json /txt)

2. start smart-proxy : "go run main.go -route proxy"

3. try a few URLs

   - curl -X GET http://localhost:9090/
   - curl -X GET http://localhost:9090/proxy/
   - curl -X GET http://localhost:9090/proxy/json
   - curl -X GET http://localhost:9090/health


# License

BSD

Feel free to use, reuse, extend, and contribute







