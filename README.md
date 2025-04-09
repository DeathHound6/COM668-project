# A.I.M.S (Automated Incident Management System)

## What is A.I.M.S?
A.I.M.S (an acronym for Automated Incident Management System) is as it is named.<br/>
This project is designed to automatically detect and raise tickets to track critical client-facing bugs.

# How does it work?
The background processor periodically scans configured app monitoring systems (e.g. Sentry, DynaTrace) and pulls any new exceptions into this system, and searches the stacktrace to determine the root cause.

The system hashes the error stacktrace with the SHA-1 algorithm to check if the error has already been handled.

## Future Improvements
The system at the moment only captures the error message, however through the use of a custom trained AI model, the system could make use of remote Version Control Systems (e.g. Github, Bitbucket) to have more insight into the whole application and determine the cause more accurately and reliably