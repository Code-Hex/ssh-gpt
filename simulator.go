package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const (
	BashMode = iota
	SFTPMode
)

const explainBashPrompt = `The SSH client we use expects to execute Linux commands on the server. From the client, Linux command line is sent.

Your task is to evaluate these command-lines and simulate the appropriate responses, as if you were the actual user connecting to that server. You will do nothing other than evaluating these commands as Bash.

Note:
- ANSI escape sequence is supported.
- Do not attempt to explain what these commands are doing or what they are for.`

const explainSFTPPrompt = `The SFTP client we use expects file access, file transfer, and file management over any reliable data stream. From the client, a list of 8-bit unsigned integers as packet data is sent.

Your task is to evaluate these packet data and simulate the appropriate responses, as if you were the actual user connecting to that server. You will do nothing other than evaluating these packets.

Note:
- Do not attempt to explain what these packets are doing or what they are for.
- All these interactions are conducted using JSON arrays with comments.
  - The comments will describe information such as the absolute filepath, and its contents.
- SFTP Protocol Version should be 3.
- The following values are defined for packet types:
  - SSH_FXP_INIT: 1
  - SSH_FXP_VERSION: 2
  - SSH_FXP_OPEN: 3
  - SSH_FXP_CLOSE: 4
  - SSH_FXP_READ: 5
  - SSH_FXP_WRITE: 6
  - SSH_FXP_LSTAT: 7
  - SSH_FXP_FSTAT: 8
  - SSH_FXP_SETSTAT: 9
  - SSH_FXP_FSETSTAT: 10
  - SSH_FXP_OPENDIR: 11
  - SSH_FXP_READDIR: 12
  - SSH_FXP_REMOVE: 13
  - SSH_FXP_MKDIR: 14
  - SSH_FXP_RMDIR: 15
  - SSH_FXP_REALPATH: 16
  - SSH_FXP_STAT: 17
  - SSH_FXP_RENAME: 18
  - SSH_FXP_READLINK: 19
  - SSH_FXP_SYMLINK: 20
  - SSH_FXP_STATUS: 101
  - SSH_FXP_HANDLE: 102
  - SSH_FXP_DATA: 103
  - SSH_FXP_NAME: 104
  - SSH_FXP_ATTRS: 105
  - SSH_FXP_EXTENDED: 200
  - SSH_FXP_EXTENDED_REPLY: 201
- Do not to make the byte sequence returned by SSH_FXP_DATA incremental in a continuous sequence. For example, if you return '[0,0,0,26,103,0,0,0,10,0,0,0,0]' in sequence 1, you should not make it incremental like '[0,0,0,26,103,0,0,0,11,0,0,0,0]' in sequence 2, and '[0,0,0,26,103,0,0,0,12,0,0,0,0]' in sequence 3.`

type Simulator struct {
	serverName string
	background string

	history       []openai.ChatCompletionMessage
	mode          int
	systemPrompts []openai.ChatCompletionMessage
	c             *openai.Client
}

func NewSimulator(mode int, serverName, background string) (*Simulator, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("OPENAI_API_KEY env is empty")
	}
	c := openai.NewClient(apiKey)
	return &Simulator{
		serverName: serverName,
		background: background,
		history:    make([]openai.ChatCompletionMessage, 0),
		mode:       mode,
		c:          c,
	}, nil
}

func (s *Simulator) Simulate(ctx context.Context, line string) (string, error) {
	s.history = append(s.history, openai.ChatCompletionMessage{
		Role:    "user",
		Content: line,
	})
	resp, err := s.c.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       openai.GPT4,
		Temperature: 0,
		Messages:    append(s.getSystemPrompts(), s.history...),
	})
	if err != nil {
		return "", err
	}
	result := resp.Choices[0].Message
	s.history = append(s.history, result)
	return result.Content, nil
}

func (s *Simulator) getSystemPrompts() []openai.ChatCompletionMessage {
	if s.systemPrompts != nil {
		return s.systemPrompts
	}

	lines := []string{
		fmt.Sprintf("Simulate a fictitious SSH server %s.\n", s.serverName),
		s.background,
	}
	if s.mode == SFTPMode {
		lines = append(lines, explainSFTPPrompt)
	} else {
		lines = append(lines, explainBashPrompt)
	}

	s.systemPrompts = []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: strings.Join(lines, "\n"),
		},
	}
	if s.mode == SFTPMode {
		s.systemPrompts = append(s.systemPrompts,
			openai.ChatCompletionMessage{
				Role:    "system",
				Content: "[0,0,0,5,1,0,0,0,3]",
				Name:    "user",
			},
			openai.ChatCompletionMessage{
				Role:    "system",
				Content: "// Initialize the session with a protocol version of 3.\n[0,0,0,5,2,0,0,0,3]",
				Name:    "assistant",
			},
		)
	}
	return s.systemPrompts
}
