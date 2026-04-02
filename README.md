# **panzerbot**

## **A. Prerequisites**

### **1. Install Required System Packages**

```bash
sudo apt install curl make libasound2-dev libopus-dev libvpx-dev libv4l-dev
```

### **2. Install Go**

```bash
# Install Go
wget https://go.dev/dl/go1.26.1.linux-arm64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.26.1.linux-arm64.tar.gz
rm go1.26.1.linux-arm64.tar.gz

# Update PATH and Check Version
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc
go version

# Install swaggo
go install github.com/swaggo/swag/cmd/swag@latest
```

### **3. Install NVM, Node.js, and npm**

```bash
wget -qO- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.4/install.sh | bash
nvm install 24
nvm use 24
```

## **B. Installation**

### **1. Backend**

* Create `.env` and change `CORS_ALLOWED_ORIGINS` to the frontend base URL

    ```bash
    cp .env.example .env
    nano .env
    ```

* Build the binary

    ```bash
    make build
    ```

* Install to systemd with

    ```bash
    make install
    ```

* Uninstall with

    ```bash
    make uninstall
    ```

### **2. Frontend**

* Setup `.env` and change `NEXT_PUBLIC_BACKEND_ORIGIN` to the backend domain and `NEXT_PUBLIC_BACKEND_PROTO` to `https` for WebRTC.

    ```bash
    cp .env.example .env
    nano .env
    ```

* Build the project

    ```bash
    npm run build
    ```

* Install to systemd with

    ```bash
    make install
    ```

* Uninstall with

    ```bash
    make uninstall
    ```