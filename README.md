# bluesnap-api-go
Golang Bluesnap API SDK
This was written for the Revel Framework for Go

# Setup
* Edit bluesnap/api.go lines 6 and 61. Change myApp to your app name
* Add OS Environment Variables for the following:
  * TEST_BS_API_USER
  * TEST_BS_PASS
  * TEST_BS_STORE_ID
  * TEST_BS_ENDPOINT
  * TEST_BS_PRODUCT_ID
  * BS_API_USER
  * BS_PASS
  * BS_STORE_ID
  * BS_ENDPOINT
  * BS_PRODUCT_ID

# Instructions
Include it where needed


```go
import(
  "myApp/bluesnap"
)
```
