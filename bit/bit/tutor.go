package bit

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/dollarshaveclub/line"
	"github.com/pkg/errors"

	"github.com/fatih/color"
)

type BitConfig struct {
	ShowMode bool
}

type Tutor struct {
	ShowMode bool
	Reader   *bufio.Reader
}

func NewTutor() (*Tutor, error) {
	config, err := readFromConfig()
	if err != nil {
		return nil, err
	}
	return &Tutor{
		ShowMode: config.ShowMode,
		Reader:   bufio.NewReader(os.Stdin),
	}, nil
}

func saveConfig(config BitConfig) error {
	homeDir := os.Getenv("HOME")
	bitConfigPath := path.Join(homeDir, ".config", "bit")
	file, err := os.Create(bitConfigPath)
	if err != nil {
		return errors.Wrap(err, "Failed to open $HOME/.config/bit")
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(&config)
	return errors.Wrap(err, "Failed to write to $HOME/.config/bit")
}

func readFromConfig() (*BitConfig, error) {
	homeDir := os.Getenv("HOME")
	bitConfigPath := path.Join(homeDir, ".config", "bit")
	_, err := os.Stat(bitConfigPath)

	if os.IsNotExist(err) {
		file, err := os.Create(bitConfigPath)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create $HOME/.config/bit")
		}
		defer file.Close()

		config := BitConfig{}
		err = json.NewEncoder(file).Encode(&config)
		return &config, err
	}

	file, err := os.Open(bitConfigPath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open $HOME/.config/bit")
	}

	var config BitConfig
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to decode $HOME/.config/bit")
	}

	return &config, nil
}

func (t *Tutor) explain(explanations ...string) *Tutor {
	if t.ShowMode {
		output := line.New(os.Stdout, "", "", line.WhiteColor)
		output.Println(strings.Join(explanations, "\n\n"))
	}
	return t
}

func (t *Tutor) finalOutput(output string) {
	if !t.ShowMode {
		fmt.Println(output)
	}
}

func (t *Tutor) hide() *Tutor {
	return &Tutor{ShowMode: false}
}

func (t *Tutor) git(args ...string) string {
	out, _ := t.gitF(args...)
	return out
}

func (t *Tutor) gitF(args ...string) (string, error) {
	if t.ShowMode {
		boldgreen := color.New(color.Bold, color.FgBlue)
		output := line.New(os.Stdout, "", "", line.WhiteColor)
		output.Println()
		output.Print("> ").Format(boldgreen).Print("git ").Cyan().Print(strings.Join(args, " "))
		t.Reader.ReadString('\n')
	}

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	cmd := exec.Command("git", append([]string{"-c", "color.ui=always"}, args...)...)
	cmd.Stdout = outBuf
	cmd.Stderr = errBuf

	err := cmd.Run()
	if err != nil {
		if errBuf.String() != "" {
			Fail(errBuf.String())
		}
		Fail(outBuf.String())
	}

	if t.ShowMode {

		fmt.Println()
		fmt.Println("========== Output ==========")
		output := outBuf.String()
		if strings.TrimSpace(output) == "" {
			output = "Nothing!"
		}
		fmt.Print(output)
		fmt.Println("============================")
		fmt.Println()

	}

	return strings.TrimSpace(outBuf.String()), err
}

func p(text string) string {
	return word_wrap(text, 80)
}

func word_wrap(text string, lineWidth int) string {
	words := strings.Fields(strings.TrimSpace(text))
	if len(words) == 0 {
		return text
	}
	wrapped := words[0]
	spaceLeft := lineWidth - len(wrapped)
	for _, word := range words[1:] {
		if len(word)+1 > spaceLeft {
			wrapped += "\n" + word
			spaceLeft = lineWidth - len(word)
		} else {
			wrapped += " " + word
			spaceLeft -= 1 + len(word)
		}
	}
	return wrapped
}
