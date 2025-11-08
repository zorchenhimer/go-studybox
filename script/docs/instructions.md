## 0x80 Play Beep

Stack Arguments:  0
Inline Arguments: 0

Vars used:

    Byte_0493

Play's an audible beep on the Square 1 channel.

## 0x81 Halt

Stack Arguments:  0
Inline Arguments: 0

Vars used:
N/A

Infinite loop that does not return.

## 0x82 Tape NMI Shenanigans

Stack Arguments:  0
Inline Arguments: 0

Vars:

    Byte_E0_TapeCtrl_Cache
    Byte_EE
    Byte_F2

    Byte_0740
    Byte_07EF
    Byte_07F3

JSRs:

    L2706_SetupNMI_ED00_LongJump
    L2721_TurnOnNMI_LongJump
    L2724_TurnOffNMI_LongJump
    L2742


## 0x83 Tape Wait

Stack Arguments:  0
Inline Arguments: 0

Vars:

    Byte_0740

JMPs to `L1329_WaitOn_EE`

## 0x84 Jump

Stack Arguments:  0
Inline Arguments: 1 Word

Vars:

    Code_Pointer
    Argument_A

Updates the script pointer to the inline address and continues script execution
from the new address.

## 0x85 Call

Stack Arguments:  0
Inline Arguments: 1 Word

Vars:

    Code_Pointer
    Argument_A
    Stack_Pointer

Pushes return address to the stack and performs a Jump to the inline script
address.

## 0x86 Return

Stack Arguments:  0 (1 Word, implied)
Inline Arguments: 0

Vars:

    Stack_Pointer
    Code_Pointer

Remove a script address from the stack and update the `Code_Pointer` to it
before continuing execution.

## 0x87 Loop

Stack Arguments:  0 (4 implied)
Inline Arguments: 0

Args manually pulled from stack:

- Limit
- Increment
- LoopVar
- LoopEntry

ArgA `Stack_Pointer-6`
ArgB `(ArgD)`
ArgC `Stack_Pointer-2`
ArgD `Stack_Pointer-4`
ArgE `ArgA+1`

Vars:

    Argument_A
    Argument_B
    Argument_C
    Argument_D
    Argument_E
    Stack_Pointer

JSRs:

    Handler_CB_Sum
    Handler_C7_LessThan

JMPs to `L49CD_LessThan`

## 0x88 Play Sound

Stack Arguments:  32 bytes (string copied to `$0700`)
Inline Arguments: 0

Plays a short SFX defined by a string.

Vars:

    Pointer_A0
    Pointer_A2

    Byte_0494
    Byte_0495
    Byte_0496
    Byte_0497
    Byte_0498
    Byte_0499
    Byte_04AC_AudioState
    Byte_04AD

Vars inside `L5C18_CopyPtrA0PtrA2`:

    Byte_0490_PtrA0Len

Vars inside `L5A6F_DecodeAudioString_EntryPoint`:

    Pointer_A2
    Byte_04AC_AudioState

    Word_049A+0 (current channel ID)
    Word_049A+1 (current channel mask)

    Byte_0497 Byte_0498 Byte_0499
        enable byte for each channel.  RTS if != 1 on entry
        written to with value of Byte_04B3

    Byte_0494 Byte_0495 Byte_0496
        Pointer_A2 low for channel.  looks like it's updated after reading a
        string.

    Byte_049F Added to note lookup index; used with O stuff?
    Byte_04B1 octave?  value from Table_04A6
    Byte_04B2 T value
    Byte_04B3 = Byte_04B1 * Byte_04B2

    Table_049C  Y#
    Table_04A0  V##
    Table_04A3  M#
    Table_04A6  ?? note related.  octave?  value after letter
    Table_04A9  O#
    Table_04AE  T#

    Table_B391  G??

JSRs:

    L5C08_DataLen_A0
    L5C18_CopyPtrA0PtrA2
    L5CC5_LongDelay
    L5A6F_DecodeAudioString_EntryPoint
        L5A8F_DecodeAudioString (if audio is turned on)
            L5C42 (for weird chars?)

Hard-coded addresses:

    $0700 (Pointer_A0)
    $0420 (Pointer_A2)
    $0441 (Pointer_A2)
    $0462 (Pointer_A2)

Copies data currently at `$0700` to three locations (`$0420`, `$0441`, `$0462`)

### String format

Three channels are encoded in this string, separated by a colon (`:`).  The
string consists of letter and number pairs.  The order of the individual pairs
in the overal string don't seem to matter.

M#  Loop & Constant volume
    number is 0 or 1.

V## Volume
    number between 0 and 15, inclusive

Y#  Duty
    number between 0 and 3, inclusive

T#  ?? $4001/$4006?
    number between 1 and 9, inclusive (verify this)

O#  Note related (timer low/high stuff)
    number between 1 and 6, inclusive

A#-G# Notes and octaves?
    number between 0 and 10, inclusive

## 0x89

Stack Arguments:  3
Inline Arguments: 0

## 0x8A Pop String to Address

Stack Arguments:  0
Inline Arguments: 1 Word

Removes 32 bytes from the stack and writes them starting to the inline address.

## 0x8B

## 0x97

Arguments: 2

ArgA
ArgB

Vars:

    Byte_0740
    Byte_44FE
    Byte_4598 = Argument_B+0
    Byte_44FD

    Argument_A
    Array_44FB+2 ($44FD)

JSRs:

    L5592_CheckForZero

Conditionally sets up NMI stuff depending on byte $0740
If ArgA >= 3, setup some more stuff

