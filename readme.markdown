### Clean Architecture Layers

![Go Backend Clean Architecture Diagram](clean_arch.png?raw=true)

1. **Entities (Core):** Represent the business objects and rules.
2. **Use Cases:** Coordinate the flow of data to and from the entities.
3. **Interface Adapters:** Convert data from the format most convenient for the use cases and entities to the format most convenient for the external agency (web, DB, etc.).
4. **Frameworks and Drivers:** Contain frameworks and tools such as the database, web framework, etc. This layer is where the actual implementation of interfaces happens.

```mermaid
graph TD
    root[Clean Architecture]
    root --> Entities_Core
    root --> Use_Cases_Application_Logic
    root --> Interface_Adapters
    root --> Frameworks_and_Drivers
```

       +-------------------+
       |   Presentation    |       (port: http, grpc)
       |-------------------|
       |   Handlers, DTOs  |       (dto)
       +--------^----------+
                |
       +--------|----------+
       |    Application    |       (service)
       |-------------------|
       |  Business Logic   |
       +--------^----------+
                |
       +--------|----------+
       |      Domain       |       (model)
       |-------------------|
       |   Entities, Repos |
       +--------^----------+
                |
       +--------|----------+
       |   Infrastructure  |       (repository)
       |-------------------|
       | Database, APIs    |
       +-------------------+

![Go Backend Clean Architecture Diagram](clean_arch_diagram.webp?raw=true)

### Explanation:

- **Domain Layer**: Contains the core business logic and entities.

  - `model`
  - `service`
  - `dto`

- **Data Layer**: Handles data persistence and retrieval.

  - `repository`

- **Presentation Layer**: Manages the interfaces through which the application interacts with the outside world.

  - `port/http`
    - `routes.go`
    - `handlers.go`
  - `port/grpc`
    - `server.go`
    - `handlers.go`

- **Server Layer**: Configures and runs the server instances.

  - `server/http/server.go`
  - `server/grpc/server.go`

- **Proto**: Contains the protocol buffer definitions.
- **Package Layer**: Utility packages.

This structured diagram provides a clear visual representation of the Clean Architecture layers and the components within each layer.

### Internal Structure

<details>
<summary>Click to expand internal structure details</summary>

**internal**

