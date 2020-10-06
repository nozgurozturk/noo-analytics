<p align="center">
  <a href="https://github.com/nozgurozturk/mercury-ui">
    <img style="-webkit-user-select: none; display: block; margin: auto; padding: env(safe-area-inset-top) env(safe-area-inset-right) env(safe-area-inset-bottom) env(safe-area-inset-left); cursor: zoom-in;" src="https://user-images.githubusercontent.com/44316679/95265189-64122180-0839-11eb-9b76-e8fd61da087b.png" width=" 398" height="130">
  </a>

  <h3 align="center">noo-analytics</h3>

  <p align="center">
    Simple backend service for personal websites
  </p>
</p>

## Install
First of all, [download](https://golang.org/dl/) and install Go. `1.14` or higher is required.

When installation is done clone the repo with command:

```bash
git clone https://github.com/nozgurozturk/noo-analytics
 ```
After that install missing dependencies with command:
```bash
go mod tidy
```

## Run the app

    go run ./cmd/server

## Database Support

noo-analytics works with **MongoDB** and **Redis**

## Environment Variables

**MongoDB Variables:**
```.env
MONGO_DB_USERNAME = 
MONGO_DB_PASSWORD = 
MONGO_DB_HOST = 
MONGO_DB_PORT = 
MONGO_DB_NAME = 
MONGO_DB_QUERY = 
```
**Redis Variables:**
```.env
REDIS_DB_HOST = 
```
**Server Variables:**
```.env
SERVER_PORT = 
ACCESS_SECRET = 
REFRESH_SECRET = 
# Hour
ACCESS_EXPIRE =
# Minute
REFRESH_EXPIRE = 
ADMIN_MAIL = 
```

## API Documentation

### Send User Trace 

#### Request

`POST /trace`
    
```json
{
  "ip": "192.168.0.1",
  "loc": "TUR",
  "a": "Mozilla/5.0 (Macintosh; Intel Mac OS X x.y; rv:42.0) Gecko/20100101 Firefox/42.0",
  "act": 1,
  "tag": "home"
}   
```    

#### Response
```json
    {}
```
    
### Find Action by Date 

#### Request

- *year, month, day and hour are objects are optional*
- *range and its properties are optional*

`POST /analytics/action/date`
    
```json
{
    "action": 1,
    "year":{
        "isInclude": true,
        "range": {
          "from": 2018,
          "to": 2020
        }
    },
    "month":{
        "isInclude": true,
        "range": {
          "from": 2018,
          "to": 2020
        }
    },
    "day":{
        "isInclude": true,
        "range": {
          "from": 2018,
          "to": 2020
        }
    },
    "hour":{
        "isInclude": true,
        "range": {
          "from": 2018,
          "to": 2020
        }
    }
}
```    

#### Response
**y**: year |
**m**: month |
**d**: day |
**h**: hour |
**v**: visitors' ip |
**uv**: unique visitors' ip

```json
[
    {
        "y": 2020,
        "m": 11,
        "d": 6,
        "h": 0,
        "v": [
            "192.168.0.1"
        ],
        "uv": [
            "192.168.0.1"
        ]
    },
    {
        "y": 2020,
        "m": 10,
        "d": 6,
        "h": 0,
        "v": [
            "192.168.0.1"
        ],
        "uv": [
            "192.168.0.1"
        ]
    },
    {
        "y": 2020,
        "m": 10,
        "d": 4,
        "h": 0,
        "v": [
            "192.168.0.1",
            "192.168.0.2",
            "192.168.0.1"
        ],
        "uv": [
            "192.168.0.1",
            "192.168.0.2"
        ]
    }
]
```

### Login 

#### Request

`POST /auth/login`
    
```json
{
  "email": "example@mail.com",
  "password": "sample_password"
}   
```    

#### Response
```json
{
    "name": "Sample Name",
    "email": "example@mail.com"
} 
```

### Sign Up 

#### Request

`POST /auth/signup`
    
```json
{
  "name": "Sample Name",
  "email": "example@mail.com",
  "password": "sample_password"
}   
```    

#### Response
```json
{
    "name": "Sample Name",
    "email": "example@mail.com"
} 
```
---
### Error

#### Response
```json
{
    "message": "Not Accessible",
    "error": "unauthorized",
    "status": 401
}
```

