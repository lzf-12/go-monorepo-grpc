#!/bin/bash

# =============================================================================
# Environment Setup Script for Ops Monorepo
# =============================================================================
# This script copies .env.example files to .env files for all services and
# the root directory to enable seamless Docker Compose deployment.
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

# Function to copy .env.example to .env if it doesn't exist
copy_env_file() {
    local dir="$1"
    local service_name="$2"
    
    if [ -f "$dir/.env.example" ]; then
        if [ -f "$dir/.env" ]; then
            print_warning "$service_name .env already exists, skipping..."
        else
            cp "$dir/.env.example" "$dir/.env"
            print_success "$service_name .env created from .env.example"
        fi
    else
        print_error "$service_name .env.example not found in $dir"
        return 1
    fi
}

# Main setup function
main() {
    print_status "Setting up environment files for Ops Monorepo..."
    echo
    
    # Get the script directory (should be /scripts)
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    ROOT_DIR="$(dirname "$SCRIPT_DIR")"
    
    # Change to root directory
    cd "$ROOT_DIR"
    
    print_status "Working directory: $ROOT_DIR"
    echo
    
    # Copy root .env file (for database configuration)
    print_status "Setting up root environment file..."
    copy_env_file "$ROOT_DIR" "Root"
    echo
    
    # Copy service .env files
    print_status "Setting up service environment files..."
    
    # Array of services
    services=("svc-user" "svc-notification" "svc-inventory" "svc-order")
    
    for service in "${services[@]}"; do
        service_dir="$ROOT_DIR/services/$service"
        if [ -d "$service_dir" ]; then
            copy_env_file "$service_dir" "$service"
        else
            print_error "Service directory not found: $service_dir"
        fi
    done
    
    echo
    print_status "Environment setup completed!"
    echo
    
    # Display next steps
    echo -e "${BLUE}==============================================================================${NC}"
    echo -e "${BLUE}NEXT STEPS:${NC}"
    echo -e "${BLUE}==============================================================================${NC}"
    echo
    echo "1. Review and customize your environment files:"
    echo "   - Root .env file (database configuration)"
    echo "   - Service .env files (service-specific settings)"
    echo
    echo "2. Update SMTP credentials in:"
    echo "   - Root .env (global SMTP settings)"
    echo "   - services/svc-notification/.env (service-specific SMTP)"
    echo
    echo "3. Change default JWT secret in:"
    echo "   - Root .env (JWT_SECRET)"
    echo "   - services/svc-user/.env (JWT_SECRET)"
    echo
    echo "4. Start the services:"
    echo "   docker-compose up"
    echo
    echo -e "${GREEN}Environment files are ready for Docker Compose!${NC}"
}

# Function to show help
show_help() {
    echo "Environment Setup Script for Ops Monorepo"
    echo
    echo "Usage: $0 [OPTIONS]"
    echo
    echo "Options:"
    echo "  -h, --help     Show this help message"
    echo "  -f, --force    Force overwrite existing .env files"
    echo
    echo "This script copies .env.example files to .env files for:"
    echo "  - Root directory (database configuration)"
    echo "  - All service directories (service-specific configuration)"
    echo
    echo "Environment files are used by Docker Compose for seamless deployment."
}

# Function to force copy (overwrite existing .env files)
force_copy_env_file() {
    local dir="$1"
    local service_name="$2"
    
    if [ -f "$dir/.env.example" ]; then
        cp "$dir/.env.example" "$dir/.env"
        print_success "$service_name .env created/updated from .env.example"
    else
        print_error "$service_name .env.example not found in $dir"
        return 1
    fi
}

# Handle command line arguments
FORCE_OVERWRITE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -f|--force)
            FORCE_OVERWRITE=true
            shift
            ;;
        *)
            print_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# If force flag is set, use force copy function
if [ "$FORCE_OVERWRITE" = true ]; then
    print_warning "Force mode enabled - will overwrite existing .env files"
    
    # Redefine copy function to force overwrite
    copy_env_file() {
        force_copy_env_file "$1" "$2"
    }
fi

# Run main function
main