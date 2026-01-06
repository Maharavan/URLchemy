# url-shortner


1. understand http.resposne and http.request
2. base62 possible way
3. get post method (handlefunc)
4. redis for caching
5. docker

POST /longurl
 ├─ Read long URL from request body
 ├─ Generate short code
 ├─ Store: map[shortCode] = longURL
 └─ Return short URL

GET /{shortCode}
 ├─ Extract shortCode from path
 ├─ Lookup in map
 └─ Redirect


