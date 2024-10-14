## Project Summary
During the development of the gRPC Todo List application, I encountered several challenges that provided valuable learning experiences, despite not completing every aspect of the project.

## Issues Faced
Docker Container Issues: I faced difficulties getting the Docker container to run properly, particularly with configuration settings it was resolved by changing mysql versions.

Graceful Shutdown: Implementing a graceful shutdown for the gRPC server proved challenging, requiring context management to ensure a clean exit.
Google Package Conflicts: I encountered issues with a Google package needed for generating a server that supports a REST API. I had to clone the actual repository into my Go source files, which later led to conflicts during go mod tidy due to having two versions of the same package with slightly different functions. I resolved this by manually deleting the conflicting modules in my Go path.

Unified Logger: I wished I had enough time to implement a unified logging solution. I would have preferred using a sugared logger from Uber for better log management.

Failing Tests: I have a few failing tests that I haven't pushed to Git yet. My main idea was to create mocks to mimic the application's actual behavior, enhancing test reliability.

CLI and gRPC Separation: In hindsight, I would have separated the CLI from the gRPC implementation further to improve the overall structure and maintainability of the code.

Understanding gRPC
gRPC is a high-performance, open-source RPC (Remote Procedure Call) framework that enables the creation of efficient, distributed applications. It uses HTTP/2 for transport, allowing for features such as multiplexing, streaming, and server push. gRPC supports multiple programming languages and promotes the use of Protocol Buffers (protobuf) for defining service contracts and message types.

One of the key advantages of gRPC is its support for reverse proxies, which allows resources to be accessed via REST APIs. This enables developers to expose gRPC services over standard HTTP, making it easier for web clients and other applications to interact with the service.

## Learning Experience
Overall, this project was a significant learning opportunity. I gained practical experience with gRPC, the Cobra CLI for command-line interactions, Docker for containerization, and integrating a REST API via gRPC-Gateway. Even though I didn’t complete everything, the insights I gained from these challenges will undoubtedly inform my future projects.

## Known defect
CLI Behavior: The CLI exhibited unexpected behavior when run in a Docker container, displaying commands in a loop without waiting for user input .


