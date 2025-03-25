AuthService Docs
    Using JWT auth with RSA
    Using Access and Refresh Tokens
    Fiber
    Gorm
    Zap
    air (for hot reload)
    swagger for documention
    4 endpoints (http://localhost:8080/swagger/index.html)
        - signup (returns Access and Refresh)
        - login (returns Access and Refresh)
        - check-token (check Access Token)
        - refresh (Check refresh, if it is valid returns new Access and Refresh)


    Creating Keys 
        - openssl rsa -in jwt-private.pem -pubout -out jwt-public.pem
        - openssl genrsa -out jwt-private.pem 2048

    Turning on
        - compose 
        - sudo docker run --rm -it  -v $(pwd):/app <container_name>


        