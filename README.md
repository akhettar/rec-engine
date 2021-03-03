# Recommendation Engine

## Overview
A Simple Colleabartive filtering recommendation engine written in [Go](https://golang.org/) using [Redis](http://redis.io) as the In memory persistence store. The Redis client Go library used is [Redigo](https://github.com/garyburd/redigo).


## Usage

1. Install [httpie](https://httpie.io/): `brew install httpie`

2. Fire the below requests

```
======== CREATING SOME RATINGS ==========

 🔮  09:40:30   👸  ...workspace/vfc/rec-engine   main 💀 💥 🍄 🎀 
👉  http post 34.89.236.249:8765/api/rate user=user2 item=item1 score:=5.6
HTTP/1.1 201 Created
Content-Length: 43
Content-Type: application/json
Date: Tue, 02 Mar 2021 08:44:49 GMT

{
    "Item": "item1",
    "Score": 5.6,
    "user": "user2"
}


 🔮  09:44:49   👸  ...workspace/vfc/rec-engine   main 💀 💥 🍄 🎀 
👉  http post 34.89.236.249:8765/api/rate user=user2 item=item2 score:=5.6
HTTP/1.1 201 Created
Content-Length: 43
Content-Type: application/json
Date: Tue, 02 Mar 2021 08:44:59 GMT

{
    "Item": "item2",
    "Score": 5.6,
    "user": "user2"
}


 🔮  09:44:59   👸  ...workspace/vfc/rec-engine   main 💀 💥 🍄 🎀 
👉  http post 34.89.236.249:8765/api/rate user=user2 item=item3 score:=5.6
HTTP/1.1 201 Created
Content-Length: 43
Content-Type: application/json
Date: Tue, 02 Mar 2021 08:45:03 GMT

{
    "Item": "item3",
    "Score": 5.6,
    "user": "user2"
}


 🔮  09:45:03   👸  ...workspace/vfc/rec-engine   main 💀 💥 🍄 🎀 
👉  http post 34.89.236.249:8765/api/rate user=user2 item=item4 score:=8.6
^[[DHTTP/1.1 201 Created
Content-Length: 43
Content-Type: application/json
Date: Tue, 02 Mar 2021 08:45:10 GMT

{
    "Item": "item4",
    "Score": 8.6,
    "user": "user2"
}


 🔮  09:45:10   👸  ...workspace/vfc/rec-engine   main 💀 💥 🍄 🎀 
👉  http post 34.89.236.249:8765/api/rate user=user1 item=item1 score:=5.6
HTTP/1.1 201 Created
Content-Length: 43
Content-Type: application/json
Date: Tue, 02 Mar 2021 08:45:21 GMT

{
    "Item": "item1",
    "Score": 5.6,
    "user": "user1"
}

====== QUERYING SUGGESTIONS FOR A GIVEN USER =========

 🔮  09:45:22   👸  ...workspace/vfc/rec-engine   main 💀 💥 🍄 🎀 
👉  http get 34.89.236.249:8765/api/suggestion/user1                      
HTTP/1.1 200 OK
Content-Length: 88
Content-Type: application/json
Date: Tue, 02 Mar 2021 08:45:31 GMT

[
    {
        "item": "item4",
        "score": 8.6
    },
    {
        "item": "item3",
        "score": 5.6
    },
    {
        "item": "item2",
        "score": 5.6
    }
]

Getting probability for a given user with an item

 👉  http get 34.89.236.249:8765/api/probability/user1/item2
HTTP/1.1 200 OK
Content-Length: 49
Content-Type: application/json
Date: Tue, 02 Mar 2021 08:49:12 GMT

{
    "item": "item2",
    "propability": 5.6,
    "user": "user1"
}

