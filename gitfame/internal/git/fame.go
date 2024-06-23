package git

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/sync/errgroup"

	"github.com/schollz/progressbar/v3"
)

type Stats struct {
	Name    string `json:"name"`
	Lines   int    `json:"lines"`
	Commits int    `json:"commits"`
	Files   int    `json:"files"`
}

type CommitterStat struct {
	Name    string
	Lines   int
	Commits []string
}

type State int

const (
	commitDescription = iota
	groupDescription

	chanSize = 100
	hashSize = 40

	prefixAuthor    = "author "
	prefixCommitter = "committer "
	prefixFilename  = "filename "
)

func Fame(files []string, rep string, revision string, useCommitter bool) ([]Stats, error) {
	g := new(errgroup.Group)

	bar := progressbar.Default(int64(len(files)))

	ch := make(chan []CommitterStat, chanSize)

	for _, file := range files {
		g.Go(func() error {
			cstats, err := fame(file, rep, revision, useCommitter)
			if err != nil {
				close(ch)
				return err
			}

			ch <- cstats
			return nil
		})
	}

	m := make(map[string]Stats)
	commits := make(map[string]map[string]struct{})

	for range files {
		cstats, ok := <-ch
		if !ok {
			break
		}

		for _, cstat := range cstats {
			mapStat, ok := m[cstat.Name]
			if !ok {
				mapStat.Name = cstat.Name
			}

			mapStat.Files += 1
			mapStat.Lines += cstat.Lines

			m[cstat.Name] = mapStat

			_, ok = commits[cstat.Name]
			if !ok {
				commits[cstat.Name] = make(map[string]struct{})
			}

			for _, commit := range cstat.Commits {
				commits[cstat.Name][commit] = struct{}{}
			}
		}

		if err := bar.Add(1); err != nil {
			return nil, fmt.Errorf("bar: %v", err)
		}

	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	for key, value := range commits {
		stat := m[key]
		stat.Commits = len(value)
		m[key] = stat
	}

	return maps.Values(m), nil
}

func fame(file, rep, revision string, useCommitter bool) ([]CommitterStat, error) {
	blameText, err := gitblame(file, rep, revision)
	if err != nil {
		return nil, fmt.Errorf("gitblame: %v", err)
	}

	if len(blameText) == 0 {
		var logText []string

		if logText, err = gitlog(file, rep, revision, useCommitter); err != nil {
			return nil, fmt.Errorf("gitlog: %v", err)
		}

		if len(logText) == 0 {
			return nil, fmt.Errorf("file hasn't info")
		}

		line := logText[0]
		if len(line) < hashSize+1 {
			return nil, fmt.Errorf("invalid git log format")
		}

		hash := line[:hashSize]
		name := line[hashSize+1:]

		return []CommitterStat{{name, 0, []string{hash}}}, nil
	}

	m := make(map[string]CommitterStat)
	commits := make(map[string]string)

	var (
		state State = groupDescription
		skip  int
		hash  string
	)

	for i := 0; i < len(blameText); i++ {
		if state == groupDescription {
			line := strings.Split(blameText[i], " ")

			if len(line) != 4 {
				return nil, fmt.Errorf("invalid git blame format")
			}

			hash = line[0]
			skip, err = strconv.Atoi(line[3])
			if err != nil {
				return nil, fmt.Errorf("atoi: %v", err)
			}

			if name, ok := commits[hash]; !ok {
				state = commitDescription
			} else {
				stat := m[name]
				stat.Name = name
				stat.Lines += skip
				m[name] = stat

				i += skip*2 - 1
			}

		} else {
			if strings.HasPrefix(blameText[i], prefixFilename) {
				state = groupDescription

				stat := m[commits[hash]]
				stat.Lines += skip
				m[commits[hash]] = stat

				i += skip*2 - 1
			}

			var prefix string

			if useCommitter && strings.HasPrefix(blameText[i], prefixCommitter) {
				prefix = prefixCommitter
			} else if !useCommitter && strings.HasPrefix(blameText[i], prefixAuthor) {
				prefix = prefixAuthor
			}

			if prefix != "" {
				name, _ := strings.CutPrefix(blameText[i], prefix)

				if stat, ok := m[name]; ok {
					stat.Commits = append(stat.Commits, hash)
					m[name] = stat
				} else {
					stat.Name = name
					stat.Commits = append(stat.Commits, hash)
					m[name] = stat
				}

				commits[hash] = name
			}
		}

	}

	return maps.Values(m), nil
}

func gitblame(file, rep, revision string) ([]string, error) {
	cmd := exec.Command("git", "blame", "--porcelain", revision, file)
	cmd.Dir = rep

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git blame %s --porcelain: %v", file, err)
	}

	text, err := splitLines(string(out))
	if err != nil {
		return nil, fmt.Errorf("split lines: %v", err)
	}

	return text, nil
}

func gitlog(file, rep, revision string, useCommitter bool) ([]string, error) {
	var format string

	if useCommitter {
		format = `--format=format:%H %cn`
	} else {
		format = `--format=format:%H %an`
	}

	cmd := exec.Command("git", "log", format, revision, "--", file)
	cmd.Dir = rep

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf(`git log --format=format:"%%H %%an" %s -- %s: %v`, revision, file, err)
	}

	text, err := splitLines(string(out))
	if err != nil {
		return nil, fmt.Errorf("split lines: %v", err)
	}

	return text, nil
}
