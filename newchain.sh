#!/bin/bash

echo "Let's create a new Ignite blockchain!"

# Prompt for the name of the blockchain
read -p "Enter the name of your blockchain: " blockchain_name

# Prompt for the port number for the Tendermint node
read -p "Enter the port number for the Tendermint node (default is 26657): " tendermint_port
tendermint_port=${tendermint_port:-26657}

# Prompt for the port number for the blockchain API
read -p "Enter the port number for the blockchain API (default is 1317): " blockchain_api_port
blockchain_api_port=${blockchain_api_port:-1317}

# Prompt for the port number for the token faucet
read -p "Enter the port number for the token faucet (default is 4500): " token_faucet_port
token_faucet_port=${token_faucet_port:-4500}

# Scaffold the new blockchain
echo "Scaffolding new blockchain '$blockchain_name'..."
ignite scaffold chain "$blockchain_name"

# Start the blockchain
echo "Starting blockchain..."
cd "$blockchain_name"
ignite chain serve 

# Build the Docker image
echo "Building Docker image..."
docker build -t $USER/$blockchain_name .

# Push the Docker image to Docker Hub
echo "Pushing Docker image to Docker Hub..."
docker login
docker push $USER/$blockchain_name

echo "Done!"
