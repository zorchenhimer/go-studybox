package main

import (
	"fmt"
	"os"
	"image"
	"image/draw"
	"image/color"
	"image/png"
	//"strings"
	"encoding/binary"

	"github.com/alexflint/go-arg"

	nesimg "github.com/zorchenhimer/go-retroimg"
	//"github.com/zorchenhimer/go-retroimg/palette"
)

type Arguments struct {
	Nametable string `arg:"--nt,required"`
	Chr string `arg:"--chr,required"`
	Output string `arg:"--output,required"`
	IsSprite bool `arg:"--sprites"`
	SpriteSheet bool `arg:"--sprite-sheet"`
	NoBg bool `arg:"--no-bg"` // don't draw hot-pink background
}

var (
	tileMissing *nesimg.Tile
)

func main() {
	args := &Arguments{}
	arg.MustParse(args)

	var err error
	tileMissing, err = nesimg.NewTileFromPlanes([][]byte{
		{0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33},
		{0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F} })
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err = run(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args *Arguments) error {
	fmt.Println("--")
	fmt.Println("Nametable:", args.Nametable)
	fmt.Println("CHR:", args.Chr)
	fmt.Println("Output:", args.Output)
	fmt.Println("IsSprite:", args.IsSprite)
	fmt.Println("SpriteSheet:", args.SpriteSheet)
	fmt.Println("--")

	chrFile, err := os.Open(args.Chr)
	if err != nil {
		return err
	}
	defer chrFile.Close()

	raw := nesimg.NewRawChr(chrFile)
	tiles, err := raw.ReadAllTiles(nesimg.BD_2bpp)
	if err != nil {
		return err
	}

	//fmt.Printf("first tile: %#v\n", tiles[0])

	//tileUniform_00 := image.NewUniformPaletted(nesimg.DefaultPal_2bpp, 0)
	//tileUniform_F0 := image.NewUniformPaletted(nesimg.DefaultPal_2bpp, 1)
	//tileUniform_0F := image.NewUniformPaletted(nesimg.DefaultPal_2bpp, 2)
	//tileUniform_FF := image.NewUniformPaletted(nesimg.DefaultPal_2bpp, 3)

	tileUniform_00, err := nesimg.NewTileFromPlanes([][]byte{
		{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00} })
	if err != nil {
		return err
	}

	tileUniform_F0, err := nesimg.NewTileFromPlanes([][]byte{
		{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00} })
	if err != nil {
		return err
	}

	tileUniform_0F, err := nesimg.NewTileFromPlanes([][]byte{
		{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF} })
	if err != nil {
		return err
	}

	tileUniform_FF, err := nesimg.NewTileFromPlanes([][]byte{
		{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF} })
	if err != nil {
		return err
	}

	tiles = append([]*nesimg.Tile{
		tileUniform_00,
		tileUniform_F0,
		tileUniform_0F,
		tileUniform_FF,
	}, tiles...)

	ntData, err := os.ReadFile(args.Nametable)
	if err != nil {
		return err
	}

	layers, err := ReadData(ntData, tiles, args.IsSprite)
	if err != nil {
		return fmt.Errorf("ReadData() error: %w", err)
	}

	uni := image.NewUniform(color.RGBA{0xFF, 0x00, 0xFF, 0xFF})
	var screen *image.RGBA
	if args.SpriteSheet {
		maxWidth := 0
		maxHeight := 0
		for _, l := range layers {
			maxWidth += l.Bounds().Dx()
			if l.Bounds().Dy() > maxHeight {
				maxHeight = l.Bounds().Dy()
			}
		}
		screen = image.NewRGBA(image.Rect(0, 0, maxWidth, maxHeight))
		if !args.NoBg {
			draw.Draw(screen, screen.Bounds(), uni, image.Pt(0, 0), draw.Over)
		}

		loc := image.Pt(0, 0)
		for _, l := range layers {
			draw.Draw(screen, l.Bounds().Add(loc), l, image.Pt(0, 0), draw.Over)
			loc = loc.Add(image.Pt(l.Bounds().Dx(), 0))
		}

	} else {
		screen = image.NewRGBA(image.Rect(0, 0, 32*8, 30*8))
		if !args.NoBg {
			draw.Draw(screen, screen.Bounds(), uni, image.Pt(0, 0), draw.Over)
		}

		for _, l := range layers {
			if args.IsSprite {
				draw.Draw(screen, l.Bounds().Add(l.Location), l, image.Pt(0, 0), draw.Over)
			} else {
				draw.Draw(screen, screen.Bounds(), l, image.Pt(0, 0), draw.Over)
			}
		}
	}

	output, err := os.Create(args.Output)
	if err != nil {
		return err
	}
	defer output.Close()

	err = png.Encode(output, screen)
	if err != nil {
		return err
	}

	//img := image.NewRGBA(image.Rect(0, 0, 32*8, 30*8))
	return nil
}

//type UniformPaletted struct {
//	Palette color.Palette
//	Index int
//}
//
//func NewUniformPaletted(palette color.Palette, idx int) *UniformPaletted {
//	if idx >= len(palette) {
//		panic("Index out of palette range")
//	}
//
//	return &UniformPaletted{
//		Palette: palette,
//		Index: idx,
//	}
//}
//
//func (u *UniformPaletted) At(x, y int) color.Color {
//	if u.Index >= len(u.Palette) {
//		panic("Index out of palette range")
//	}
//
//	return u.Palette[u.Index]
//}
//
//func (u *UniformPaletted) ColorModel() color.Model {
//	return u.Palette
//}
//
//func (u *UniformPaletted) Bounds() image.Rectangle {
//	// Copied from image.Uniform.Bounds()
//	return image.Rectangle{image.Point{-1e9, -1e9}, image.Point{1e9, 1e9}}
//}

func BytesToInt(raw []byte) int {
	if len(raw) > 2 {
		panic("only 8 and 16 bit numbers for now")
	}

	if len(raw) == 1 {
		return int(raw[0])
	}

	return int(raw[1])<<8 | int(raw[0])
}

type DataHeader struct {
	PaletteOffset uint16
	ArgB uint16
	ImageCount uint8
}

func (h DataHeader) String() string {
	return fmt.Sprintf("{DataHeader PaletteOffset:$%04X ArgB:$%04X ImageCount:%d}",
		h.PaletteOffset,
		h.ArgB,
		h.ImageCount,
	)
}

type ImageHeader struct {
	Width uint8
	Height uint8
	AttrLength uint16

	// in pixels
	X uint8
	Y uint8
}

func (h ImageHeader) String() string {
	return fmt.Sprintf("{ImageHeader Width:%d[%02X] Height:%d[%02X] AttrLength:$%04X XCoord:%d[%02X] YCoord:%d[%02X]}",
		h.Width, h.Width,
		h.Height, h.Height,
		h.AttrLength,
		h.X, h.X,
		h.Y, h.Y,
	)
}

func ReadData(raw []byte, tiles []*nesimg.Tile, isSprites bool) ([]*Layer, error) {
	//raw, err := io.ReadAll(r)
	//if err != nil {
	//	return nil, err
	//}

	if isSprites {
		tiles = tiles[4:]
	}

	dataHeader := &DataHeader{}
	_, err := binary.Decode(raw, binary.LittleEndian, dataHeader)
	if err != nil {
		return nil, err
	}
	dataHeader.PaletteOffset += 4

	fmt.Println(dataHeader)

	imgHeaders := []*ImageHeader{}
	for i := 0; i < int(dataHeader.ImageCount); i++ {
		head := &ImageHeader{}
		_, err = binary.Decode(raw[4+(i*6)+1:], binary.LittleEndian, head)
		if err != nil {
			return nil, err
		}
		fmt.Println(head)
		imgHeaders = append(imgHeaders, head)
	}

	palettes := []color.Palette{}
	colorIds := [][]string{}
	for i := 0; i < 4; i++ {
		p := color.Palette{}
		idlist := []string{}
		for j := 0; j < 4; j++ {
			v := raw[int(dataHeader.PaletteOffset)+(i*4)+j]
			if v == 0x3D {
				v = 0x0F
			}
			c, ok := NesColors[v]
			if !ok {
				return nil, fmt.Errorf("Color value $%02X invalid", v)
			}
			p = append(p, c)
			idlist = append(idlist, fmt.Sprintf("%02X", v))
		}
		palettes = append(palettes, p)
		colorIds = append(colorIds, idlist)
	}

	fmt.Println("Palettes:")
	for i := 0; i < len(palettes); i++ {
		fmt.Printf("%s: %v\n", colorIds[i], palettes[i])
	}

	nesTileSize := image.Rect(0, 0, 8, 8)

	offset := 5 + len(imgHeaders) * 6
	layers := []*Layer{}
	for _, head := range imgHeaders {
		rect := image.Rect(0, 0, 32, 30)
		if isSprites {
			rect = image.Rect(0, 0, int(head.Width), int(head.Height))
		}

		l := NewLayer(
			rect,
			nesTileSize,
			palettes,
		)
		l.Transparency = isSprites
		l.Location = image.Pt(int(head.X), int(head.Y))
		l.IsSprite = isSprites
		//tileIds := []int{}

		if isSprites {
			for idx, b := range raw[offset:offset+(int(head.Width)*int(head.Height))] {
				id := int(uint(b)) // is this required to ignore negatives?
				if id < len(tiles) {
					l.Tiles[idx] = tiles[id]
				}
			}
		} else {
			for idx, b := range raw[offset:offset+(int(head.Width)*int(head.Height))] {
				row := idx / int(head.Width)
				col := idx % int(head.Width)
				arrIdx := ((row+int(head.Y/8)) * 32) + (col+int(head.X/8))
				//fmt.Printf("%d ", arrIdx)
				id := int(uint(b)) // is this required to ignore negatives?
				if id < len(tiles) {
					l.Tiles[arrIdx] = tiles[id]
					//tileIds = append(tileIds, id)
				}
			}
		}

		//for idx, t := range tileIds {
		//	fmt.Printf("[%d] %02X\n", idx, t)
		//}

		offset += int(head.Width)*int(head.Height)

		if isSprites {
			l.Attributes = raw[offset:offset+int(head.AttrLength)]
		} else {
			// FIXME: This will break on background chunks that aren't a full screen.
			//        Need to verify the anchor point in the firmware for this case (for
			//        when the tile anchor point isn't aligned to an attribute byte).
			//        Partial screen attributes will also be a different size than 8*8.
			l.SetAttributes(raw[offset:offset+int(head.AttrLength)])
		}
		offset += int(head.AttrLength)

		layers = append(layers, l)
	}

	if offset != int(dataHeader.PaletteOffset) {
		fmt.Printf("[[ offset != dataHeader.PaletteOffset: %04X != %04X ]]\n", offset, dataHeader.PaletteOffset)
	}

	// Looks like the data after the palettes is just padding?  If there's a non-zero
	// byte in there, print a warning.
	palend := int(dataHeader.PaletteOffset) + (8*4)
	if palend < len(raw) {
		nonzero := false
		for i := palend; i < len(raw); i++ {
			if raw[i] != 0x00 {
				nonzero = true
				break
			}
		}
		if nonzero {
			fmt.Printf("[[ %d bytes of extra data ]]\n", len(raw) - palend)
		}
	}

	return layers, nil
}

type Layer struct {
	Tiles []*nesimg.Tile
	Attributes []byte
	Palettes []color.Palette
	TileSize image.Rectangle

	Location image.Point

	Width int
	Height int
	Transparency bool // true for sprites
	IsSprite bool

	Solid bool
}

func NewLayer(layerSize image.Rectangle, tileSize image.Rectangle, palettes []color.Palette) *Layer {
	return &Layer{
		Tiles: make([]*nesimg.Tile, layerSize.Dx()*layerSize.Dy()),
		Attributes: make([]byte, layerSize.Dx()*layerSize.Dy()),
		TileSize: tileSize,
		Location: image.Pt(0, 0),
		Width: layerSize.Dx(),
		Height: layerSize.Dy(),
		Palettes: palettes,
	}
}

func (l *Layer) At(x, y int) color.Color {
	if l.Solid {
		return color.RGBA{0x00, 0x00, 0x00, 0xFF}
	}

	width, height := l.TileSize.Dx(), l.TileSize.Dy()

	row := y / height
	col := x / width
	tx  := x % width
	ty  := y % height

	tileIdx := (row*l.Width)+col

	if l.Tiles[tileIdx] == nil {
		return color.RGBA{0x00, 0x00, 0x00, 0x00}
		//return color.RGBA{0xFF, 0x00, 0xFF, 0xFF}
	}

	colorIdx := l.Tiles[tileIdx].ColorIndexAt(tx, ty)
	if l.Transparency && colorIdx == 0 {
		return color.RGBA{0x00, 0x00, 0x00, 0x00}
	}

	palIdx := l.Attributes[tileIdx]
	return l.Palettes[palIdx][colorIdx]
}

func (l *Layer) Bounds() image.Rectangle {
	return image.Rect(0, 0, l.TileSize.Max.X*l.Width, l.TileSize.Max.Y*l.Height)
}

func (l *Layer) ColorModel() color.Model {
	return color.RGBAModel
}

func (sc *Layer) SetAttributes(data []byte) error {
	if len(data) != 64 {
		return fmt.Errorf("Attribute data must be 64 bytes")
	}

	sc.Attributes = make([]byte, 32*30)
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			src := row*8+col
			start := (row*32)*4 + (col*4)

			raw := data[src]

			br := (raw >> 6) & 0x03
			bl := (raw >> 4) & 0x03
			tr := (raw >> 2) & 0x03
			tl := raw & 0x03

			//if row == 0 && col == 0 {
			//	fmt.Printf("br:%02X bl:%02X tr:%02X tl:%02X\n", br, bl, tr, tl)
			//}

			sc.Attributes[start+0+(32*0)] = tl
			sc.Attributes[start+1+(32*0)] = tl
			sc.Attributes[start+0+(32*1)] = tl
			sc.Attributes[start+1+(32*1)] = tl

			sc.Attributes[start+2+(32*0)] = tr
			sc.Attributes[start+3+(32*0)] = tr
			sc.Attributes[start+2+(32*1)] = tr
			sc.Attributes[start+3+(32*1)] = tr

			if row < 7 {
				sc.Attributes[start+0+(32*2)] = bl
				sc.Attributes[start+1+(32*2)] = bl
				sc.Attributes[start+0+(32*3)] = bl
				sc.Attributes[start+1+(32*3)] = bl

				sc.Attributes[start+2+(32*2)] = br
				sc.Attributes[start+3+(32*2)] = br
				sc.Attributes[start+2+(32*3)] = br
				sc.Attributes[start+3+(32*3)] = br
			}
		}
	}

	return nil
}

