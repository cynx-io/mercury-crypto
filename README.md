# Mercury Crypto Service

A gRPC-based cryptocurrency intelligence service written in Go.  
This service fetches data from external APIs like CoinGecko and wraps the results in a structured format through defined protobuf contracts.

## Features

- Coin search by name or symbol
- Token risk analysis:
   - Mint/pause functions
   - Honeypot detection
   - Ownership and fee control flags
   - Holder distribution
   - Social and metadata info
- Standardized response format
- gRPC communication

## Prerequisites

- Go 1.21 or higher
- Protocol Buffers compiler (protoc)
- MySQL/MariaDB

## Installation

1. Clone the repository:
```bash
git clone https://github.com/cynxees/mercury-crypto.git
cd mercury-crypto
```

2. Install dependencies:
```bash
go mod download
```

3. Generate proto files:
```bash
make proto
```

4. Build the application:
```bash
make build
```

## Running the Service

```bash
make run
```

## Development

1. Generate proto files:
```bash
make proto
```

2. Build the application:
```bash
make build
```

3. Run the application:
```bash
make run
```

4. Clean generated files:
```bash
make clean
```

## License

MIT License