- **address**
  - **model**
    - [address.go](https://github.com/ahmed-eid-faried/golang-backend-clean-arch/blob/main/internal/address/model/address.go): Defines the data structures for address entities.
  - **dto**
    - [address.go](https://github.com/ahmed-eid-faried/golang-backend-clean-arch/blob/main/internal/address/dto/address.go): Contains Data Transfer Objects (DTOs) for address operations.
  - **repository**
    - [address.go](https://github.com/ahmed-eid-faried/golang-backend-clean-arch/blob/main/internal/address/repository/address.go): Provides the repository interface and implementations for address persistence.
  - **service**
    - [address.go](https://github.com/ahmed-eid-faried/golang-backend-clean-arch/blob/main/internal/address/service/address.go): Implements the business logic for address-related operations.
  - **port**
    - **http**
      - [routes.go](https://github.com/ahmed-eid-faried/golang-backend-clean-arch/blob/main/internal/address/port/http/routes.go): Defines the HTTP routes for address endpoints.
      - [handlers.go](https://github.com/ahmed-eid-faried/golang-backend-clean-arch/blob/main/internal/address/port/http/handlers.go): Implements the HTTP handlers for address-related requests.
    - **grpc**
      - [server.go](https://github.com/ahmed-eid-faried/golang-backend-clean-arch/blob/main/internal/address/port/grpc/server.go): Handles initialization, registration, configuration, and starting of the gRPC server for address services.
      - [handlers.go](https://github.com/ahmed-eid-faried/golang-backend-clean-arch/blob/main/internal/address/port/grpc/handlers.go): Implements the gRPC service methods and business logic for address services.
  - **server**
    - **http**
      - [server.go](https://github.com/ahmed-eid-faried/golang-backend-clean-arch/blob/main/internal/server/http/server.go): Integrates address HTTP routes into the main HTTP server.
    - **grpc**
      - [server.go](https://github.com/ahmed-eid-faried/golang-backend-clean-arch/blob/main/internal/server/grpc/server.go): Integrates address gRPC services into the main gRPC server.

**cmd**

- **api**
  - [main.go](https://github.com/ahmed-eid-faried/golang-backend-clean-arch/blob/main/cmd/api/main.go): Bootstraps and runs the application, including address models, gRPC server, and HTTP server.

**proto**

- **address**
  - [address.proto](https://github.com/ahmed-eid-faried/golang-backend-clean-arch/blob/main/proto/address/address.proto): Defines the Protocol Buffers (proto) schema for the address service and generates the corresponding gRPC code.

</details>

<h2>Project File Structure</h2>

<p>This section provides an overview of the file structure and the purpose of each file in the project.</p>

<table>
  <thead>
    <tr>
      <th>File Path</th>
      <th>Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><a href="https://github.com/ahmed-eid-faried/golang-backend-clean-arch/blob/main/internal/address/model/address.go">internal/address/model/address.go</a></td>
      <td>Defines the data structures for address entities.</td>
    </tr>
    <tr>
      <td><a href="https://github.com/ahmed-eid-faried/golang-backend-clean-arch/blob/main/internal/address/dto/address.go">internal/address/dto/address.go</a></td>
      <td>Contains Data Transfer Objects (DTOs) for address operations.</td>
    </tr>
    <tr>
      <td><a href="https://github.com/ahmed-eid-faried/golang-backend-clean-arch/blob/main/internal/address/repository/address.go">internal/address/repository/address.go</a></td>
      <td>Provides the repository interface and implementations for address persistence.</td>
    </tr>
    <tr>
      <td><a href="https://github.com/ahmed-eid-faried/golang-backend-clean-arch/blob/main/internal/address/service/address.go">internal/address/service/address.go</a></td>
      <td>Implements the business logic for address-related operations.</td>
    </tr>
    <tr>
      <td><a href="https://github.com/ahmed-eid-faried/golang-backend-clean-arch/blob/main/internal/address/port/http/routes.go">internal/address/port/http/routes.go</a></td>
      <td>Defines the HTTP routes for address endpoints.</td>
    </tr>
    <tr>
      <td><a href="https://github.com/ahmed-eid-faried/golang-backend-clean-arch/blob/main/internal/address/port/http/handlers.go">internal/address/port/http/handlers.go</a></td>
      <td>Implements the HTTP handlers for address-related requests.</td>
    </tr>
    <tr>
      <td><a href="https://github.com/ahmed-eid-faried/golang-backend-clean-arch/blob/main/internal/address/port/grpc/server.go">internal/address/port/grpc/server.go</a></td>
      <td>Handles initialization, registration, configuration, and starting of the gRPC server for address services.</td>
    </tr>
    <tr>
      <td><a href="https://github.com/ahmed-eid-faried/golang-backend-clean-arch/blob/main/internal/address/port/grpc/handlers.go">internal/address/port/grpc/handlers.go</a></td>
      <td>Implements the gRPC service methods and business logic for address services.</td>
    </tr>
    <tr>
      <td><a href="https://github.com/ahmed-eid-faried/golang-backend-clean-arch/blob/main/proto/address/address.proto">proto/address/address.proto</a></td>
      <td>Defines the Protocol Buffers (proto) schema for the address service and generates the corresponding gRPC code.</td>
    </tr>
    <tr>
      <td><a href="https://github.com/ahmed-eid-faried/golang-backend-clean-arch/blob/main/internal/server/http/server.go">internal/server/http/server.go</a></td>
      <td>Integrates address HTTP routes into the main HTTP server.</td>
    </tr>
    <tr>
      <td><a href="https://github.com/ahmed-eid-faried/golang-backend-clean-arch/blob/main/internal/server/grpc/server.go">internal/server/grpc/server.go</a></td>
      <td>Integrates address gRPC services into the main gRPC server.</td>
    </tr>
    <tr>
      <td><a href="https://github.com/ahmed-eid-faried/golang-backend-clean-arch/blob/main/cmd/api/main.go">cmd/api/main.go</a></td>
      <td>Bootstraps and runs the application, including address models, gRPC server, and HTTP server.</td>
    </tr>
  </tbody>
</table>

### Visual Representation

```mermaid
graph TD
    style proto fill:#ff9
    root[golang-backend-clean-arch]
    root --> cmd
    root --> docs
    root --> internal
    root --> pkg
    root --> test
    root --> proto
    root --> templates

    subgraph cmd
        main[main.go]
    end

    subgraph internal
        address
        server
    end

    subgraph address
        model[model/address.go]
        dto[dto/address.go]
        repository[repository/address.go]
        service[service/address.go]
        port_http[port/http]
        port_grpc[port/grpc]
    end

    subgraph port_http
        http_routes[routes.go]
        http_handlers[handlers.go]
    end

    subgraph port_grpc
        grpc_server[server.go]
        grpc_handlers[handlers.go]
    end

    subgraph server
        http_server[http/server.go]
        grpc_server[grpc/server.go]
    end

    subgraph proto
        address_proto[address.proto]
    end
    sequenceDiagram
        pkg-->config[config]
        pkg-->dbs[dbs]
        pkg-->jtoken[jtoken]
        pkg-->middleware[middleware]
        pkg-->paging[paging]
        pkg-->redis[redis]
        pkg-->response[response]
        pkg-->utils[utils]
```
