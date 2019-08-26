# supermanDetector
## Service for detecting suspicious IP Access.

The service accepts a JSON POST request with data about an IP Access with the following structure:
```
{"username": "bob", "unix_timestamp": 1514764006, "event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e46", "ip_address": "203.2.218.214"}
```
An SQLite database is used to store the IP Access along with the coordinates of the IP Access. The database contains the following fields. 

#### UUID identifying the login attempt
#### Username
#### Unix epoch timestamp in seconds
#### IPv4 address
#### Latitude
#### Longitude
#### Location Radius 

When a new IP Access is posted to the API, a database query for previous and subsequent IP Accesses by the user executes. The query results go into a JSON _`IP Access Report`_ containing data about the current, previous, and subsequent IP Access. The report also contains the best guess about whether the previous or subsequent IP Access is suspicious. The _`IP Access Report`_ is returned to the API caller and has the following structure:
```
{"currentGeo":{"lat":39.211,"lon":-76.8362,"radius":5},"travelToCurrentGeoSuspicious":false,"travelFromCurrentGeoSuspicious":true,"precedingIpAccess":{"ip":"","speed":0,"lat":0,"lon":0,"radius":0,"timestamp":0},"subsequentIpAccess":{"ip":"203.2.218.214","speed":11331113,"lat":-33.8919,"lon":151.1554,"radius":500,"timestamp":1514764006}}
```

If a user would need to travel over 500 km an hour to each destination to sign in - the IP Access is considered suspicious. The service can check the "travel speed" of a user by examining the timestamp and coordinates of the previous and subsequent IP Accesses.

#### The service relies on the following dependencies:

_For the Object Relational Mapper_

github.com/jinzhu/gorm v1.9.10 

_For the Geo IP coordinate lookups_

github.com/oschwald/geoip2-golang v1.3.0
github.com/oschwald/maxminddb-golang v1.3.1 // indirect

_For calculating the distance between coordinates_ 

github.com/umahmood/haversine v0.0.0-20151105152445-808ab04add26

## To run the service, use one of two options.

### Option One - Clone the repository and build the binary 
Clone or download this repository into you Go workspace
```
$ git clone https://github.com/NEPDAVE/supermanDetector.git
```
`cd` into the `supermanDetector` directory, build the binary, and run the program
```
$ cd supermanDetector
$ go build
$ ./supermanDetector
```
To create an IP Access report run `curl -X POST -d '{"username": "bob", "unix_timestamp": 1514764007, "event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e47", "ip_address": "203.2.218.214"}' http://localhost:5000/v1/ipaccess`

### Option Two - Use Docker 
```
$ docker pull dstreck/superman-detector
$ docker run -p 5000:5000 dstreck/superman-detector:latest
$ curl -X POST -d '{"username": "bob", "unix_timestamp": 1514764007, "event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e47", "ip_address": "203.2.218.214"}' http://localhost:5000/v1/ipaccess`
```

Addtional Curl commands for manual testing:
```
Hong_Kong 
curl -X POST -d '{"username": "bob", "unix_timestamp": 1514764000, "event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e40", "ip_address": "119.28.48.231"}' http://localhost:5000/v1/ipaccess

New York 
curl -X POST -d '{"username": "bob", "unix_timestamp": 1514764001, "event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e41", "ip_address": "206.81.252.6"}' http://localhost:5000/v1/ipaccess

Moscow 
curl -X POST -d '{"username": "bob", "unix_timestamp": 1514764002, "event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e42", "ip_address": "31.173.221.5"}' http://localhost:5000/v1/ipaccess

Sydney - 6,4
curl -X POST -d '{"username": "bob", "unix_timestamp": 1514764006, "event_uuid": "85ad929a-db03-4bf4-9541-8f728fa12e46", "ip_address": "203.2.218.214"}' http://localhost:5000/v1/ipaccess

```
