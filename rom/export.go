package rom

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"git.zorchenhimer.com/Zorchenhimer/go-studybox/build-script"
)

func (sb *StudyBox) Export(directory string, includeAudio bool) error {
	sbj := StudyBoxJson{
		Version: 1,
		Pages:   []jsonPage{},
		Audio:   directory + "/audio" + sb.Audio.ext(),
	}

	bscript := []build.Token{
		&build.TokenStrValue{ValType: "rom", Value: filepath.Base(directory+".studybox")},
		&build.TokenNumValue{ValType: "version", Value: 1},
		&build.TokenStrValue{
			ValType: "fullaudio",
			Value: "audio"+sb.Audio.ext(),
		},
	}

	// for delay resets and data file names
	var prevTok build.Token

	// A "Page" here does not correspond to the entered "Page" number on the
	// title screen.  These are really segments.  The "Page" that is entered on
	// the title screen is stored in the header of a segment.  Multiple
	// segments can have the same "Page" number.
	for pidx, page := range sb.Data.Pages {
		jp := jsonPage{
			AudioOffsetLeadIn: page.AudioOffsetLeadIn,
			AudioOffsetData:   page.AudioOffsetData,
			Data:              []jsonData{},
		}

		file, err := os.Create(fmt.Sprintf("%s/segment-%02d_packet-0000.txt", directory, pidx))
		if err != nil {
			return err
		}
		fmt.Fprintln(file, page.InfoString())
		file.Close()

		var dataStartId int
		jData := jsonData{}
		rawData := []byte{}

		for i, packet := range page.Packets {
			switch p := packet.(type) {
			case *packetHeader:
				jData.Type = "header"
				jData.Values = []int{int(p.PageNumber)}

				jp.Data = append(jp.Data, jData)
				jData = jsonData{}

				bscript = append(bscript, &build.TokenNumValue{ValType: "page", Value: int(p.PageNumber)})
				bscript = append(bscript, &build.TokenAudioOffsets{
					LeadIn: uint64(page.AudioOffsetLeadIn),
					Data:   uint64(page.AudioOffsetData),
				})
				prevTok = nil

			case *packetDelay:
				jData.Type = "delay"
				jData.Values = []int{p.Length}

				prevTok = &build.TokenDelay{Value: int(p.Length)}
				bscript = append(bscript, prevTok)

			case *packetWorkRamLoad:
				jData.Type = "script"
				jData.Values = []int{int(p.bankId), int(p.loadAddressHigh)}
				dataStartId = i

				prevTok = &build.TokenData{
					Bank: int(p.bankId),
					Addr: int(p.loadAddressHigh),
				}
				bscript = append(bscript, prevTok)

			case *packetPadding:
				jData.Type = "padding"
				jData.Values = []int{p.Length}
				jData.Reset = false

				jp.Data = append(jp.Data, jData)
				jData = jsonData{}

				prevTok = nil
				bscript = append(bscript, &build.TokenNumValue{
					ValType: "padding",
					Value: int(p.Length),
				})

			case *packetMarkDataStart:
				jData.Values = []int{int(p.ArgA), int(p.ArgB)}
				jData.Type = p.dataType()
				dataStartId = i

				prevTok = &build.TokenData{
					Bank: int(p.ArgA),
					Addr: int(p.ArgB),
				}
				bscript = append(bscript, prevTok)

			case *packetMarkDataEnd:
				jData.Reset = p.Reset

				if jData.Values == nil || len(jData.Values) == 0 {
					fmt.Printf("[WARN] No data at page %d, dataStartId: %d\n", pidx, dataStartId)
					jp.Data = append(jp.Data, jData)
					jData = jsonData{}
					continue
				}

				switch jData.Type {
				case "pattern":
					jData.File = fmt.Sprintf("%s/segment-%02d_packet-%04d_chrData.chr", directory, pidx, dataStartId)
					d := prevTok.(*build.TokenData)
					d.ValType = "pattern"

				case "nametable":
					jData.File = fmt.Sprintf("%s/segment-%02d_packet-%04d_ntData.dat", directory, pidx, dataStartId)
					d := prevTok.(*build.TokenData)
					d.ValType = "tiles"

				case "script":
					jData.File = fmt.Sprintf("%s/segment-%02d_packet-%04d_scriptData.dat", directory, pidx, dataStartId)
					d := prevTok.(*build.TokenData)
					d.ValType = "script"

					//script, err := DissassembleScript(scriptData)
					//if err != nil {
					//	fmt.Println(err)
					//} else {
					//	fmt.Printf("Script OK Page %02d @ %04d\n", pidx, dataStartId)
					//	err = script.WriteToFile(fmt.Sprintf("%s/script_page%02d_%04d.txt", directory, pidx, dataStartId))
					//	if err != nil {
					//		return fmt.Errorf("Unable to write data to file: %v", err)
					//	}
					//}

				case "delay":
					jp.Data = append(jp.Data, jData)
					jData = jsonData{}
					continue

				default:
					return fmt.Errorf("[WARN] unknown end data type: %s\n", jData.Type)
				}

				if prevTok != nil {
					if jData.Type == "delay" {
						d := prevTok.(*build.TokenDelay)
						d.Reset = p.Reset
					} else {
						d := prevTok.(*build.TokenData)
						d.File = filepath.Base(jData.File)
					}
				}

				err = os.WriteFile(jData.File, rawData, 0666)
				if err != nil {
					return fmt.Errorf("Unable to write data to file [%q]: %v", jData.File, err)
				}

				jp.Data = append(jp.Data, jData)
				jData = jsonData{}
				rawData = []byte{}

			case *packetBulkData:
				if rawData == nil {
					rawData = []byte{}
				}
				rawData = append(rawData, p.Data...)

			default:
				return fmt.Errorf("Encountered an unknown packet: %s segment: %d", p.Asm(), pidx)
			}
		}

		sbj.Pages = append(sbj.Pages, jp)
	}

	if sb.Audio == nil {
		return fmt.Errorf("Missing audio!")
	}

	if includeAudio {
		err := sb.Audio.WriteToFile(directory + "/audio")
		if err != nil {
			return fmt.Errorf("Error writing audio file: %v", err)
		}
	}

	rawJson, err := json.MarshalIndent(sbj, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(directory+".json", rawJson, 0666)
	if err != nil {
		return err
	}

	bfile, err := os.Create(directory+".sbb")
	if err != nil {
		return err
	}
	defer bfile.Close()

	for _, tok := range bscript {
		if tok.Type() == "page" {
			_, err = fmt.Fprintln(bfile, "")
			if err != nil {
				return fmt.Errorf("error writing bscript file: %w", err)
			}
		}
		_, err = fmt.Fprintln(bfile, tok.Text())
		if err != nil {
			return fmt.Errorf("error writing bscript file: %w", err)
		}
	}

	return nil
}
