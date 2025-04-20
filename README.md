# FILE-CHANGES-MONITOR

*COMPANY*: CODTECH IT SOLUTIONS

*NAME*: KANCHAN VILAS JADHAV

*INTERN ID*: CT04DA375

*DOMAIN*: CYBERSECURITY AND EH

*DURATION*: 4 WEEKS

*MENTOR*: NEELA SANTOSH

## DESCRIPTION OF TASK ##: This Python script is designed to monitor file changes within a specified directory. It detects and logs file creations, modifications, and deletions in real-time.

Core Modules Used:

os and time: For file system navigation and timestamping.

hashlib: To generate file hashes and detect changes.

json: To store and compare file state snapshots.

argparse: For command-line arguments.


How It Works:

1. Directory Scanning:

The script scans all files within a given folder and records their last modified time, size, and content hash (MD5).

This metadata is stored in a dictionary and saved to a .json file as a baseline snapshot.


2. Change Detection:

On subsequent runs or on an interval loop, the script re-scans the directory and compares the current state with the previous snapshot.

It identifies:

New files (present now but not before),

Deleted files (in snapshot but not in the current scan),

Modified files (same file but changed timestamp or content hash).


3. Output:

Changes are printed in the console in color-coded logs:

Green for added,

Red for deleted,

Yellow for modified.


The updated snapshot is then saved again for future comparison.



4. Use Case:

Run manually or set as a background task to monitor important folders like project directories, system logs, etc.



5. Git Integration:

The script itself and its tracked changes can be version-controlled using Git. You can commit changes and push them to GitHub for remote backups or collaboration.


Typical Workflow:

1. Place the script in your desired folder.


2. Run it using:

python script.py


3. Modify, add, or delete files in that folder.


4. Re-run the script to detect and log those changes.


This tool is useful for developers, analysts, or anyone who wants to track file activity in a folder without manually checking everything.

# OUTPUT

