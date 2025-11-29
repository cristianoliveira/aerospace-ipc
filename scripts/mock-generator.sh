#!/usr/bin/env bash

set -e

# This script generates mock files for all interfaces in the pkg directory.
# It uses the `mockgen` tool to generate the mocks.
#
# Usage:
#  ./scripts/mock-generator.sh [output_dir]
#
# If output_dir is not provided, it defaults to $PWD/mocks

OUTPUT_DIR=${1:-"$PWD/mocks"}
PKG_DIR="$PWD/pkg"

if [ ! -d "$PKG_DIR" ]; then
    echo "Error: pkg directory not found at $PKG_DIR"
    exit 1
fi

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

echo "Generating mocks for all interfaces in $PKG_DIR"
echo "Output directory: $OUTPUT_DIR"
echo ""

# Find all .go files in pkg directory (excluding test files)
while IFS= read -r -d '' file; do
    # Skip test files
    if [[ "$file" == *_test.go ]]; then
        continue
    fi
    
    # Check if file contains interfaces
    if ! grep -q "^type [A-Z].* interface" "$file"; then
        continue
    fi
    
    # Get relative path from pkg directory
    rel_path="${file#$PKG_DIR/}"
    
    # Get directory and filename
    file_dir=$(dirname "$rel_path")
    file_name=$(basename "$file" .go)
    
    # Create output directory structure
    if [ "$file_dir" != "." ]; then
        output_subdir="$OUTPUT_DIR/$file_dir"
    else
        output_subdir="$OUTPUT_DIR"
    fi
    mkdir -p "$output_subdir"
    
    # Get package name from the directory
    if [ "$file_dir" != "." ]; then
        # For nested packages, use the last directory name
        package_name=$(basename "$file_dir")
    else
        # For root pkg files, use the package name from the file
        package_name=$(grep "^package " "$file" | head -1 | awk '{print $2}')
    fi
    
    # Generate mock file name
    mock_file_name="${file_name}_mock.go"
    mock_file_path="$output_subdir/$mock_file_name"
    
    # Generate the mock file using mockgen
    echo "Generating mock for $rel_path"
    echo "  -> $mock_file_path (package: ${package_name}_mock)"
    
    if ! mockgen -source="$file" -destination="$mock_file_path" -package="${package_name}_mock"; then
        echo "  ERROR: Failed to generate mock for $rel_path"
        exit 1
    fi
    
    echo ""
done < <(find "$PKG_DIR" -type f -name "*.go" -print0)

echo "Mock generation completed successfully!"
