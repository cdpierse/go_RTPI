# go_RTPI

 ## *In Development*

 _`go_RTPI`_ is simple wrapper developed in Go to expose the RTPI  (Real Time Passenger Information) REST web services API endpoints, and provide a simple interface for Go developers to query this service in Go. Additionally this service implements a number of new services built on top of the existing endpoints to provide users with more powerful queries. 

 


## Stop Endpoints
 - `/stops` - Type `GET` - returns all stops from RTPI 
 - `/stops/id/{stop_id}` - Type `GET` - returns stop with corresponding {stop_id}
 - `/stops/stop_name/{stop_name}` - Type `GET` - returns stop(s) that are an exact match {stop_name} on either the stops's full name or short name. 
 -  `/stops/stop_name_fuzzy/{stop_name}` - Type `GET` - returns stop(s) that are a fuzzy match with {stop_name} on either the stops's full name or short name. 
 - `/stops/operator/{operator_name}` - Type `GET` - returns stop(s) that match with the provided {operator_name} e.g. bac (Dublin Bus), BE (Bus Eireann). 
 - `/stops/distance` - Type `GET` takes mandatory `latitude` and `longitude` query params as well as an optional `max_distance` query param which defaults to 500 meters. Endpoint returns all stops within max_distance of the latitude and longitude location specified.
 - `/stops/distance/{stop_id}` - Type `GET` takes an optional `max_distance` query param which defaults to 500 meters. Endpoint returns all stops within max_distance of the stop that matches {stop_id}