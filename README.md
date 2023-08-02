# Lls
Lls is a link shortener with built-in statistics functionality that allows for convenient link shortening. It is written in GO and uses the gin-gonic framework, as well as mangodb as its database, It is open-source and can run on most platforms.

It provides convenient APIs that enable users to quickly access the link shortening functions of the application. Additionally, the application uses request Captcha and request throttling to avoid malicious use and improve the security and stability of the system.

The application also features two management interfaces, which allow users to query statistics and delete previously created links, and has password protection. Furthermore, the application is built using Docker technology, allowing for quick start-up and improving portability and scalability.

Lls is an efficient, secure, and user-friendly application that can be used to shorten links and track link usage.

## Run with docker
This was built using Docker-Compose. In order to run in your machine, just clone the repository and copy your foreground project to . /resources/ui/ and run:
* cp ./resources/statik/mongo-init.js ./resources/statik/app.ini ./
* perl -p -i -e "s/VFSNnSFLvfOwFnBh/***<font color="red">{Set_database_password_here}</font>***/g" ./mongo-init.js ./app.ini
* sudo docker-compose up --build -d (first time, then you go up without the --build, which is much faster)

If you don't have docker installed in your computer, please download and install it at:
https://www.docker.com/

## File structure description:
````
├── /app.ini
|   ├── Server configuration file (Created automatically on first boot).
|
├── /logs/
|   ├── Directory to store server log files, including access logs and error logs.
|
├── /db-data/
|   ├── Directory to store MongoDB database data. This folder is used to persist MongoDB data in Docker containers.
|
├── /resources/
|   ├── /ui/
|   |   ├── Directory to store static files for the frontend, such as HTML, CSS, and JavaScript files.
|   ├── /statik/
|   |   ├── Directory to store resource files that the server needs.
|
└── /statik/
    └── /statik.go
        ├── File that stores the Go code generated by Statik after packaging. This file contains all files under the /resources/ directory.
````
### File and folder description:

#### /app.ini
This file is the server configuration file, which contains various parameters and settings required for the server to run.

#### /logs/
This folder is used to store server log files, including access logs and error logs.

#### /db-data/
This folder is used to persist MongoDB data in Docker containers. In this folder, you can persist MongoDB database data so that it can continue to be used after the container is restarted.

#### /resources/ui/
This folder stores static files for the frontend, such as HTML, CSS, and JavaScript files.

#### /resources/statik/
This folder stores resource files that the server needs, such as images and fonts.

#### /statik/statik.go
This file is a Go code file generated automatically by the Statik tool, which is used to embed all files under the /resources/ directory into Go code so that they can be accessed as static files.


## API Instructions
### Captcha
To perform a create/manage operation you need to create Captcha first, just http GET to `http://localhost:8040/api/captcha`, The API will return the following:
```json5
{
  "code":0,
  "data":{
    "pic":"data:image/png;base64,....."
  },
  "detail":"",
  "fail":false,
  "message":"",
  "success":true,
  "type":""
}
```

Then, a Base64-encoded challenge image and a cookie identifying the Session are returned

### Generate

To shorten a URL, just http POST to `http://localhost:8040/api/generate_link` with the following json payload (example):

```json5
{
  "link":"http://127.0.0.1:8040/", //Original URL
  "captcha":"8", //Captcha answer
  "pwd": "" //Shortened Access Password
}
```

The api will return the following:

```json5
{
  "code":0,
  "data":{
    "hash":"18nfqL", //shortened URL Hash
    "token":"IKmXKMrVtBOvdibt" //Manage Password
  },
  "detail":"",
  "fail":false,
  "message":"",
  "success":true,
  "type":""
}
```

The token is your subsequent credentials for managing the link, and the hash is the shortened URL Hash

### Redirect

just HTTP GET request to `http://localhost:8040/s/:hash`. This action will lead you to the original URL that was shortened. Here's an example:

`http://localhost:8040/s/18nfqL?...Parameters`

The redirection supports the following parameters:

| Parameter | Description |
| --------- | ----------- |
| `pwd`     | This is the password required for accessing the shortened URL. If an incorrect or empty password is provided, you will be redirected to `/#/PasswordRedirect/:hash` |
| `soft`    | This parameter indicates whether a soft redirect is required. If used, you will be redirected to `/#/SoftRedirect/:hash` |
| `detect`  | This parameter activates the Detect mode. Instead of a redirection, data will be returned in JSON format |


If an incorrect or empty password is provided for a password-protected shortened URL, you will be redirected to `/#/PasswordRedirect/:hash`. The front-end will handle the subsequent logic.

If the `soft` parameter is used, you will be redirected to `/#/SoftRedirect/:hash`.

If the `detect` parameter is used, the API will return the following JSON data:

```json5
{
  "code":0,
  "data":{
    "hash":"18nfqL", //shortened URL Hash
    "url":"http://127.0.0.1:8040/", //Original URL
  },
  "detail":"",
  "fail":false,
  "message":"",
  "success":true,
  "type":""
}
```
In this mode, no redirection will occur. Instead, the data will be returned in JSON format.
### Statistics

The application will record the visit and write it to the database, just http POST to `http://localhost:8040/api/stats_link` with the following json payload (example):

```json5
{
  "hash": "18nfqL", //shortened URL Hash
  "token": "IKmXKMrVtBOvdibt", //Manage Password
  "captcha": "25" //Captcha answer
  "page": 1, // Page number of current visit(A positive integer)
  "size": 50 //Size per page(integers from 1-100)
}
```
The api will return the following:

```json5
{
  "code":0,
  "data":{
    "current":1, //current page
    "size":50, //Size set by request
    "pages":1, //Total number of pages
    "total":1, //Total number of results
    "records":[
      {
        "Hash":"18nfqL", //HASH of the query
        "IP":"127.0.0.1", //IP of the visitor
        "Header":{ //Request Header of the visitor
          "Accept-Encoding":[
            "gzip, deflate, br"
          ]
        },
        "Country":"Local Address", //The country indicated by the visitor's IP
        "Area":"Local Address", //The area indicated by the visitor's IP
        "Browser":"Chrome", //The Browser indicated by the visitor's UA
        "BrowserVersion":"109.0.0",//The Browser Version indicated by the visitor's UA
        "OS":"Windows", //The OS indicated by the visitor's UA
        "OSVersion":"10", //The OS Version indicated by the visitor's UA
        "Device":"Other", //The Device indicated by the visitor's UA
        "Created":1675143659 //Access time (seconds timestamp)
      }
    ]
  },
  "detail":"",
  "fail":false,
  "message":"",
  "success":true,
  "type":""
}
```
It will show detailed data about the URL accessed.

### Delete
If the link needs to be removed, just http POST to `http://localhost:8040/api/delete_link` with the following json payload (example):

```json5
{
  "hash": "18nfqL", //shortened URL Hash
  "token": "IKmXKMrVtBOvdibt", //Manage Password
  "captcha": "32" //Captcha answer
}
```
The api will return the following:

```json5
{
  "code":0,
  "data":null,
  "detail":"",
  "fail":false,
  "message":"",
  "success":true,
  "type":""
}
```
The link will be marked for deletion, but note that it can still be queried for statistics using the administrative password.

