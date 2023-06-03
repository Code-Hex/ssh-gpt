# Fictitious SSH Server - Powered by OpenAI GPT-4

This is a program to simulate a fictitious SSH Server using Open AI's [Create chat completion API](https://platform.openai.com/docs/api-reference/chat/create) ([GPT-4](https://openai.com/gpt-4)). This has been created for fun and experimental purposes.

Please do not run this server in a production environment, as it has not been designed with security considerations! Keep it confined to your own machine for playing around :bow:

The current list of features that can be executed is as follows:

- Execute Bash on the ssh server.
- SFTP server.

https://github.com/Code-Hex/ssh-gpt/assets/6500104/5ff51752-364a-43d7-a8c8-853d6c86564d

## How to Use

This program is written in [Go](https://go.dev/) and requires Go 1.20 or higher for building. Please clone the repository and build.

You will need your Open AI API key to run this program. Set the API key in the environment variable `OPENAI_API_KEY` and execute.

```sh
$ git clone https://github.com/Code-Hex/ssh-gpt.git
$ cd ssh-gpt
$ go build .
$ OPENAI_API_KEY="..." ./ssh-gpt # Replace the ... part with the API key and execute!
```

The options available for this program are as follows:

```
-s  hostname of fictitious ssh server (default:"ssh.netflix.net")
-bg background for fictitious ssh server (default:"This is exposed by Netflix to share exclusive materials with their fans, includings a bunch of folders, movies, txt files and more.")
```

### Bash Simulation

After setting up an SSH server using the `ssh-gpt` program, you can connect to a server that simulates Bash execution by using an ssh client as follows:

```sh
$ ssh 127.0.0.1 -p 2222 -o "StrictHostKeyChecking=no"
> ls
exclusive_materials/  movies/  txt_files/  fan_art/  behind_the_scenes/
> 
```

### SFTP Simulation

Currently, it only supports file transfers using the `scp` command. However, it will fail most of the time. The `ssh-gpt` program will output logs to the screen during execution, so you can observe how it processes during transfers.

Example:

```sh
$ scp -v -o "StrictHostKeyChecking=no" -P 2222 -r 127.0.0.1:/home/netflix_fan/trivia_and_easter_eggs/trivia1.txt trivia1.txt
```

You can also connect with a GUI client like [Cyberduck](https://blog.cyberduck.io/), but it will probably fail to get the directory list. That's a challenge for the future. (I'm waiting for someone's nice idea! :smile:)

Currently, the simulation is successful for most SFTP interactions, but it mainly fails in the content generation part, such as `SSH_FXP_HANDLE` and `SSH_FXP_DATA`.

## LICENSE

MIT

## Author

[@codehex](https://twitter.com/codehex)