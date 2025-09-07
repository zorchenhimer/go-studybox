package script

import (
	"fmt"
)

var InstrMap map[byte]*Instruction

func init() {
	InstrMap = make(map[byte]*Instruction)
	for _, i := range Instructions {
		InstrMap[i.Opcode] = i
	}
}

var Instructions []*Instruction = []*Instruction{
	&Instruction{ 0x80, 0, 0, 0, false,  "play_beep"},
	&Instruction{ 0x81, 0, 0, 0, false,  "halt"},
	&Instruction{ 0x82, 0, 0, 0, false,  "tape_nmi_shenanigans"},
	&Instruction{ 0x83, 0, 0, 0, false,  "tape_wait"},

	// Jump to the inline word, in the VM
	&Instruction{ 0x84, 0, 2, 0, false,  "jump_abs"},

	// Call a routine at the inline word, in the VM
	&Instruction{ 0x85, 0, 2, 0, false,  "call_abs"},

	// Return from a previous call
	&Instruction{ 0x86, 0, 0, 0, false,  "return"},

	&Instruction{ 0x87, 0, 0, 0,  false, "loop"},
	&Instruction{ 0x88, 0, 0, 0,  false, "play_sound"},
	&Instruction{ 0x89, 3, 0, 0,  false, ""},
	&Instruction{ 0x8A, 0, 2, 0,  false, "pop_string_to_addr"},
	&Instruction{ 0x8B, 1, 0, 0,  false, ""},
	&Instruction{ 0x8C, 0, 0, 1,  false, "string_length"},
	&Instruction{ 0x8D, 0, 0, 1,  false, "string_to_int"},
	&Instruction{ 0x8E, 0, 0, 16, false, "string_concat"},
	&Instruction{ 0x8F, 0, 0, 1,  false,  "strings_equal"},

	&Instruction{ 0x90, 0, 0, 1,  false,  "strings_not_equal"},
	&Instruction{ 0x91, 0, 0, 1,  false,  "string_less_than"},
	&Instruction{ 0x92, 0, 0, 1,  false,  "string_less_than_equal"},
	&Instruction{ 0x93, 0, 0, 1,  false,  "string_greater_than_equal"},
	&Instruction{ 0x94, 0, 0, 1,  false,  "string_greater_than"},

	// Sets some tape NMI stuff if the byte at $0740 is not zero.
	// Will call 0x82 tape_nmi_shenanigans if $0740 != 0
	&Instruction{ 0x95, 1, 0, 0,  false,  "tape_nmi_shenigans_set"},

	&Instruction{ 0x96, 0, 2, 0,  true,   "set_word_4E"},
	&Instruction{ 0x97, 2, 0, 0,  false,  ""},
	&Instruction{ 0x98, 1, 0, 0,  false,  ""},
	&Instruction{ 0x99, 1, 0, 0,  false,  ""},
	&Instruction{ 0x9A, 0, 0, 0,  false,  ""},
	&Instruction{ 0x9B, 0, 0, 0,  false,  "halt"},
	&Instruction{ 0x9C, 0, 0, 0,  false,  "toggle_44FE"},
	&Instruction{ 0x9D, 2, 0, 0,  false,  "something_tape"},

	// Calls 0xEB draw_overlay.  Draws the whole screen from data previously
	// loaded from the tape.
	&Instruction{ 0x9E, 2, 0, 0,  false,  "draw_and_show_screen"},

	&Instruction{ 0x9F, 6, 0, 0,  false,  ""},

	&Instruction{ 0xA0, 2, 0, 1,  false,  ""},
	&Instruction{ 0xA1, 1, 0, 0,  false,  ""},
	&Instruction{ 0xA2, 1, 0, 0,  false,  "buffer_palette"},

	// Possibly a sprite setup routine.  loads up some CHR data and some palette
	// data.
	&Instruction{ 0xA3, 1, 0, 0,  false,  "sprite_setup"},
	&Instruction{ 0xA4, 3, 0, 0,  false,  ""},
	&Instruction{ 0xA5, 1, 0, 0,  false,  "set_470A"},
	&Instruction{ 0xA6, 1, 0, 0,  false,  "set_470B"},

	// jump to the inline address, in assembly, not in the VM
	// (built-in ACE, lmao)
	// Will not jump to anything at or above $8000 or below $5000.
	// Addresses in $5000-$5FFF use $470A as the bank ID
	// Addresses in $6000-$7FFF use $470B as the bank ID
	&Instruction{ 0xA7, 0, 0, 0,  false,   "call_asm"},

	&Instruction{ 0xA8, 5, 0, 0,  false,  ""},
	&Instruction{ 0xA9, 1, 0, 0,  false,  ""},
	&Instruction{ 0xAA, 1, 0, 0,  false,  ""},
	&Instruction{ 0xAB, 1, 0, 0,  false,  "long_call"},
	&Instruction{ 0xAC, 0, 0, 0,  false,  "long_return"},
	&Instruction{ 0xAD, 1, 0, 1,  false,  "absolute"},
	&Instruction{ 0xAE, 1, 0, 1,  false,  "compare"},
	&Instruction{ 0xAF, 0, 0, 1,  false,  ""},

	&Instruction{ 0xB0, 1, 0, 16, false, ""},
	&Instruction{ 0xB1, 1, 0, 16, false, "to_hex_string"},
	&Instruction{ 0xB2, 0, 0, 1,  false,  ""},
	&Instruction{ 0xB3, 7, 0, 0,  false,  ""}, // possible 16-bit inline?
	&Instruction{ 0xB4, 0, 0, 0,  false,  ""},
	&Instruction{ 0xB5, 0, 0, 0,  false,  ""},
	&Instruction{ 0xB6, 0, 0, 0,  false,  ""},

	// Uses the inline word as a pointer, and pushes the byte value at that
	// address to the stack.
	&Instruction{ 0xB7, 0, 2, 0,  false,  "deref_ptr_inline"},

	// Pushes the inline word to the stack
	&Instruction{ 0xB8, 0, 2, 0,  true,  "push_word"},
	&Instruction{ 0xB9, 0, 2, 0,  false, "push_word_indexed"},
	&Instruction{ 0xBA, 0, 2, 0,  false, "push"},
	&Instruction{ 0xBB, 0, -1, 0, false, "push_data"},
	&Instruction{ 0xBC, 0, 2, 0,  false, "push_string_from_table"},

	// Pops a byte off the stack and stores it at the inline address.
	&Instruction{ 0xBD, 0, 2, 0,  false,  "pop_into"},

	&Instruction{ 0xBE, 0, 2, 0,  false,  "write_to_table"},
	&Instruction{ 0xBF, 0, 2, 0,  false,  "jump_not_zero"},

	// One byte off stack; jumps to inline if byte is zero
	&Instruction{ 0xC0, 1, 2, 0,  false,  "jump_zero"},
	&Instruction{ 0xC1, 1, -2, 0, false,  "jump_switch"},
	&Instruction{ 0xC2, 1, 0, 1,  false,  "equals_zero"},
	&Instruction{ 0xC3, 2, 0, 1,  false,  "and_a_b"},
	&Instruction{ 0xC4, 2, 0, 1,  false,  "or_a_b"},
	&Instruction{ 0xC5, 2, 0, 1,  false,  "equal"},

	// Two bytes off stack; result pushed back; 1 if A == B, 0 if A != B
	&Instruction{ 0xC6, 2, 0, 1,  false,  "not_equal"},

	&Instruction{ 0xC7, 2, 0, 1,  false,  "less_than"},
	&Instruction{ 0xC8, 2, 0, 1,  false,  "less_than_equal"},
	&Instruction{ 0xC9, 2, 0, 1,  false,  "greater_than"},
	&Instruction{ 0xCA, 2, 0, 1,  false,  "greater_than_equal"},
	&Instruction{ 0xCB, 2, 0, 1,  false,  "sum"},
	&Instruction{ 0xCC, 2, 0, 1,  false,  "subtract"},
	&Instruction{ 0xCD, 2, 0, 1,  false,  "multiply"},
	&Instruction{ 0xCE, 2, 0, 1,  false,  "signed_divide"},
	&Instruction{ 0xCF, 1, 0, 1,  false,  "negate"},

	&Instruction{ 0xD0, 1, 0, 1,  false,  "modulus"},
	&Instruction{ 0xD1, 2, 0, 1,  false,  "expansion_controller"},
	&Instruction{ 0xD2, 2, 0, 1,  false,  ""},
	&Instruction{ 0xD3, 2, 0, 16, false,  ""},
	&Instruction{ 0xD4, 3, 0, 0,  false,  "set_cursor_location"},

	// Wait for ArgA itterations.  "itterations" is undefined as of now. (data from tape?)
	&Instruction{ 0xD5, 1, 0, 0,  false,  "wait_for_tape"},

	&Instruction{ 0xD6, 1, 0, 16, false,  "truncate_string"},
	&Instruction{ 0xD7, 1, 0, 16, false,  "trim_string"},
	&Instruction{ 0xD8, 1, 0, 16, false,  "trim_string_start"},
	&Instruction{ 0xD9, 2, 0, 16, false,  "trim_string_start"},
	&Instruction{ 0xDA, 1, 0, 16, false,  "to_int_string"},
	&Instruction{ 0xDB, 3, 0, 0,  false,  ""},
	&Instruction{ 0xDC, 5, 0, 0,  false,  ""},

	// ArgA, ArgB: X,Y of corner A
	// ArgC, ArgD: X,Y of corner B
	// ArgE: fill value.  This is an index into
	//       the table at $B451.
	// Fills a box with a tile
	&Instruction{ 0xDD, 5, 0, 0,  false,  "fill_box"},

	&Instruction{ 0xDE, 3, 0, 0,  false,  ""},
	&Instruction{ 0xDF, 3, 0, 0,  false,  ""},

	// Divide and return remainder
	&Instruction{ 0xE0, 2, 0, 1,  false,  "modulo"},

	&Instruction{ 0xE1, 4, 0, 0,  false,  ""},
	&Instruction{ 0xE2, 7, 0, 0,  false,  "setup_sprite"},

	// Pops a word off the stack, uses it as a pointer, and pushes the byte
	// value at that address to the stack.
	&Instruction{ 0xE3, 1, 0, 1,  false,  "deref_ptr_stack"},
	&Instruction{ 0xE4, 2, 0, 0,  false,  "swap_ram_bank"},
	&Instruction{ 0xE5, 1, 0, 0,  false,  "disable_sprite"},

	// Will call 0x82 tape_nmi_shenanigans if $0740 != 0
	&Instruction{ 0xE6, 1, 0, 0,  false,  "tape_nmi_setup"},

	&Instruction{ 0xE7, 7, 0, 0,  false,  "draw_metasprite"},
	&Instruction{ 0xE8, 1, 0, 0,  false,  "setup_tape_nmi"},
	&Instruction{ 0xE9, 0, 1, 0,  false,  "setup_loop"},
	&Instruction{ 0xEA, 0, 0, 0,  false,  "string_write_to_table"},

	// Reads and saves tiles from the PPU, then draws over them.
	// This is used to draw dialog boxes, so saving what it overwrites
	// so it can re-draw them later makes sense.
	// Not sure what the arguments actually mean.
	// ArgB and ArgC are probably coordinates.
	&Instruction{ 0xEB, 4, 0, 0,  false,  "draw_overlay"},

	&Instruction{ 0xEC, 2, 0, 0,  false,  "scroll"},
	&Instruction{ 0xED, 1, 0, 0,  false,  "disable_sprites"},

	&Instruction{ 0xEE, 1, -3, 0, false,  "call_switch"},
	&Instruction{ 0xEF, 6, 0, 0,  false,  ""},

	&Instruction{ 0xF0, 0, 0, 0,  false,  "disable_sprites"},
	&Instruction{ 0xF1, 4, 0, 0,  false,  ""},
	&Instruction{ 0xF2, 0, 0, 0,  false,  "halt"},
	&Instruction{ 0xF3, 0, 0, 0,  false,  "halt"},
	&Instruction{ 0xF4, 0, 0, 16, false,  "halt"},
	&Instruction{ 0xF5, 1, 0, 1,  false,  "halt"},
	&Instruction{ 0xF6, 1, 0, 0,  false,  "halt"},
	&Instruction{ 0xF7, 0, 0, 0,  false,  "halt"},
	&Instruction{ 0xF8, 2, 0, 0,  false,  "halt"},
	&Instruction{ 0xF9, 0, 0, 1,  false,  ""},
	&Instruction{ 0xFA, 0, 0, 1,  false,  ""},
	&Instruction{ 0xFB, 1, 0, 0,  false,  "jump_arg_a"},
	&Instruction{ 0xFC, 2, 0, 1,  false,  ""},
	&Instruction{ 0xFD, 0, 0, 16, false,  "halt"},
	&Instruction{ 0xFE, 4, 2, 0,  false,  "draw_rom_char"},

	// code handler is $FFFF
	&Instruction{ 0xFF, 0, 0, 0,  false,  "break_engine"},
}

type Instruction struct {
	Opcode    byte
	ArgCount  int  // stack arguments
	OpCount   int  // inline operands.  length in bytes.
				   // -1: nul-terminated
				   // -2: first byte is count, followed by that number of words
				   // -3: like -2, but with no default.  code continues after list on OOB
	RetCount  int  // return count
	InlineImmediate bool // don't turn the inline value into a variable
	Name      string
}

func (i Instruction) String() string {
	if i.Name != "" {
		//return fmt.Sprintf("$%02X_%s", i.Opcode, i.Name)
		return i.Name
	}

	return fmt.Sprintf("unknown_0x%02X", i.Opcode)
	//return "unknown"
}

