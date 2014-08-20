/*
Copyright (c) 2014 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ls

import (
	"flag"
	"fmt"
	"io"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/govc/cli"
	"github.com/vmware/govmomi/govc/flags"
	"github.com/vmware/govmomi/govc/flags/list"
	"github.com/vmware/govmomi/vim25/mo"
)

type ls struct {
	*flags.ListFlag

	Long bool
}

func init() {
	cli.Register(&ls{})
}

func (cmd *ls) Register(f *flag.FlagSet) {
	f.BoolVar(&cmd.Long, "l", false, "Long listing format")
}

func (cmd *ls) Process() error { return nil }

func (cmd *ls) PathRelativeTo() (govmomi.Reference, error) {
	dc, err := cmd.Datacenter()
	if err != nil {
		return nil, err
	}
	return dc, nil
}

func (cmd *ls) Run(f *flag.FlagSet) error {
	args := f.Args()
	if len(args) == 0 {
		args = []string{"."}
	}

	es, err := cmd.ListSlice(args, true, cmd.PathRelativeTo)
	if err != nil {
		return err
	}

	lr := listResult{
		Elements: es,
		Long:     cmd.Long,
	}

	return cmd.WriteResult(lr)
}

type listResult struct {
	Elements []list.Element `json:"elements"`

	Long bool `json:"-"`
}

func (l listResult) WriteTo(w io.Writer) error {
	var err error

	for _, e := range l.Elements {
		if !l.Long {
			fmt.Fprintf(w, "%s\n", e.Path)
			continue
		}

		switch e.Object.(type) {
		case mo.Folder:
			if _, err = fmt.Fprintf(w, "%s/\n", e.Path); err != nil {
				return err
			}
		case mo.Datacenter:
			if _, err = fmt.Fprintf(w, "%s (Datacenter)\n", e.Path); err != nil {
				return err
			}
		case mo.VirtualMachine:
			if _, err = fmt.Fprintf(w, "%s (VirtualMachine)\n", e.Path); err != nil {
				return err
			}
		case mo.Network:
			if _, err = fmt.Fprintf(w, "%s (Network)\n", e.Path); err != nil {
				return err
			}
		case mo.ComputeResource:
			if _, err = fmt.Fprintf(w, "%s (ComputeResource)\n", e.Path); err != nil {
				return err
			}
		case mo.Datastore:
			if _, err = fmt.Fprintf(w, "%s (Datastore)\n", e.Path); err != nil {
				return err
			}
		}
	}

	return nil
}
