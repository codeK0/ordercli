# FILE-CHANGES-MONITOR

*COMPANY*: CODTECH IT SOLUTIONS

*NAME*: KANCHAN VILAS JADHAV

*INTERN ID*: CT04DA375

*DOMAIN*: CYBERSECURITY AND EH

*DURATION*: 4 WEEKS

*MENTOR*: NEELA SANTOSH

# DESCRIPTION OF TASK #: File Change Monitoring Script in Python

This Python script is built to continuously monitor and track changes within a specific directory on your system. Its primary function is to detect file creations, modifications, and deletions in real-time or between intervals. This is especially useful for developers, system administrators, and cybersecurity analysts who need to keep an eye on important files and folders.

# Core Python Modules Used:
-os and time: These are used for navigating the file system, retrieving file attributes such as modification times, and implementing time-based functionality like delay or timestamps.

-hashlib: Generates MD5 hash values for file contents. This helps in identifying files that have been modified based on their content, even if their names or sizes remain the same.

-json: Stores and loads snapshots of the file system in a structured format that is easy to compare later.

-argparse: Handles command-line arguments to allow flexible script execution and directory specification.

# Editor used: VS Studio

# How It Works:
1.Initial Directory Scanning: When the script is first run, it scans the target directory and creates a snapshot of all existing files. For each file, it collects metadata including file path, last modified timestamp, file size, and an MD5 hash of the content.

2.Snapshot Saving: This metadata is stored in a Python dictionary and saved as a JSON file (often named snapshot.json) in the script's directory. This serves as a baseline to detect future changes.

3.Change Detection: On subsequent runs, the script re-scans the directory and loads the previous snapshot. It compares each file's current state to its prior state to determine if it:

Was added (exists now but not before),

Was deleted (existed before but not anymore),

Was modified (exists in both scans but with different timestamps or hashes).

4.Logging and Output: Detected changes are printed to the terminal with color-coded logs for better clarity:

Green: New files

Red: Deleted files

Yellow: Modified files

5.Snapshot Update: After logging, the script saves the new state of the folder by updating the snapshot JSON, preparing it for the next run.

# Use Cases:
-Monitoring source code directories for unauthorized changes

-Watching log or config folders on servers

-Auditing file system behavior over time

-Creating lightweight file integrity monitoring tools

# Git Integration (Optional):
The script itself, along with its snapshot files, can be version-controlled using Git. This allows users to track both code changes and detected file changes, and even push them to GitHub for backup or collaboration.

=Typical Workflow:
-Place the script in or point it to the directory you want to monitor.

-Run it using the command:
python script.py [optional_directory_path]

-Make changes to files: add, delete, or modify.

-Run the script again to view a log of those changes.

# OUTPUT
<img width="406" alt="Image" src="https://github.com/user-attachments/assets/9ede3aca-df1b-402b-a51a-4ed0d33febf3" />
