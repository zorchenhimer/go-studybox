# The Image Format

Images are split up into two segments on the tapes: tile data and CHR image
data.  The CHR data is the simplest to explain.  It is exactly the same as your
standard NES CHR.  There are no headers or extra metadata.  (i think the
firmware counts tiles and stores that somewhere for the engine to use)

Images can be either nametable data or sprite data.  Both follow the same data
format.  The only difference in decoding is nametable data expects the first
four CHR tiles to be solid colors (one for each palette color), while the
sprites do not.

Multi-byte values are little endian.

    type DataHeader struct {
        // Offset from the start of the ImageCount byte
        PaletteOffset uint16

        // Unused I think?
        ArgB uint16

        // Number of images in this segment
        ImageCount uint8
    }

    type ImageHeader struct {
        // in tiles
        Width uint8
        Height uint8

        // in bytes
        AttrLength uint16

        // in pixels
        X uint8
        Y uint8
    }

Each file starts with a `DataHeader` and is followed by one or more
`ImageHeader`s.  `DataHeader.ImageCount` containts the number of images (and
headers) are contained in this segment.  After the last `ImageHeader` is tile and
attribute data for every image in sequence.  The length of data for each image
is `Width * Height + AttrLength`.  Immediately following the data for the last
image is the palette data.  This data is 16 bytes (either all of the BG
palettes, or all of the sprite palettes).

## Full Backgrounds

This image is a full background.  The palette offset is `$0407 + 4 = $040B` and
there is a single image in this segment (`$01` at offset `$04`).

The image header starts at offset `$05`: `$20 $1E $40 $00 $00 $00`.  The tile
data is immediately following.  Attribute data starts at `$03CA` and is 64
bytes long (`$40` at offset `$07`).

This image is 32x30 tiles in size starting at an X/Y coordinate of (0, 0).
there are a few padding bytes after the palette data (27 bytes in this case).
This padding doesn't seem to be required.

    00000000  07 04 2b 00 01 20 1e 40  00 00 00 02 02 02 02 02  |..+.. .@........|
    00000010  02 02 02 02 02 02 02 02  02 02 02 02 02 02 02 02  |................|
    ...
    skipping to attribute and palette data
    ...
    000003b0  02 02 02 02 02 02 02 02  02 02 02 02 02 02 02 02  |................|
    000003c0  02 02 02 02 02 02 02 02  02 02 02 00 00 00 00 00  |................|
    000003d0  00 00 00 00 00 00 00 00  40 50 10 04 05 01 00 00  |........@P......|
    000003e0  00 05 01 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
    000003f0  00 00 00 00 00 00 00 00  40 10 00 00 00 00 00 00  |........@.......|
    00000400  04 01 00 00 00 00 00 00  00 00 00 0f 12 20 16 0f  |............. ..|
    00000410  27 20 16 0f 27 20 12 0f  28 20 12 00 00 00 00 00  |' ..' ..( ......|
    00000420  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
    00000430  00 00 00 00 00 00                                 |......|

## Partial Backgrounds & Metasprites

A partial background is laid out exactly like a full background.  The only
difference is that the X/Y coordinates will most likely not be (0, 0) and
`Width * Height != 960`.  The attribute data, however, will still be 64 bytes
long.

The only difference with sprite data is that there is one byte of attribute
data for each tile.  Sprites with 6 tiles will have 6 bytes of attribute data.

An examle of a segment with three metasprites follows.

### Complete data

    00000000  79 00 2b 00 03 03 05 0f  00 40 38 05 06 1e 00 78  |y.+......@8....x|
    00000010  28 03 02 06 00 60 40 00  01 02 03 04 05 06 07 08  |(....`@.........|
    00000020  09 0a 0b 0c 0d 0e 00 02  02 02 02 03 03 03 03 02  |................|
    00000030  03 02 00 00 02 0c 0f 10  11 12 00 13 14 15 16 17  |................|
    00000040  18 19 1a 1b 1c 1d 1e 1f  20 21 22 23 24 25 26 27  |........ !"#$%&'|
    00000050  28 29 2a 00 00 00 00 00  00 00 00 00 00 00 00 00  |()*.............|
    00000060  00 00 01 00 00 00 00 01  01 01 01 01 01 01 01 01  |................|
    00000070  00 2b 2c 2d 2e 2f 30 00  00 00 00 00 00 3d 01 17  |.+,-./0......=..|
    00000080  37 3d 01 26 37 3d 01 20  37 3d 01 20 02 00 00 00  |7=.&7=. 7=. ....|
    00000090  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
    000000a0  00 00 00 00 00 00 00 00                           |........|

### Segmented & Labeled

    Main header: Palettes at $007D ($0079 + 4), 3 images.
    00000000  79 00 2b 00 03 .. .. ..  .. .. .. .. .. .. .. ..  |y.+.............|

    Image #1 Header: 3x5 tiles, 15 attribute bytes, at coord (64, 48)
    00000000  .. .. .. .. .. 03 05 0f  00 40 38 .. .. .. .. ..  |.........@8.....|

    Image #2 Header: 5x6 tiles, 30 attribute bytes, at coord (120, 40)
    00000000  .. .. .. .. .. .. .. ..  .. .. .. 05 06 1e 00 78  |...............x|
    00000010  28 .. .. .. .. .. .. ..  .. .. .. .. .. .. .. ..  |(...............|

    Image #3 Header: 3x2 tiles, 6 attribute bytes, at coord (96, 64)
    00000010  .. 03 02 06 00 60 40 ..  .. .. .. .. .. .. .. ..  |.....`..........|

    Image #1 Tile Data
    00000010  .. .. .. .. .. .. .. 00  01 02 03 04 05 06 07 08  |......@.........|
    00000020  09 0a 0b 0c 0d 0e .. ..  .. .. .. .. .. .. .. ..  |................|

    Image #1 Attribute Data
    00000020  .. .. .. .. .. .. 00 02  02 02 02 03 03 03 03 02  |................|
    00000030  03 02 00 00 02 .. .. ..  .. .. .. .. .. .. .. ..  |................|

    Image #2 Tile Data
    00000030  .. .. .. .. .. 0c 0f 10  11 12 00 13 14 15 16 17  |................|
    00000040  18 19 1a 1b 1c 1d 1e 1f  20 21 22 23 24 25 26 27  |........ !"#$%&'|
    00000050  28 29 2a .. .. .. .. ..  .. .. .. .. .. .. .. ..  |()*.............|

    Image #2 Attribute Data
    00000050  .. .. .. 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
    00000060  00 00 01 00 00 00 00 01  01 01 01 01 01 01 01 01  |................|
    00000070  00 .. .. .. .. .. .. ..  .. .. .. .. .. .. .. ..  |................|

    Image #3 Tile Data
    00000070  .. 2b 2c 2d 2e 2f 30 ..  .. .. .. .. .. .. .. ..  |.+,-./0.........|

    Image #3 Attribute Data
    00000070  .. .. .. .. .. .. .. 00  00 00 00 00 00 .. .. ..  |................|

    Palette Data
    00000070  .. .. .. .. .. .. .. ..  .. .. .. .. .. 3d 01 17  |.............=..|
    00000080  37 3d 01 26 37 3d 01 20  37 3d 01 20 02 .. .. ..  |7=.&7=. 7=. ....|

    Padding
    00000080  .. .. .. .. .. .. .. ..  .. .. .. .. .. 00 00 00  |................|
    00000090  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
    000000a0  00 00 00 00 00 00 00 00                           |........|
