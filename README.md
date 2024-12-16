### Zip api

This is a web application for processing zip archives

the main features:
1. To get info about zip archive
2. To archive the files into zip archive
3. Send files to specified emails

### API
1. Get info about archive: 
```
POST /api/archive/information HTTP/1.1
Content-Type: multipart/form-data; boundary=-{some-random-boundary}

-{some-random-boundary}
Content-Disposition: form-data; name="file"; filename="my_archive.zip"
Content-Type: application/zip

{Binary data of ZIP file}
-{some-random-boundary}--
```

2. Archive files into zip archive

```
POST /api/archive/files HTTP/1.1
Content-Type: multipart/form-data; boundary=-{some-random-boundary}

-{some-random-boundary}
Content-Disposition: form-data; name="files[]"; filename="document.docx"
Content-Type: application/vnd.openxmlformats-officedocument.wordprocessingml.document

{Binary data of file}
-{some-random-boundary}
Content-Disposition: form-data; name="files[]"; filename="avatar.png"
Content-Type: image/png

{Binary data of file}
-{some-random-boundary}--
```

3. Send the file to the email

```
POST /api/archive/files HTTP/1.1
Content-Type: multipart/form-data; boundary=-{some-random-boundary}

-{some-random-boundary}
Content-Disposition: form-data; name="files[]"; filename="document.docx"
Content-Type: application/vnd.openxmlformats-officedocument.wordprocessingml.document

{Binary data of file}
-{some-random-boundary}
Content-Disposition: form-data; name="files[]"; filename="avatar.png"
Content-Type: image/png

{Binary data of file}
-{some-random-boundary}--
```

### How to up?

Clone the repository
```
git clone https://github.com/fallen-fatalist/zip-api
```

Run the main function
```
go run main.go
```

### Configuration

* The application can launch the application in the specific port\for that purpose set "PORT" environment variable to desired port.

* The application can set the maximum size of the request body\ for that purpose set "BODYLIMIT" environment variable to desired body limit in bytes.


### Progress

- [x] Router

- [x] Archive information
    - [x] Archive information controller
    - [x] Archive information service
    - [x] Added docx signature and XML mimetype changing
    
- [ ] Archive formation 
    - [x] Archive files service 
    - [x] Archive files controller
    - [ ] Divide service and controller

- [ ] Move validation from controllers to services

- [ ] Email 
    - [ ] SMTP sender
    - [ ] Mail credentials env variables
    - [ ] Email controller
