# filesystem backup

> A simple filesystem backup for our source code projects that archive specified folders 
and save a snapshot of them every time we make a change. The change could be when we tweak a file and save it, when we add new files and folders, or even when we delete a file.

A listing of some high-level acceptance criteria for our solution and the approach we want to take:

* The solution should create a snapshot of our files at regular intervals as we make changes to our source code projects

* We want to control the interval at which the directories are checked for changes

* Code projects are primarily text-based, so zipping the directories to generate archives will save a lot of space

* We will build this project quickly, while keeping a close watch over where we might want to make improvements later

* Any implementation decisions we make should be easily modified if we decide to change our implementation in the future

* We will build two command-line tools: the backend daemon that does the work and a user interaction utility that will let us list, add, and remove paths from the backup service
