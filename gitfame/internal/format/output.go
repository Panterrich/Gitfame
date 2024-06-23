package format

import (
	"fmt"
	"slices"
	"strings"

	"gitlab.com/slon/shad-go/gitfame/internal/flags"
	"gitlab.com/slon/shad-go/gitfame/internal/git"
)

type cmp = func(a, b git.Stats) int

var (
	cmpByLines cmp = func(a, b git.Stats) int {
		if a.Lines != b.Lines {
			return b.Lines - a.Lines
		}

		if a.Commits != b.Commits {
			return b.Commits - a.Commits
		}

		if a.Files != b.Files {
			return b.Files - a.Files
		}

		return strings.Compare(a.Name, b.Name)
	}

	cmpByCommits cmp = func(a, b git.Stats) int {
		if a.Commits != b.Commits {
			return b.Commits - a.Commits
		}

		if a.Lines != b.Lines {
			return b.Lines - a.Lines
		}

		if a.Files != b.Files {
			return b.Files - a.Files
		}

		return strings.Compare(a.Name, b.Name)
	}

	cmpByFiles cmp = func(a, b git.Stats) int {
		if a.Files != b.Files {
			return b.Files - a.Files
		}

		if a.Lines != b.Lines {
			return b.Lines - a.Lines
		}

		if a.Commits != b.Commits {
			return b.Commits - a.Commits
		}

		return strings.Compare(a.Name, b.Name)
	}
)

func Output(stats []git.Stats, format flags.Format, order flags.Order) error {
	switch order {
	case flags.OrderByLines:
		slices.SortFunc(stats, cmpByLines)
	case flags.OrderByCommits:
		slices.SortFunc(stats, cmpByCommits)
	case flags.OrderByFiles:
		slices.SortFunc(stats, cmpByFiles)
	default:
		return fmt.Errorf("unknown order")
	}

	switch format {
	case flags.FormatTabular:
		return OutputTabular(stats)
	case flags.FormatCSV:
		return OutputCSV(stats)
	case flags.FormatJSON:
		return OutputJSON(stats)
	case flags.FormatJSONLines:
		return OutputJSONLines(stats)
	default:
		return fmt.Errorf("unknown format")
	}
}
