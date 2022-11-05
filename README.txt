Simple service to manage config files.

All files stored with versions:
    How:
        example : file_v1.0.json
    Rules:
        The first digit is responsible for "create" version.
        The second digit is responsible for "update" version.
        The server assigns the versions.

You can Create config from file by post query.
    How:
        example : curl http://localhost:8080/config -F "file=@data.json"
        or use client library
    Rules and hints:
        You shouldn't use your version in file name, this is server work
        You can use a file name that is already in use. File will be saved with different version.
        You can save any json file, but if you want to use it and read - your file should have "service" field.

You can Read config(-s) by get query:
    How:
        /config without params - return all configs
        /config?service=*** - return config by service value

You can Update config
    How:
        Update option works like create option. The difference is that the file is saved with the new "update" version.
    Rules and hints:
        like create option and...
        You can update only those files that are stored on the server

You can Delete config

All warnings and errors are logging in /tmp/GoCloudTest