## 0xA9

Arguments: 1

ArgA    

Some sort of nametable restore operation?

Writes smaller updates to the nametable for animation purposes.

## 0x9D Something Tape (draw screen?)

Arguments: 2

ArgA -> X
AgrB

Vars:

    Byte_44FE
    Array_44F8, X
    Array_44CF, X
    Array_457A, X

If ArgB != 0

    ldx ArgA

    lda #1
    sta Array_44F8, X
    sta Array_44CF, X

    lda #0
    Array_457A, X

else, check for zero from $44F9 through $44FC.
    if non-zero, store #1 into `Byte_FF4E` and return.
    else:

        ldx ArgA
        lda #1
        sta Array_44F8, X
        lda #0
        sta Byte_44F8

Then wait for `Array_44F8, X` to become zero.  We're waiting for the IRQ to
finish writing data from the tape to RAM.

After this, jump into opcode `0x9E` after the argument parsing to draw a
screen.

## 0x9E Draw And Show Screen

Arguments: 2

ArgA
ArgB

If ArgB == 0, clear out a bunch of arguments and call `Handler_DD`.  Clears
`Byte_44FE` afterwards and returns.

If ArgB != 0, make sure data has been loaded off the tape and is ready to draw.
This is done by checking `Word_FF49, X` for a value of 1.  If 1, increment it
and store it in `Byte_44FE` before returning.  Continue otherwise.

### Drawing

    lda #0
    sta Byte_0750
    sta Byte_4579
    sta Byte_457A ("Array_457A")

    lda ArgA
    sta Byte_4C


### Data layout

Data in RAM, starting at CPU address $5000

    $5000 Word offset to palette data.  Added to #$5004.
    $5002 Word
    $5004 Byte loop counter apparently?? adds #6 to #$5004 this many times in a
               pointer.  additional header data?

    // generic image header data
    $5005 Byte Width (data row len (tile data len))
    $5006 Byte Height (row count (title count))
                Data length = Width*Height (stored in Word_61 and Word_6AFE)

    $5007 Word data length? offset offset to attr data?

    $5008 Byte X/Y coords
    $5009 Byte X/Y coords


    $500B Data start

## 0xBB Push String / Push Data

Stack Arguments:  0
Inline Arguments: 32 bytes max (NULL terminated)

Push a NULL terminated string to the stack.  This opperation will increment the
stack pointer by 32 bytes, always.  Even if more than 32 bytes are pushed to
the stack.  The code pointer is properly incremented, so long as there is a
NULL byte within 255 bytes of the start of data.

## 0xC1 Jump Switch

Stack Arguments:  1
Inline Arguments: 3+

ArgA    

## 0xC5 Equal

If ArgA == ArgB
    push 1 to stack
If ArgA != ArgB
    push 0 to stack

## 0xC6 Not Equal

Stack Arguments: 2
Stack Result:    1

If ArgA != ArgB
    push 1 to stack
If ArgA == ArgB
    push 0 to stack

## 0xC7 Less Than

If ArgA < ArgB
    push 1 to stack
If ArgA >= ArgB
    push 0 to stack

## 0xC8 Less Than or Equal

If ArgA <= ArgB
    push 1 to stack
If ArgA > ArgB
    push 0 to stack

## 0xC9 Greater Than

If ArgA > ArgB
    push 1 to stack
If ArgA <= ArgB
    push 0 to stack

## 0xCA Greater Than or Equal To

If ArgA >= ArgB
    push 1 to stack
If ArgA < ArgB
    push 0 to stack

## 0xCB Add

ArgA = ArgA + ArgB

## 0xCC Subtract

ArgA = ArgA - ArgB

## 0xCD Multiply

ArgA-ArgB = ArgA * ArgB

## 0xCE Signed Divide

ArgA = ArgA / ArgB

## 0xCF Negate

ArgA = 0 - ArgA

## 0xD1 Controller Stuff

Stack Arguments:  2
Inline Arguments: 0

Vars:

    Argument_C+0 = Argument_A+0
    Argument_A = Word_B1 + Argument_B
    Argument_B


# 0xD4 Set Cursor Location

Stack Arguments: 3

Some sort of setup for 0xFE Draw Rom Character.

    Byte_0606 = ArgA
    Byte_0604 = ArgB
    Byte_0605 = ArgC

## 0xE0 Modulo

ArgA = ArgA % ArgB

# 0xE7 Draw Metasprite

Stack Arguments: 7

ArgA    sprite ID? (read as byte, but high is used as temp) Some sort of list size or count (header count?)
        This is used as a table lookup into a metasprite pointer table at $6980
ArgB    (byte) X Coord
ArgC    (byte) Y Coord
ArgD    (byte?) Palette override.  If positive, bottom two bits are used directly for palette index.
ArgE    (byte) switch of some sort.  sets X to $00 if zero, $20 if not zero
ArgF    (byte) Sprite flip.  Uses two middle bits of value (`%0001_1000`)
ArgG    (byte) Extra args on HW stack??

Header data for metasprites.  This location is pointed to by a table at $6980.

Width
Height
Count
Palette??

There is a table at $0140 that keeps track of sprite allocations.  Each byte
corresponds to a hardware sprite and the value corresponds to a metasprite that
that hardware sprite is a part of.

# 0xFE Draw Rom Character

Stack Arguments:  4
Inline Arguments: 1 Word

Observed drawing a 1bpp kanji character taken from the SBX ROM charset.
Inline word seems to be an ID or index for the character to draw.

