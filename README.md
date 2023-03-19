## Auth Middleware

A middleware to add the user ID to the request if there is a valid token

## 1. Install

1.. **Traefik dynamic config** in [docker-compose.yml](./example/docker-compose.yml)

1.1 Start middleware

    version: "3.x"
    services:
      traefik:
        ...
        command:
          - --experimental.localPlugins.auth-middleware.modulename=middleware/auth-middleware
          ...

1.2 Active middleware

    version: "3.x"
    services:
      traefik:
        ...
        labels:
        
	      # secretKey must read from docker secret file.
          - traefik.http.middlewares.my-auth-middleware.plugin.auth-middleware.secretKey=qwertyuiop
          
          # The header containing the jwt token
          # Authorization : Bearer eyJhbGciOiJIUzI1NiI...
          - traefik.http.middlewares.my-auth-middleware.plugin.auth-middleware.getAuthorizationHeaderKey=Authorization
          
          # The header where you want the user ID to be placed
          # userID : 12345678
          - traefik.http.middlewares.my-auth-middleware.plugin.auth-middleware.setUserIDHeaderKey=userID
          
          ...

1.3 Add local plugin 

    version: "3.x"
    services:
      traefik:
        ...
        volumes:
        
          # Root path of the middleware project:/plugins-local/src/middleware/auth-middleware
          - $GOPATH:/plugins-local
          
          ...

## 2. Usage

2.. Service config  in  [docker-compose.yml](./example/docker-compose.yml)

    version: "3.x"
    services:
      example:
        ...
        labels:
          - traefik.http.routers.example.middlewares=my-auth-middleware
          ...


