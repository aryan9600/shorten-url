## Shorten URL

This is a short project that I made to explore GoLang. The idea was to create my own bit.ly using GoLang.
It uses BoltDB buckets to store tokens and corresponding URLs in a map like {token: URL}.
After registering your token and URL, simply go to ``www.example.com/yourToken`` and you will be redirected to your desired URL!

**Usage:**

``git clone https://github.com/aryan9600/shorten-url.git``

``cd shorten-url``

``go run main.go``
   
| **Endpoint**      | **Method** | **Body** |
| ----------- | ----------- |------|
| /api/v1/shorten      | POST       | {"path: "token", "url": "URL"} |
| /api/v1/shorten   | GET, DELETE        |{"path": "token"}|



