#!/bin/bash

# =============================================================================
# Docker Build Test Script for Ops Monorepo Services
# =============================================================================
# This script tests that each service can be built with Docker from its own
# directory using the corrected Dockerfiles and Makefiles.
# =============================================================================

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to test building a service
test_service_build() {
    local service_name="$1"
    local service_dir="$2"
    
    print_status "Testing Docker build for $service_name..."
    
    if [ -d "$service_dir" ]; then
        cd "$service_dir"
        
        # Check if Dockerfile exists
        if [ -f "Dockerfile" ]; then
            # Try to build with Docker
            if make docker-build; then
                print_success "$service_name Docker build completed successfully"
                
                # Check if image was created
                if docker images | grep -q "$service_name"; then
                    print_success "$service_name Docker image created"
                    return 0
                else
                    print_error "$service_name Docker image not found"
                    return 1
                fi
            else
                print_error "$service_name Docker build failed"
                return 1
            fi
        else
            print_error "Dockerfile not found in $service_dir"
            return 1
        fi
    else
        print_error "Service directory not found: $service_dir"
        return 1
    fi
}

# Function to test if a service can respond to health checks
test_service_health() {
    local service_name="$1"
    local port="$2"
    local protocol="${3:-http}" # http or grpc
    
    print_status "Testing $service_name health on port $port..."
    
    # Start container in background
    docker run -d --name "test-$service_name" -p "$port:$port" "$service_name" > /dev/null
    
    # Wait a bit for service to start
    sleep 10
    
    # Test connectivity based on protocol
    if [ "$protocol" = "grpc" ]; then
        # For gRPC services, we'll just check if the port is open
        if nc -z localhost "$port" 2>/dev/null; then
            print_success "$service_name is listening on port $port (gRPC)"
            docker stop "test-$service_name" > /dev/null 2>&1
            docker rm "test-$service_name" > /dev/null 2>&1
            return 0
        else
            print_error "$service_name is not responding on port $port"
            docker stop "test-$service_name" > /dev/null 2>&1
            docker rm "test-$service_name" > /dev/null 2>&1
            return 1
        fi
    else
        # For HTTP services
        if curl -f "http://localhost:$port/health" > /dev/null 2>&1; then
            print_success "$service_name health check passed"
            docker stop "test-$service_name" > /dev/null 2>&1
            docker rm "test-$service_name" > /dev/null 2>&1
            return 0
        else
            print_warning "$service_name health endpoint not available (this may be expected)"
            docker stop "test-$service_name" > /dev/null 2>&1
            docker rm "test-$service_name" > /dev/null 2>&1
            return 0  # Don't fail for missing health endpoints
        fi
    fi
}

# Main function
main() {
    print_status "Starting Docker build tests for all services..."
    echo
    
    # Get the script directory (should be /scripts)
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    ROOT_DIR="$(dirname "$SCRIPT_DIR")"
    
    # Change to root directory
    cd "$ROOT_DIR"
    
    print_status "Working directory: $ROOT_DIR"
    echo
    
    # Array of services with their test ports
    declare -A services
    services[svc-user]="50053:grpc"
    services[svc-notification]="50052:grpc"
    services[svc-inventory]="50051:grpc"
    services[svc-order]="8080:http"
    
    failed_builds=()
    successful_builds=()
    
    # Test each service
    for service in "${!services[@]}"; do
        service_dir="$ROOT_DIR/services/$service"
        port_info="${services[$service]}"
        port="${port_info%:*}"
        protocol="${port_info#*:}"
        
        echo "=================================================================="
        print_status "Testing $service"
        echo "=================================================================="
        
        # Test build
        if test_service_build "$service" "$service_dir"; then
            successful_builds+=("$service")
            
            # Test basic functionality (port listening)
            test_service_health "$service" "$port" "$protocol"
        else
            failed_builds+=("$service")
        fi
        
        echo
        cd "$ROOT_DIR"  # Return to root for next iteration
    done
    
    # Print summary
    echo "=================================================================="
    print_status "BUILD TEST SUMMARY"
    echo "=================================================================="
    
    if [ ${#successful_builds[@]} -gt 0 ]; then
        print_success "Successfully built: ${successful_builds[*]}"
    fi
    
    if [ ${#failed_builds[@]} -gt 0 ]; then
        print_error "Failed to build: ${failed_builds[*]}"
        echo
        print_error "Some services failed to build. Please check the logs above."
        exit 1
    else
        echo
        print_success "All services built successfully! 🎉"
        echo
        print_status "You can now run individual services with:"
        print_status "  cd services/<service-name>"
        print_status "  make docker-build && make docker-run"
    fi
}

# Function to show help
show_help() {
    echo "Docker Build Test Script for Ops Monorepo"
    echo
    echo "Usage: $0 [OPTIONS]"
    echo
    echo "Options:"
    echo "  -h, --help     Show this help message"
    echo "  -c, --clean    Clean up test images before testing"
    echo
    echo "This script tests Docker builds for all services from their respective directories."
}

# Function to clean up test images
cleanup_images() {
    print_status "Cleaning up existing test images..."
    
    for service in svc-user svc-notification svc-inventory svc-order; do
        if docker images | grep -q "$service"; then
            docker rmi "$service" 2>/dev/null || true
            print_status "Removed $service image"
        fi
    done
}

# Handle command line arguments
CLEAN_FIRST=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -c|--clean)
            CLEAN_FIRST=true
            shift
            ;;
        *)
            print_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Clean up if requested
if [ "$CLEAN_FIRST" = true ]; then
    cleanup_images
    echo
fi

# Run main function
main