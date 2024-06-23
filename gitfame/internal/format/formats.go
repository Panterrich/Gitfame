package format

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"gitlab.com/slon/shad-go/gitfame/internal/git"
)

func OutputTabular(stats []git.Stats) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	fmt.Fprintln(w, "Name\tLines\tCommits\tFiles")

	for _, stat := range stats {
		_, err := fmt.Fprintf(w, "%s\t%d\t%d\t%d\n",
			stat.Name, stat.Lines, stat.Commits, stat.Files)

		if err != nil {
			return fmt.Errorf("tabular: %v", err)
		}
	}

	err := w.Flush()
	if err != nil {
		return fmt.Errorf("tabular: %v", err)
	}

	return nil
}

func OutputCSV(stats []git.Stats) error {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	err := w.Write([]string{"Name", "Lines", "Commits", "Files"})
	if err != nil {
		return fmt.Errorf("csv: %v", err)
	}

	for _, stat := range stats {

		err := w.Write([]string{
			stat.Name,
			strconv.Itoa(stat.Lines),
			strconv.Itoa(stat.Commits),
			strconv.Itoa(stat.Files),
		})

		if err != nil {
			return fmt.Errorf("csv: %v", err)
		}
	}

	return nil
}

func OutputJSON(stats []git.Stats) error {
	out, err := json.Marshal(stats)
	if err != nil {
		return fmt.Errorf("json marshal: %v", err)
	}

	fmt.Println(string(out))

	return nil
}

func OutputJSONLines(stats []git.Stats) error {
	for _, stat := range stats {
		out, err := json.Marshal(stat)
		if err != nil {
			return fmt.Errorf("json marshal indent: %v", err)
		}

		fmt.Println(string(out))
	}

	return nil
}
