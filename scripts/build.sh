#!/bin/bash
set -e

NAMESPACE="testing"
IMAGE="shilucloud/csi-driver-hostpath-on-steriod:local"
DRIVER_NAME="csi.driver.hostpath.on.steriod"

# colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

pass() { echo -e "${GREEN}✅ $1${NC}"; }
fail() { echo -e "${RED}❌ $1${NC}"; exit
1; }
info() { echo -e "${YELLOW}➜ $1${NC}"; }



# ─────────────────────────────────────────
# STEP 1: BUILD AND LOAD IMAGE
# ─────────────────────────────────────────
info "Building image..."
docker build -t $IMAGE .
kind load docker-image $IMAGE --name kind
pass "Image built and loaded"