var NesColors map[byte]color.Color = map[byte]color.Color{
	0x00: color.RGBA{0x66, 0x66, 0x66, 0xFF},
	0x10: color.RGBA{0xAD, 0xAD, 0xAD, 0xFF},
	0x20: color.RGBA{0xFF, 0xFF, 0xEF, 0xFF},
	0x30: color.RGBA{0xFF, 0xFF, 0xEF, 0xFF},

	0x01: color.RGBA{0x00, 0x2A, 0x88, 0xFF},
	0x11: color.RGBA{0x15, 0x5F, 0xD9, 0xFF},
	0x21: color.RGBA{0x64, 0xB0, 0xFF, 0xFF},
	0x31: color.RGBA{0xC0, 0xDF, 0xFF, 0xFF},

	0x02: color.RGBA{0x14, 0x12, 0xA7, 0xFF},
	0x12: color.RGBA{0x42, 0x40, 0xFF, 0xFF},
	0x22: color.RGBA{0x92, 0x90, 0xFF, 0xFF},
	0x32: color.RGBA{0xD3, 0xD2, 0xFF, 0xFF},

	0x03: color.RGBA{0x3B, 0x00, 0xA4, 0xFF},
	0x13: color.RGBA{0x75, 0x27, 0xFE, 0xFF},
	0x23: color.RGBA{0xC6, 0x76, 0xFF, 0xFF},
	0x33: color.RGBA{0xE8, 0xC8, 0xFF, 0xFF},

	0x04: color.RGBA{0x5C, 0x00, 0x7E, 0xFF},
	0x14: color.RGBA{0xA0, 0x1A, 0xCC, 0xFF},
	0x24: color.RGBA{0xF3, 0x6A, 0xFF, 0xFF},
	0x34: color.RGBA{0xFB, 0xC2, 0xFF, 0xFF},

	0x05: color.RGBA{0x6E, 0x00, 0x40, 0xFF},
	0x15: color.RGBA{0xB7, 0x1E, 0x7B, 0xFF},
	0x25: color.RGBA{0xFE, 0x6E, 0xCC, 0xFF},
	0x35: color.RGBA{0xFE, 0xC4, 0xEA, 0xFF},

	0x06: color.RGBA{0x6C, 0x06, 0x00, 0xFF},
	0x16: color.RGBA{0xB5, 0x31, 0x20, 0xFF},
	0x26: color.RGBA{0xFE, 0x81, 0x70, 0xFF},
	0x36: color.RGBA{0xFE, 0xCC, 0xC5, 0xFF},

	0x07: color.RGBA{0x56, 0x1D, 0x00, 0xFF},
	0x17: color.RGBA{0x99, 0x4E, 0x00, 0xFF},
	0x27: color.RGBA{0xEA, 0x9E, 0x22, 0xFF},
	0x37: color.RGBA{0xF7, 0xD8, 0xA5, 0xFF},

	0x08: color.RGBA{0x33, 0x35, 0x00, 0xFF},
	0x18: color.RGBA{0x6B, 0x6D, 0x00, 0xFF},
	0x28: color.RGBA{0xBC, 0xBE, 0x00, 0xFF},
	0x38: color.RGBA{0xE4, 0xE5, 0x94, 0xFF},

	0x09: color.RGBA{0x0B, 0x48, 0x00, 0xFF},
	0x19: color.RGBA{0x38, 0x87, 0x00, 0xFF},
	0x29: color.RGBA{0x88, 0xD8, 0x00, 0xFF},
	0x39: color.RGBA{0xCF, 0xEF, 0x96, 0xFF},

	0x0A: color.RGBA{0x00, 0x52, 0x00, 0xFF},
	0x1A: color.RGBA{0x0C, 0x93, 0x00, 0xFF},
	0x2A: color.RGBA{0x5C, 0xE4, 0x30, 0xFF},
	0x3A: color.RGBA{0xBD, 0xF4, 0xAB, 0xFF},

	0x0B: color.RGBA{0x00, 0x4F, 0x08, 0xFF},
	0x1B: color.RGBA{0x00, 0x8F, 0x32, 0xFF},
	0x2B: color.RGBA{0x45, 0xE0, 0x82, 0xFF},
	0x3B: color.RGBA{0xB3, 0xF3, 0xCC, 0xFF},

	0x0C: color.RGBA{0x00, 0x40, 0x4D, 0xFF},
	0x1C: color.RGBA{0x00, 0x7C, 0x8D, 0xFF},
	0x2C: color.RGBA{0x48, 0xCD, 0xDE, 0xFF},
	0x3C: color.RGBA{0xB5, 0xEB, 0xF2, 0xFF},

	0x0D: color.RGBA{0x00, 0x00, 0x00, 0xFF},
	0x1D: color.RGBA{0x00, 0x7C, 0x8D, 0xFF},
	0x2D: color.RGBA{0x4F, 0x4F, 0x4F, 0xFF},
	0x3D: color.RGBA{0xB8, 0xB8, 0xB8, 0xFF},

	0x0E: color.RGBA{0x00, 0x00, 0x00, 0xFF},
	0x1E: color.RGBA{0x00, 0x00, 0x00, 0xFF},
	0x2E: color.RGBA{0x00, 0x00, 0x00, 0xFF},
	0x3E: color.RGBA{0x00, 0x00, 0x00, 0xFF},

	0x0F: color.RGBA{0x00, 0x00, 0x00, 0xFF},
	0x1F: color.RGBA{0x00, 0x00, 0x00, 0xFF},
	0x2F: color.RGBA{0x00, 0x00, 0x00, 0xFF},
	0x3F: color.RGBA{0x00, 0x00, 0x00, 0xFF},
}

