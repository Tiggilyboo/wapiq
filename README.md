# WAPIQ (Web Application Interface Query)

## What?
So, let's say you need to interface with all kinds of different web APIs, for example:
  * Google Places API
  * TripAdviser API
  * Yelp API
  * A Crypto currency exchange API (ex. Bitfinex)

You, for whatever reason beyond general laziness can't be bothered to write yet another web api interaction, and properly parse the damned thing... So here's my solution to this, or at least some attempt.

## How?

### Set up
```sh
git clone https://github.com/Tiggilyboo/wapiq
cd wapiq
go build
./wapiq -f ./examples/GooglePlaces.wapiq
```
Returns JSON from CLI
```json
{  
   "0":[  
      {  
         "address":"32 The Promenade, King Street Wharf 5, Sydney",
         "id":"ChIJrTLr-GyuEmsRBfy61i59si0",
         "location":{  
            "lat":-33.867591,
            "lng":151.201196
         },
         "name":"Australian Cruise Group",
         "types":[  
            "travel_agency",
            "restaurant",
            "food",
            "point_of_interest",
            "establishment"
         ]
      },
      {  
        ...
      }
   ]
}
```

### Ok, what just got run?
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

Read more about WAPIQ on my site: [here](http://simonwillshire.com/blog/WAPIQ/)

### Example Configuration

Here is the simple GooglePlaces API example we just ran:
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

**Explain, what just happened?**

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

### Request URL Variables

Many APIs use variables within a request URL, WAPIQ is capable of using variable output in the URL by wrapping  `{}` braces around them. An example use case is provided in the `Yelp.wapiq` example, but here it is in short:

```wapiq

"GetReviews" GET {
  path `/v3/businesses/{id}/reviews`
  head [
    `access_token`
  ]
  query [
    `locale`
  ]
};
```
**Note:** ***These can be used in the API path, or request (GET/POST) paths***

In the above, `{id}` is the placeholder variable which gets replaced if the variable is provided in the query:

```wapiq
/GetReviews FOR Reviews WHERE`
  id `SOME_BUSINESS_ID`
;
```

### Array @ Expressions

Some APIs use arrays with an expected set of information (without naming the key of the value in the pair). In order to handle these responses, WAPIQ uses the `@` character when specifying the index in the array, or the depth of the array. Here is an example that is used when handling responses from the Bitfinex exchange API (See more in ```examples/Bitfinex.wapiq```):

```wapiq
...
"Trade" MAP "Bitfinex" {
  "trades" {
    "Id"              @0,0
    "Mts"             @0,1
    "Amount"          @0,2
    "Price"           @0,3
  }
  "funding" {
    "Id"              @0,0
    "Mts"             @0,1
    "Amount"          @0,2
    "Rate"            @0,3
    "Period"          @0,4
  }
};
```

In the above example, ```@0,0``` is used to read a json value from a response like this one, where <> are used as placeholders to identify the indexes:

```json
  [
    [
      <ID>,
      <MTS>,
      <AMOUNT>,
      <PRICE>
    ],
    [
      ...
    ],
    ...
  ]
```

The first ```@0``` specifies the first embedded array object, ```[[<here>],...]```, and the second identifies the first element (at index 0) inside that array.

If we wanted to just read the first element in an array, we just use ```@0```, and ```@1``` is the proceding element, etc.

### Includes
So let's say you're implementing a particularily intricate API using WAPIQ, but you are struggling with the separation of concerns of the script. For example, you are implementing a public or common portion of the API, and want to separate the authenticated endpoints from the common. In order to handle this, WAPIQ can use include statements using the `^` character followed by the filename (excluding the wapiq extension).

As an example, the Bitfinex example is split into Bitfinex.wapiq and Bitfinex_Auth.wapiq:
```wapiq
# Include common API from Bitfinex.wapiq
^Bitfinex

"wallets" POST {
  path "auth/wallets"
  head [
    `bfx-apikey`
    `bfx-apisecret`
    `bfx-nonce`
  ]
};
```

At the moment WAPIQ does not support inheritance of objects (potential in the future), but this at least allows file separation.
