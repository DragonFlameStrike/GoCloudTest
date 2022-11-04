Simple service to manage config files.

You can Create config from file by post query.
    example : curl http://localhost:8080/config -F "file=@data.json" -vvv
    or use client library
You can Read config(-s) by get query:
    /config without params - return all configs
    /config?service=*** - return config by service value
You can Update config
You can Delete config

All warnings and errors are logging in /tmp/GoCloudTest
