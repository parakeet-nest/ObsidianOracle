# ObsidianOracle
> 🚧 work in progress

A ChatBot to discuss with your Obsidian Content.

> This project is developed with the Go library [Parakeet](https://github.com/parakeet-nest/parakeet)

## Requirements

- Ollama
- Obsidian
- Docker Compose

## Setup

- Create a `.env` file (see example: `demo.env`)
- Set the `OBSIDIAN_VAULT_PATH` variable
- Set the other variables depending on your needs

## Start

### Create the embeddings and start the application

```bash
docker compose --profile generation --profile application up 
```

### Embeddings generation only

```bash
docker compose --profile generation up 
```

### Start the application only

```bash
docker compose --profile application up 
```

Then open http://localhost:9090/

## Architecture

```mermaid
flowchart TB
    subgraph profiles ["Profiles"]
        generation["generation"]
        application["application"]
    end

    subgraph core ["Core Services"]
        elasticsearch["elasticsearch"]
        elasticsearch_settings["elasticsearch_settings"]
        kibana["kibana"]
    end

    subgraph embeddings ["Embeddings Generation"]
        download-llm-embeddings["download-local-llm-embeddings"]
        create-embeddings["create-embeddings"]
    end

    subgraph app ["Application"]
        download-llm["download-local-llm"]
        backend["backend"]
        frontend["frontend"]
    end

    %% Profile associations
    generation --> elasticsearch
    generation --> elasticsearch_settings
    generation --> kibana
    generation --> download-llm-embeddings
    generation --> create-embeddings

    application --> elasticsearch
    application --> elasticsearch_settings
    application --> kibana
    application --> download-llm-embeddings
    application --> download-llm
    application --> backend
    application --> frontend

    %% Dependencies
    elasticsearch_settings --> elasticsearch
    kibana --> elasticsearch_settings
    create-embeddings --> download-llm-embeddings
    create-embeddings --> kibana
    create-embeddings --> elasticsearch
    backend --> download-llm-embeddings
    backend --> download-llm
    backend --> kibana
    backend --> elasticsearch
    frontend --> backend



    class generation,application profile
    class elasticsearch,elasticsearch_settings,kibana core
    class download-llm-embeddings,create-embeddings embeddings
    class download-llm,backend,frontend app
```