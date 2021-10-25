# Bluesnap Golang SDK
This was written for the Revel Framework
This is not actively maintained, but can be updated if sponsered to do so

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

# Example Code
Example for getting a users cards
```go
if len(account.BSShopperId) == 0 {
		return c.RenderJson(Cards{Message: getSuccessMessage("No cards on file."), Cards: &bluesnap.CardList{
			Values: make([]*bluesnap.Card, 0),
		}})

} else {
		b := &bluesnap.BlueSnap{}
		b.Init()
		b.Test = !isProd
		cList, err := b.GetCards(account.BSShopperId)
		if err != nil {
			return c.RenderJson(Cards{Message: getSuccessMessage("Error getting cards on file."), Cards: &bluesnap.CardList{
				Values: make([]*bluesnap.Card, 0),
			}})
		}
		return c.RenderJson(Cards{Message: getSuccessMessage("Cards on file"), Cards: cList})
}
```
