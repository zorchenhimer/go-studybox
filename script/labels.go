package script

import (
	"fmt"
	"io"
	"os"
	"encoding/json"
	"strconv"
)

type Label struct {
	Address int
	Name string
	Comment string
	FarLabel bool
}

type JsonLabel struct {
	Address string
	Name string
	Comment string
	FarLabel bool
}

func (l Label) JsonLabel() JsonLabel {
	return JsonLabel{
		Address: fmt.Sprintf("0x%X", l.Address),
		Name: l.Name,
		Comment: l.Comment,
		FarLabel: l.FarLabel,
	}
}

func (l JsonLabel) Label() (*Label, error) {
	addr, err := strconv.ParseInt(l.Address, 0, 32)
	if err != nil {
		return nil, fmt.Errorf("Invalid address: %q", l.Address)
	}

	return &Label{
		Address: int(addr),
		Name: l.Name,
		Comment: l.Comment,
		FarLabel: l.FarLabel,
	}, nil
}

func AutoLabel(address int) *Label {
	return &Label{
		Address: address,
		Name: fmt.Sprintf("L%04X", address),
	}
}

func AutoLabelVar(address int) *Label {
	return &Label{
		Address: address,
		Name: fmt.Sprintf("Var_%04X", address),
	}
}

func AutoLabelFar(address int) *Label {
	return &Label{
		Address: address,
		Name: fmt.Sprintf("F%04X", address),
		FarLabel: true,
	}
}

func NewLabel(address int, name string) *Label {
	return &Label{
		Address: address,
		Name: name,
	}
}

func (s *Script) LabelsFromJsonFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return s.LabelsFromJson(file)
}

func (s *Script) LabelsFromJson(r io.Reader) error {
	lbls := []JsonLabel{}
	dec := json.NewDecoder(r)
	err := dec.Decode(&lbls)
	if err != nil {
		return err
	}

	if s.Labels == nil {
		s.Labels = make(map[int]*Label)
	}

	for _, lbl := range lbls {
		l, err := lbl.Label()
		if err != nil {
			return err
		}

		s.Labels[l.Address] = l
	}

	return nil
}

func (s *Script) WriteLabelsToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return s.WriteLabels(file)
}

func (s *Script) WriteLabels(w io.Writer) error {
	slice := []JsonLabel{}
	for _, lbl := range s.Labels {
		slice = append(slice, lbl.JsonLabel())
	}

	raw, err := json.MarshalIndent(slice, "", "\t")
	if err != nil {
		return err
	}

	_, err = w.Write(raw)
	return err
}
