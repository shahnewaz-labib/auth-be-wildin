# auth-be-wildin

## Architecture Diagram

```mermaid
flowchart TB
    %% Styles
    classDef client fill:#a8d1f0,stroke:#000
    classDef server fill:#90EE90,stroke:#000
    classDef auth fill:#98FB98,stroke:#000
    classDef database fill:#FFE4B5,stroke:#000
    classDef endpoint fill:#F0F8FF,stroke:#000
    
    %% Client Layer
    Client["Client Browser/API Client\n(with Cookies)"]:::client
    
    %% Server Layer
    subgraph Server["HTTP Server Layer"]
        MainServer["HTTP Server\n(Port 9000)"]:::server
        Register["POST /register"]:::endpoint
        Login["POST /login"]:::endpoint
        Me["GET /me"]:::endpoint
        Logout["POST /logout"]:::endpoint
    end
    
    %% Auth Layer
    subgraph Auth["Authentication Layer"]
        AuthService["Authentication Service"]:::auth
        SessionMgmt["Session Management"]:::auth
        PasswordHash["Password Handling"]:::auth
    end
    
    %% Database Layer
    subgraph DB["Database Layer"]
        UserDB[("User Storage")]:::database
        SessionDB[("Session Storage")]:::database
    end
    
    %% Connections
    Client <--> MainServer
    MainServer --> Register
    MainServer --> Login
    MainServer --> Me
    MainServer --> Logout
    
    Register --> AuthService
    Login --> AuthService
    Me --> AuthService
    Logout --> AuthService
    
    AuthService --> SessionMgmt
    AuthService --> PasswordHash
    
    SessionMgmt --> SessionDB
    AuthService --> UserDB
    
    %% Click Events
    click MainServer "https://github.com/shahnewaz-labib/auth-be-wildin/blob/master/main.go"
    click AuthService "https://github.com/shahnewaz-labib/auth-be-wildin/blob/master/auth/auth.go"
    click UserDB "https://github.com/shahnewaz-labib/auth-be-wildin/blob/master/db/db.go"
    click SessionDB "https://github.com/shahnewaz-labib/auth-be-wildin/blob/master/db/db.go"
    
    %% Legend
    subgraph Legend
        L1["Client Interface"]:::client
        L2["Server Components"]:::server
        L3["Auth Components"]:::auth
        L4["Database Storage"]:::database
        L5["API Endpoints"]:::endpoint
    end
```

## Run the server

```go
go run main.go
```

## Test the server

```curl
curl -X POST http://localhost:9000/register \
 -H "Content-Type: application/json" \
 -d '{"username": "username", "password": "password"}'

curl -X POST http://localhost:9000/login \
    -H "Content-Type: application/json" \
    -d '{"username": "username", "password": "password"}' \
    -c cookies.txt

curl -X GET http://localhost:9000/me \
    -b cookies.txt

curl -X POST http://localhost:9000/logout \
    -b cookies.txt
```
