import hashlib
import os
import json

# This file will store previously computed file hashes
HASH_FILE = 'file_hashes.json'

# Function to calculate SHA256 hash of a file
def get_file_hash(file_path):
    hasher = hashlib.sha256()
    with open(file_path, 'rb') as f:
        # Read the file in chunks to handle large files
        while chunk := f.read(8192):
            hasher.update(chunk)
    return hasher.hexdigest()

# Function to walk through a directory and hash every file
def scan_directory(directory):
    file_hashes = {}
    for root, _, files in os.walk(directory):
        for name in files:
            path = os.path.join(root, name)
            file_hashes[path] = get_file_hash(path)
    return file_hashes

# Load previously stored hashes from disk (if they exist)
def load_previous_hashes():
    if os.path.exists(HASH_FILE):
        with open(HASH_FILE, 'r') as f:
            return json.load(f)
    return {}

# Save current file hashes to disk for next time
def save_hashes(hashes):
    with open(HASH_FILE, 'w') as f:
        json.dump(hashes, f, indent=2)

# Compare the current hashes to the old ones to detect changes
def detect_changes(current, previous):
    added = [f for f in current if f not in previous]
    removed = [f for f in previous if f not in current]
    changed = [f for f in current if f in previous and current[f] != previous[f]]
   
    return added, removed, changed

# Main function to tie everything together
def main(directory):
    print(f"Scanning: {directory}")
    current_hashes = scan_directory(directory)
    previous_hashes = load_previous_hashes()
   
    added, removed, changed = detect_changes(current_hashes, previous_hashes)
   
    # Report results to the user
    if added or removed or changed:
        print("Changes detected:")
        if added:
            print("  Added:")
            for f in added:
                print("   ", f)
        if removed:
            print("  Removed:")
            for f in removed:
                print("   ", f)
        if changed:
            print("  Modified:")
            for f in changed:
                print("   ", f)
    else:
        print("No changes detected.")
   
    # Save current state for future comparison
    save_hashes(current_hashes)

# Run the script with a directory provided as a command-line argument
if __name__ == "__main__":
    import sys
    # Option 1:Get folder path from command-line arguments
    if len(sys.argv) > 1:
        target_directory = sys.argv[1]
    else:
        # Option 2:Hardcore your test folder path here
        target_directory = r"G:Tasks Folder\.git"

if not os.path.exists(target_directory):
        print("Error:Directory does not exist:", target_directory)        
else:
    main(target_directory)         