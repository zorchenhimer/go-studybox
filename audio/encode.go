package audio

// TODO:
// - Stereo with the recorded audio, not just data.
// - Configurable lead-in silence (from start of audio)
// - Configurable segment gap lengths (silence between segments)

import (
	"math"
	"slices"
	"io"
	"bytes"
	"fmt"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"

	"git.zorchenhimer.com/Zorchenhimer/go-studybox/rom"
)

const (
	SampleRate uint32 = 44_100 // TODO: verify sample rate with SBX audio
	Amplitude  int = 16_000
)

var (
	BitRate    int = 4890
)

func EncodeRom(w io.WriteSeeker, sbx *rom.StudyBox) error {
	if sbx == nil {
		return fmt.Errorf("nil rom")
	}

	if sbx.Audio == nil {
		return fmt.Errorf("Missing audio")
	}

	if sbx.Audio.Format != rom.AUDIO_WAV {
		return fmt.Errorf("unsupported audio format: %s", sbx.Audio)
	}

	if len(sbx.Data.Pages) == 0 {
		return fmt.Errorf("no pages")
	}

	if len(sbx.Data.Pages[0].Packets) == 0 {
		return fmt.Errorf("no packets")
	}

	wavreader := bytes.NewReader(sbx.Audio.Data)
	decoder := wav.NewDecoder(wavreader)
	if !decoder.IsValidFile() {
		return fmt.Errorf(".studybox file does not contain a valid wav file")
	}

	decoder.ReadInfo()

	if decoder.SampleRate != SampleRate {
		return fmt.Errorf("SampleRate mismatch. Expected %d; found %d", SampleRate, decoder.SampleRate)
	}

	afmt := &audio.Format{
		NumChannels: 1,
		SampleRate: int(decoder.SampleRate),
	}

	writer := wav.NewEncoder(
		w,
		int(decoder.SampleRate),
		int(decoder.BitDepth),
		int(decoder.NumChans),
		1)
	defer writer.Close()

	var err error
	runningSamples := int64(0)

	prevPageLeadIn := 0

	for _, page := range sbx.Data.Pages {
		if prevPageLeadIn > page.AudioOffsetLeadIn {
			return fmt.Errorf("out of order pages (AudioOffsetLeadIn)")
		}

		prevPageLeadIn = page.AudioOffsetLeadIn

		padLen := int64(page.AudioOffsetLeadIn)-runningSamples
		fmt.Printf("padLen: %d = %d - %d\n", padLen, int64(page.AudioOffsetLeadIn), runningSamples)
		if padLen < 0 {
			padLen = 100_000
		}

		err = generatePadding(writer, afmt, padLen)
		if err != nil {
			return fmt.Errorf("generatePadding() error: %w", err)
		}
		runningSamples += padLen

		sampleCount, err := encodePage(writer, afmt, page)
		if err != nil {
			return fmt.Errorf("encodePage() error: %w", err)
		}
		runningSamples += sampleCount
	}

	writer.Close()

	return nil
}

func generatePadding(writer *wav.Encoder, afmt *audio.Format, length int64) error {
	fmt.Println("generatePadding() length:", length)
	buf := &audio.IntBuffer{
		Format: afmt,
		SourceBitDepth: writer.BitDepth,
		Data: make([]int, length),
	}

	return writer.Write(buf)
}

func encodePage(writer *wav.Encoder, afmt *audio.Format, page *rom.Page) (int64, error) {
	runningFract := float64(0)

	//dataLeadLen := int(float64(page.AudioOffsetData - page.AudioOffsetLeadIn) / float64(83))
	samplesPerFlux := float64(writer.SampleRate) / float64(BitRate) / 2
	fmt.Println("samplesPerFlux:", samplesPerFlux)
	dataLeadLen := (page.AudioOffsetData - page.AudioOffsetLeadIn) / int((samplesPerFlux * 9)) / 2
	fmt.Println("dataLeadLen:", dataLeadLen-9)
	lead := &rawData{
		data: slices.Repeat([]byte{0}, dataLeadLen-9),
		pageOffset: int64(page.AudioOffsetData),
	}

	fmt.Println("lead length:", len(lead.data))

	runningFract = lead.encode(samplesPerFlux, 0)
	data := []*rawData{}
	data = append(data, lead)

	for _, packet := range page.Packets {
		d := &rawData{
			data: packet.RawBytes(),
			pageOffset: int64(page.AudioOffsetData),
		}

		runningFract = d.encode(samplesPerFlux, runningFract)

		data = append(data, d)
	}

	sampleCount := int64(0)
	for _, d := range data {
		buf := &audio.IntBuffer{
			Format: afmt,
			SourceBitDepth: writer.BitDepth,
			Data: d.samples,
		}

		err := writer.Write(buf)
		if err != nil {
			return sampleCount, err
		}

		sampleCount += int64(len(d.samples))
	}

	return sampleCount, nil
}

type rawData struct {
	pageOffset int64
	realOffset int64
	data []byte
	samples []int
}

func (d *rawData) encode(samplesPerFlux, runningFract float64) float64 {
	bits := NewBitData(d.data)
	flux := []byte{}
	prev := 0

	for {
		bit, more := bits.Next()
		if !more {
			break
		}

		peek := bits.Peek()

		if bit == 1 {
			flux = append(flux, 1)
		} else {
			flux = append(flux, 0)
		}

		if bit == peek && bit == 0 {
			// clock flux change
			flux = append(flux, 1)
		} else {
			// no clock flux change
			flux = append(flux, 0)
		}
	}

	first := true
	for _, f := range flux {
		amp := Amplitude
		if f == 1 {
			if prev <= 0 {
				prev = 1
			} else {
				prev = -1
			}
			//amp = int(float64(Amplitude) * 1.5)
			first = true
		}

		//if i == 0 {
		//	amp = Amplitude + 1000
		//} else {
		//	amp = Amplitude
		//}

		spf, fract := math.Modf(samplesPerFlux)
		runningFract += fract
		if runningFract >= 1 {
			runningFract -= 1
			spf++
		}

		//d.samples = append(d.samples, slices.Repeat([]int{prev*Amplitude}, int(spf))...)
		for i := 0; i < int(spf); i++ {
			if first {
				d.samples = append(d.samples, prev*int(float64(amp)*1.5))
			}
			d.samples = append(d.samples, prev*amp)
			first = false
			//amp -= 500
		}
		//d.samples = append(d.samples, slices.Repeat([]int{prev*Amplitude}, int(samplesPerFlux))...)
		first = false
	}

	return runningFract
}
