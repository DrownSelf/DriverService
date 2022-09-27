# Driver Service
***
Service provides functionality to work with driver request.

# External requirements
***
    GO 1.18
    Docker
    Docker-compose
    Godotenv
    gocql
    swaggo/swag

# Configuration 

The service could be configured by providing env vars.
                
| Name            | Meaning                | Example   |
|-----------------|------------------------|-----------|
| CASSANDRA_NODE1 | ip of first node       | 172.0.0.1 |
| CASSANDRA_NODE2 | id of second node      | 172.0.0.2 |
| CASSANDRA_NODE3 | id of third node       | 172.0.0.3 |
| SECRET_KEY      | Key for encode session | your_key  |
| APP_PORT        | Port of application    | 8080      |