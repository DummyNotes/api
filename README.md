# Notes

A simple notes CRUD api connected to a DynamoDB backend.

This application is NOT production ready!

# Docker

Build the image:
```shell
docker build . -t notes
```

Start the image:
```shell
docker run -it --rm -p 8080:8080 -v /home/user/.aws:/root/.aws notes
```
Add `-e AWS_PROFILE=$AWS_PROFILE` to the docker command if you are using profiles.
