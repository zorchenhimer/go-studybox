# extract-imgs

Extract images from an unpacked `.studybox` ROM file.  Requires the tile
packets as well as the CHR packets to build an image.  Handles both nametable
and sprite data.

# just-stats

Decodes scripts similar to `script-decode`, but does not save the output.
Instead, scripts are decoded in bulk and instruction usage is recorded to an
output file.

# sbutil

Pack and unpack `.studybox` ROM files.  Unpacking extracts all of the data from
the ROM into a subdirectory and writes a `.json` file with metadata.  Packing
does the reverse using the `.json` metadata file.

# sbx2wav

Encode a `.studybox` ROM into a WAV audio file.  Conversion is currently a bit
shaky and hasn't been confirmed to work on hardware.  Timing between the data
and the recorded audio could also use a little more work.

# script-decode

Decode script segments from an unpacked `.studybox` ROM file.  Labels and a
CDL file are supported.  Two modes are available: dumb and `--smart` decoding.

Dumb decoding is the default and will attempt to decode every byte in the file
as script data. This will try and decode variable data as script data which is
usually undesired.

`--smart` decoding starts at given entry points and decodes scripts by
following the logic of the script and recording branches as new entry points.
By default the only entry point is the top of the script (third byte in the
file), but additional entry points can be given in the CDL file.
