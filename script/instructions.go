package script

import (
)

var InstrMap map[byte]*Instruction

func init() {
	InstrMap = make(map[byte]*Instruction)
	for _, i := range Instructions {
		InstrMap[i.Opcode] = i
	}
}

var Instructions []*Instruction = []*Instruction{
	&Instruction{ 0x80, 0, 0, 0,  "play_beep"},
	&Instruction{ 0x81, 0, 0, 0,  "halt"},
	&Instruction{ 0x82, 0, 0, 0,  "tape_nmi_shenanigans"},
	&Instruction{ 0x83, 0, 0, 0,  "tape_wait"},
	&Instruction{ 0x84, 0, 2, 0,  "jump_abs"},
	&Instruction{ 0x85, 0, 2, 0,  "call_abs"},
	&Instruction{ 0x86, 0, 0, 0,  "return"},
	&Instruction{ 0x87, 0, 0, 0,  "loop"},
	&Instruction{ 0x88, 0, 0, 0,  "play_sound"},
	&Instruction{ 0x89, 3, 0, 0,  ""},
	&Instruction{ 0x8A, 0, 2, 0,  "pop_string_to_addr"},
	&Instruction{ 0x8B, 1, 0, 0,  ""},
	&Instruction{ 0x8C, 0, 0, 1,  "string_length"},
	&Instruction{ 0x8D, 0, 0, 1,  "string_to_int"},
	&Instruction{ 0x8E, 0, 0, 16, "string_concat"},
	&Instruction{ 0x8F, 0, 0, 1,  "strings_equal"},

	&Instruction{ 0x90, 0, 0, 1,  "strings_not_equal"},
	&Instruction{ 0x91, 0, 0, 1,  "string_less_than"},
	&Instruction{ 0x92, 0, 0, 1,  "string_less_than_equal"},
	&Instruction{ 0x93, 0, 0, 1,  "string_greater_than_equal"},
	&Instruction{ 0x94, 0, 0, 1,  "string_greater_than"},
	&Instruction{ 0x95, 1, 0, 0,  ""},
	&Instruction{ 0x96, 0, 2, 0,  "set_word_4E"},
	&Instruction{ 0x97, 2, 0, 0,  ""},
	&Instruction{ 0x98, 1, 0, 0,  ""},
	&Instruction{ 0x99, 1, 0, 0,  ""},
	&Instruction{ 0x9A, 0, 0, 0,  ""},
	&Instruction{ 0x9B, 0, 0, 0,  "halt"},
	&Instruction{ 0x9C, 0, 0, 0,  "toggle_44FE"},
	&Instruction{ 0x9D, 2, 0, 0,  "something_tape"},
	&Instruction{ 0x9E, 2, 0, 0,  ""},
	&Instruction{ 0x9F, 6, 0, 0,  ""},

	&Instruction{ 0xA0, 2, 0, 1,  ""},
	&Instruction{ 0xA1, 1, 0, 0,  ""},
	&Instruction{ 0xA2, 1, 0, 0,  "buffer_palette"},
	&Instruction{ 0xA3, 1, 0, 0,  ""},
	&Instruction{ 0xA4, 3, 0, 0,  ""},
	&Instruction{ 0xA5, 1, 0, 0,  "set_470A"},
	&Instruction{ 0xA6, 1, 0, 0,  "set_470B"},
	&Instruction{ 0xA7, 0, 0, 0,  "call_asm"}, // built-in ACE, lmao
	&Instruction{ 0xA8, 5, 0, 0,  ""},
	&Instruction{ 0xA9, 1, 0, 0,  ""},
	&Instruction{ 0xAA, 1, 0, 0,  ""},
	&Instruction{ 0xAB, 1, 0, 0,  "long_call"},
	&Instruction{ 0xAC, 0, 0, 0,  "long_return"},
	&Instruction{ 0xAD, 1, 0, 1,  "absolute"},
	&Instruction{ 0xAE, 1, 0, 1,  "compare"},
	&Instruction{ 0xAF, 0, 0, 1,  ""},

	&Instruction{ 0xB0, 1, 0, 16, ""},
	&Instruction{ 0xB1, 1, 0, 16, "to_hex_string"},
	&Instruction{ 0xB2, 0, 0, 1,  ""},
	&Instruction{ 0xB3, 7, 0, 0,  ""}, // possible 16-bit inline?
	&Instruction{ 0xB4, 0, 0, 0,  ""},
	&Instruction{ 0xB5, 0, 0, 0,  ""},
	&Instruction{ 0xB6, 0, 0, 0,  ""},
	&Instruction{ 0xB7, 0, 2, 0,  "deref_ptr"},
	&Instruction{ 0xB8, 0, 2, 0,  "push_word"},
	&Instruction{ 0xB9, 0, 2, 0,  "push_word_indexed"},
	&Instruction{ 0xBA, 0, 2, 0,  "push"},
	&Instruction{ 0xBB, 0, -1, 0, "push_data"},
	&Instruction{ 0xBC, 0, 2, 0,  "push_string_from_table"},
	&Instruction{ 0xBD, 0, 2, 0,  "pop"},
	&Instruction{ 0xBE, 0, 2, 0,  "write_to_table"},
	&Instruction{ 0xBF, 0, 2, 0,  "jump_not_zero"},

	&Instruction{ 0xC0, 1, 2, 0,  "jump_zero"},
	&Instruction{ 0xC1, 1, -2, 0, "jump_switch"},
	&Instruction{ 0xC2, 1, 0, 1,  "equals_zero"},
	&Instruction{ 0xC3, 2, 0, 1,  "and_a_b"},
	&Instruction{ 0xC4, 2, 0, 1,  "or_a_b"},
	&Instruction{ 0xC5, 2, 0, 1,  "equal"},
	&Instruction{ 0xC6, 2, 0, 1,  "not_equal"},
	&Instruction{ 0xC7, 2, 0, 1,  "less_than"},
	&Instruction{ 0xC8, 2, 0, 1,  "less_than_equal"},
	&Instruction{ 0xC9, 2, 0, 1,  "greater_than"},
	&Instruction{ 0xCA, 2, 0, 1,  "greater_than_equal"},
	&Instruction{ 0xCB, 2, 0, 1,  "sum"},
	&Instruction{ 0xCC, 2, 0, 1,  "subtract"},
	&Instruction{ 0xCD, 2, 0, 1,  "multiply"},
	&Instruction{ 0xCE, 2, 0, 1,  "signed_divide"},
	&Instruction{ 0xCF, 1, 0, 1,  "negate"},

	&Instruction{ 0xD0, 1, 0, 1,  "modulus"},
	&Instruction{ 0xD1, 2, 0, 1,  "expansion_controller"},
	&Instruction{ 0xD2, 2, 0, 1,  ""},
	&Instruction{ 0xD3, 2, 0, 16, ""},
	&Instruction{ 0xD4, 3, 0, 0,  ""},
	&Instruction{ 0xD5, 1, 0, 0,  "wait_for_tape"},
	&Instruction{ 0xD6, 1, 0, 16, "truncate_string"},
	&Instruction{ 0xD7, 1, 0, 16, "trim_string"},
	&Instruction{ 0xD8, 1, 0, 16, "trim_string_start"},
	&Instruction{ 0xD9, 2, 0, 16, "trim_string_start"},
	&Instruction{ 0xDA, 1, 0, 16, "to_int_string"},
	&Instruction{ 0xDB, 3, 0, 0,  ""},
	&Instruction{ 0xDC, 5, 0, 0,  ""},
	&Instruction{ 0xDD, 5, 0, 0,  ""},
	&Instruction{ 0xDE, 3, 0, 0,  ""},
	&Instruction{ 0xDF, 3, 0, 0,  ""},

	&Instruction{ 0xE0, 2, 0, 1,  "signed_divide"},
	&Instruction{ 0xE1, 4, 0, 0,  ""},
	&Instruction{ 0xE2, 7, 0, 0,  "setup_sprite"},
	&Instruction{ 0xE3, 1, 0, 1,  "get_byte_at_arg_a"},
	&Instruction{ 0xE4, 2, 0, 0,  "swap_ram_bank"},
	&Instruction{ 0xE5, 1, 0, 0,  "disable_sprite"},
	&Instruction{ 0xE6, 1, 0, 0,  "tape_nmi_setup"},
	&Instruction{ 0xE7, 7, 0, 0,  ""},
	&Instruction{ 0xE8, 1, 0, 0,  "setup_tape_nmi"},
	&Instruction{ 0xE9, 0, 1, 0,  "setup_loop"},
	&Instruction{ 0xEA, 0, 0, 0,  "string_write_to_table"},
	&Instruction{ 0xEB, 4, 0, 0,  ""},
	&Instruction{ 0xEC, 2, 0, 0,  "scroll"},
	&Instruction{ 0xED, 1, 0, 0,  "disable_sprites"},
	&Instruction{ 0xEE, 1, -3, 0,  "call_switch"},
	&Instruction{ 0xEF, 6, 0, 0,  ""},

	&Instruction{ 0xF0, 0, 0, 0,  "disable_sprites"},
	&Instruction{ 0xF1, 4, 0, 0,  ""},
	&Instruction{ 0xF2, 0, 0, 0,  "halt"},
	&Instruction{ 0xF3, 0, 0, 0,  "halt"},
	&Instruction{ 0xF4, 0, 0, 16, "halt"},
	&Instruction{ 0xF5, 1, 0, 1,  "halt"},
	&Instruction{ 0xF6, 1, 0, 0,  "halt"},
	&Instruction{ 0xF7, 0, 0, 0,  "halt"},
	&Instruction{ 0xF8, 2, 0, 0,  "halt"},
	&Instruction{ 0xF9, 0, 0, 1,  ""},
	&Instruction{ 0xFA, 0, 0, 1,  ""},
	&Instruction{ 0xFB, 1, 0, 0,  "jump_arg_a"},
	&Instruction{ 0xFC, 2, 0, 1,  ""},
	&Instruction{ 0xFD, 0, 0, 16, "halt"},
	&Instruction{ 0xFE, 4, 0, 0,  ""},
	&Instruction{ 0xFF, 0, 0, 0,  "break_engine"}, // code handler is $FFFF
}

type Instruction struct {
	Opcode   byte
	ArgCount int // stack arguments
	OpCount  int // inline operands.  length in bytes.
				 // -1: nul-terminated
				 // -2: first byte is count, followed by that number of words
				 // -3: like -2, but with one additional word
	RetCount int // return count
	Name     string
}

func (i Instruction) String() string {
	if i.Name != "" {
		//return fmt.Sprintf("$%02X_%s", i.Opcode, i.Name)
		return i.Name
	}

	//return fmt.Sprintf("$%02X_unknown", i.Opcode)
	return "unknown"
}

