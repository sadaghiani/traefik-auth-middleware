# Traefik auth middleware

A middleware for validation JWT and passing claims to the header.

## 1. Install

**Traefik config** in [example](./example/)

### 1.1.1 Start middleware local

    version: "3.x"
    services:
      traefik:
        ...
        command:
          - --experimental.localPlugins.auth-middleware.modulename=github.com/sadaghiani/traefik-auth-middleware
          ...
        volumes:
          - $GOPATH/src/github.com/sadaghiani/traefik-auth-middleware:/plugins-local/src/github.com/sadaghiani/traefik-auth-middleware
          ...

**OR**

### 1.1.2 Start middleware public

    version: "3.x"
    services:
      traefik:
        ...
        command:
          - --experimental.plugins.auth-middleware.modulename=github.com/sadaghiani/traefik-auth-middleware
          - --experimental.plugins.auth-middleware.version=v1.1.0
          ...

Note :

- The "traefik-plugin" topic must be set to github repository.


### 1.2 Active middleware

    version: "3.x"
    services:
      traefik:
        ...
        labels:
        
          - traefik.http.middlewares.my-auth-middleware.plugin.auth-middleware.secretKey=165d042e466fcff0ce6914dae33e5b93 
          - traefik.http.middlewares.my-auth-middleware.plugin.auth-middleware.nameOfAuthorizationHeader=Authorization
          - traefik.http.middlewares.my-auth-middleware.plugin.auth-middleware.nameOfUserIDClaim=userID
          - traefik.http.middlewares.my-auth-middleware.plugin.auth-middleware.nameOfRoleIDClaim=roleID
          - traefik.http.middlewares.my-auth-middleware.plugin.auth-middleware.nameForUserIDHeader=x-user-id
          - traefik.http.middlewares.my-auth-middleware.plugin.auth-middleware.nameForRoleIDHeader=x-role-id
          
          ...


Note :

- The "secretKey" must be set to secret manager.


## 2. Usage

Service config in [example](./example/)

    version: "3.x"
    services:
      example:
        ...
        labels:
          - traefik.http.routers.example.middlewares=my-auth-middleware
          ...


## References

- https://plugins.traefik.io/create
- https://github.com/traefik/plugindemo
