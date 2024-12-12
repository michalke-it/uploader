# uploader
Simple HTTP server to upload files.

I needed an HTTP server that allows me to upload files from my phone to a directory. Directory has to be mounted to the /app/uploads directory in the container.
Example:
```
docker run -ti -v ./uploads/:/app/uploads -p 8080:8080 michalkeit/uploader
```

Files can then be uploaded either via HTTP POST:
```
curl -F 'file=@Wiederrufsformular_DE.pdf' http://localhost:8080/upload
```

or simply by browsing to http://localhost:8080/
