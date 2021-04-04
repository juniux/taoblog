package client

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/movsb/taoblog/protocols"
	"github.com/spf13/cobra"
)

// Backup backups all blog database.
func (c *Client) Backup(cmd *cobra.Command) {
	compress := true

	resp, err := c.grpcClient.Backup(c.token(), &protocols.BackupRequest{Compress: compress})
	if err != nil {
		panic(err)
	}

	var reader io.ReadCloser

	if compress {
		r, err := zlib.NewReader(bytes.NewReader(resp.Data))
		if err != nil {
			panic(err)
		}
		reader = r
	} else {
		reader = ioutil.NopCloser(bytes.NewReader(resp.Data))
	}

	defer reader.Close()

	var w io.Writer

	bStdout, err := cmd.Flags().GetBool(`stdout`)
	if err != nil {
		panic(err)
	}
	if !bStdout {
		name := time.Now().Format(`taoblog-2006-01-02.db`)
		fp, err := os.Create(name)
		if err != nil {
			panic(err)
		}
		defer fmt.Printf("Filename: %s\n", name)
		defer fp.Close()
		w = fp
	} else {
		w = os.Stdout
	}

	if _, err := io.Copy(w, reader); err != nil {
		panic(err)
	}
}
