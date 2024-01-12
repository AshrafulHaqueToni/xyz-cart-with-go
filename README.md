## **xyz-cart-with-go Project**

### **This is the basic flow of this project**

![Pasted image 1](https://github.com/AshrafulHaqueToni/xyz-cart-with-go/assets/48568933/74f487c3-dcd4-443f-8ea0-7fe761337f9c)

```bash
# You can start the project with below commands
docker-compose up -d
go run main.go
```
#### **SIGNUP FUNCTION API CALL (POST REQUEST)** ####
##### For the middleware you have to set a valid token #####
http://localhost:8080/users/signup
```json
{
  "firstName": "Ashraful",
  "lastName": "Toni",
  "email": "ashrafulhaquetoni@gmail.com",
  "password": "demo1234",
  "phone": "016000000000"
}
```
Response: "Successfully Signed Up!!"

#### **LOGIN FUNCTION API CALL (POST REQUEST)** ####

http://localhost:8080/users/login

```json
{
  "email": "ashrafulhaquetoni@gmail.com",
  "password": "demo1234"
}
```
response will be like this

```json
{
    "_id": "65a1752e7fa5e80e1c984935",
    "firstName": "Ashraful",
    "lastName": "Toni",
    "password": "$2a$14$wOA.Vpxvyo5hOUt7WpAvSuNRd46NGsoKIfjnITZhxNzQKlJAUpr3m",
    "email": "ashrafulhaquetoni@gmail.com",
    "phone": "016000000000",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6ImFzaHJhZnVsaGFxdWV0b25pQGdtYWlsLmNvbSIsIkZpcnN0TmFtZSI6IkFzaHJhZnVsIiwiTGFzdE5hbWUiOiJUb25pIiwiVWlkIjoiNjVhMTc1MmU3ZmE1ZTgwZTFjOTg0OTM1IiwiZXhwIjoxNzA1MTY2NTEwfQ.a40BjNjgAFoj3hUEtTkFhTkZeD5oieHnzJNyFNYXo-k",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6IiIsIkZpcnN0TmFtZSI6IiIsIkxhc3ROYW1lIjoiIiwiVWlkIjoiIiwiZXhwIjoxNzA1Njg0OTEwfQ.e-dGlvde0vwYu6Ck94q14i01tNongLAI2OxcjcxV0rM",
    "created_at": "2024-01-12T17:21:50Z",
    "updated_at": "2024-01-12T17:21:50Z",
    "user_id": "65a1752e7fa5e80e1c984935",
    "user_cart": [],
    "address_details": [],
    "order_status": []
}
```

#### **Admin add Product Function (POST REQUEST)** ####

http://localhost:8080/admin/addproduct

```json
{
  "product_name": "gowatch v1",
  "price": 12000,
  "rating": 5,
  "image": "gowatch.jpg"
}
```

Response : "Successfully added our Product Admin!!"

#### **View all the Products in db GET REQUEST** ####

http://localhost:8080/users/productview

Response

```json
[
    {
        "_id": "65a177127fa5e80e1c98493a",
        "product_name": "gowatch v1",
        "price": 12000,
        "rating": 5,
        "image": "gowatch.jpg"
    },
    {
        "_id": "65a177a67fa5e80e1c98493d",
        "product_name": "gomobile v1",
        "price": 10000,
        "rating": 4,
        "image": "gomobile.jpg"
    }
]
```

#### **Search Product by regex function (GET REQUEST)** ####

defines the word search sorting

http://localhost:8080/users/search?name=go

response:

```json
[
    {
        "_id": "65a177127fa5e80e1c98493a",
        "product_name": "gowatch v1",
        "price": 12000,
        "rating": 5,
        "image": "gowatch.jpg"
    },
    {
        "_id": "65a177a67fa5e80e1c98493d",
        "product_name": "gomobile v1",
        "price": 10000,
        "rating": 4,
        "image": "gomobile.jpg"
    }
]
```

- **Adding the Products to the Cart (GET REQUEST)**

  http://localhost:8000/addtocart?id=xxxproduct_idxxx&userID=xxxxxxuser_idxxxxxx

  Corresponding mongodb query

- **Removing Item From the Cart (GET REQUEST)**

  http://localhost:8000/addtocart?id=xxxxxxx&userID=xxxxxxxxxxxx

- **Listing the item in the users cart (GET REQUEST) and total price**

  http://localhost:8000/listcart?id=xxxxxxuser_idxxxxxxxxxx

- **Addding the Address (POST REQUEST)**

  http://localhost:8000/addadress?id=user_id**\*\***\***\*\***

  The Address array is limited to two values home and work address more than two address is not acceptable

```json
{
  "house_name": "white house",
  "street_name": "white street",
  "city_name": "washington",
  "pin_code": "332423432"
}
```

- **Editing the Home Address(PUT REQUEST)**

  http://localhost:8000/edithomeaddress?id=xxxxxxxxxxuser_idxxxxxxxxxxxxxxx

- **Editing the Work Address(PUT REQUEST)**

  http://localhost:8000/editworkaddress?id=xxxxxxxxxxuser_idxxxxxxxxxxxxxxx

- **Delete Addresses(GET REQUEST)**

  http://localhost:8000/deleteaddresses?id=xxxxxxxxxuser_idxxxxxxxxxxxxx

  delete both addresses

- **Cart Checkout Function and placing the order(GET REQUEST)**

  After placing the order the items have to be deleted from cart functonality added

  http://localhost:8000?id=xxuser_idxxx

- **Instantly Buying the Products(GET REQUEST)**
  http://localhost:8000?userid=xxuser_idxxx&pid=xxxxproduct_idxxxx