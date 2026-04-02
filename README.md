# **panzerbot**

## **A. Prerequisites**

### **1. Install Required System Packages**

```bash
sudo apt install libasound2-dev libopus-dev libvpx-dev libv4l-dev
```

### **2. Install Go**

```bash
# Install Go
wget https://go.dev/dl/go1.26.1.linux-arm64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.26.1.linux-arm64.tar.gz
rm go1.26.1.linux-arm64.tar.gz

# Update PATH and Check Version
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
go version
```

### **3. Install NVM, Node.js, and npm**

```bash
```
