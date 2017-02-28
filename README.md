# WAPIQ (Web Application Interface Query)

## What?
So, let's say you need to interface with all kinds of different web APIs, for example:
  * Google Places API
  * TripAdviser API
  * Yelp API

You, for whatever reason beyond general laziness can't be bothered to write yet another web api interaction, and properly parse the damned thing... So here's my solution to this, or at least some attempt.

## How?

Example: Google Places Configuration
```wapiq
"GooglePlaces" API {
  path `https://maps.googleapis.com/maps/api/place/`
  args {
    "key" `YOUR_API_KEY`
  }
};
"search" GET {
  path `nearbysearch/json`
  type `json`
  head []
  query [
    `location`
    `radius`
    `types`
    `name`
  ]
  body []
};
"Place" MAP "GooglePlaces" {
  "search" {
    "id"        `results.place_id`
    "name"      `results.name`
    "types"     `results.types`
    "location"  `results.geometry.location`
    "address"   `results.vicinity`
  }
};
```

**Explain what just happened?**

* **"api name" API { ... }:** Sets up a new API with the quoted name (In this case *GooglePlaces* ).
  * **path** Sets the APIs base uri to use when making any requests.
  * **args** Sets constants that can be used for any request, handy for API keys
    * **"key" \`YOUR_API_KEY\`** : Declares a new constant named *"key"* with the value *YOUR_API_KEY* .
* **"action name" GET { ... }:** Sets up a new API action, in this case a HTML GET request called *search*.
  * **path** : Sets the action's path to be appended after the APIs path.
  * **type** : Sets the action's return type, in this case *json* we expect output.
  * **head** : (Optional) Sets the action's possible header parameters to be sent, in this case none are set.
  * **query** : (Optional) Sets the action's possible query parameters to be appended after the path (URL encoded with ?,&,=).
  * **body** : (Optional) Sets the action's possible POST or body parameters to be sent, in this case none are set.
* **"Place" MAP { ... }:** Sets up a new API map, in this case one called *Place*
  * **"search" { ... }:** Sets up an action mapping for our search action, which defines:
    * **"id" \`results.place_id\`** Maps the *Place* field *id* to the json location \`results.place_id\`

**JSON Locations**
As previously mentioned, a json mapping is defined by the route to parse out the value you want to have serialized in your object. However, you *do not* need to reference array indexes, as WAPIQ maps all queried values (ie. we just care about the absolute path, not the specific object index returned from the request.)

Currently WAPIQ can only be called from Golang, so if you do not want to override default mappings shown above, you can configure a specific struct with direct mappings if you only want to map to a single predictable request mapping.

**Golang Mapping Example**
```go
type Location struct {
  Lat float32       `json:"lat" wapiq:"results.geometry.location.lat"`
  Lon float32       `json:"lon" wapiq:"results.geometry.location.lon"`
}
type Place struct {
  ID int64          `json:"id" wapiq:"results.place_id"`
  Name string       `json:"name" wapiq:"results.name"`
  Types []string    `json:"types" wapiq:"results.types"`
  Location Location `json:"location" wapiq:"results.geometry.location"`
  Address string    `json:"address" wapiq:"results.vicinity"`
}
```

### Ok, so how do I request data now?

```wapiq
/search FOR Place WHERE
  name `cruise`
  location `-33.8670,151.1957`
  radius `500`
  types `food`
  ;
```

So if you've ever used a query language, its very similar, as expected, this query will return a `[]Place` with the following criteria from the default mapping supplied.

Fires a request behind the scenes:
> https://maps.googleapis.com/maps/api/place/nearbysearch/json?key=AIzaSyCZmDlZXIlhlkDbHzAfffvWGWQa1LliZvE&location=-33.8670%2C151.1957&name=cruise&radius=500&types=food

Which returns (Console verbose print from Golang):
```
[map[
  id:ChIJrTLr-GyuEmsRBfy61i59si0
  name:Australian Cruise Group
  types:[
    travel_agency
    restaurant
    food
    point_of_interest
    establishment
  ]
  location:map[lat:-33.867591 lng:151.201196] address:32 The Promenade, King Street Wharf 5, Sydney
]
map[
  id:ChIJLfySpTOuEmsRMFymbMkVkOE
  name:Vagabond Cruises
  types:[
    travel_agency
    restaurant
    food
    point_of_interest
    establishment
  ]
  location:map[lat:-33.87104699999999 lng:151.189736]
  address:37 Bank Street, Pyrmont
]
map[...]
...
```

Read more about WAPIQ [here](http://simonwillshire.com/blog/wapiq/)
