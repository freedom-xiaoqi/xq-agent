---
name: file-organizer
description: A skill to organize files in a directory by grouping them into subfolders based on their extensions.
metadata:
  openclaw:
    requires:
      bins: []
---

# File Organizer Skill

This skill helps you organize files in a specific directory. It scans the directory and moves files into subfolders named after their file extensions (e.g., `.jpg` files go to `images/`, `.pdf` to `documents/`, etc., or simply by extension name like `jpg/`, `pdf/`).

## Capabilities

The agent can perform this task using its built-in tools:
1.  **List files**: Use `file_list` to see what files are in the directory.
2.  **Move files**: Use `shell_run` to move files.
    *   **Windows (PowerShell)**: The agent runs commands in PowerShell.
        *   Create Directory: `New-Item -ItemType Directory -Force "path/to/folder"` or `mkdir "path/to/folder"`
        *   Move File: `Move-Item -Path "source_file" -Destination "target_folder/"` or `move "source" "dest"`
        *   **Important**: Always quote paths to handle spaces (e.g., `"C:\My Documents\file.txt"`).
    *   **Linux/Mac (Bash)**:
        *   Create Directory: `mkdir -p "path/to/folder"`
        *   Move File: `mv "source" "destination"`

## Usage Examples

**User**: "Organize the files in my downloads folder."

**Agent's Plan**:
1.  List files in the target folder using `file_list`.
2.  Analyze the file extensions.
3.  Create necessary subdirectories.
    *   Windows: `shell_run` -> `mkdir "D:\Downloads\pdf"`
4.  Move files to their respective folders.
    *   Windows: `shell_run` -> `move "D:\Downloads\report.pdf" "D:\Downloads\pdf\"`

**User**: "Clean up the desktop by putting all screenshots into a Screenshots folder."

**Agent's Plan**:
1.  List files on the Desktop.
2.  Identify files that look like screenshots (e.g., contain "Screenshot" in name or are .png).
3.  Create "Screenshots" folder if it doesn't exist.
4.  Move those files using `Move-Item` (Windows) or `mv` (Linux/Mac).
