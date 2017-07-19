# webhook-proxy

Docker Golang container to securely proxy Bitbucket webhook requests to Jenkins.

This project is for those of us that run Jenkins on a private network and would like to 
 take advantage of event driven builds via Bitbucket webhooks without having to expose
 Jenkins to the outside world.  This service is designed to be public facing to accept 
 and parse the webhook request.  Only valid whitelisted requests will be forwarded to 
 the Jenkins instance.
 
Configuration options: TODO

Container will listen on port 8080
