package main

// Contains a set of names associated with tags.
type NameDictionary map[string]Properties

func (d NameDictionary) Add(name string, p Properties) NameDictionary {
	if entry, ok := d[name]; ok {
		// Names are treated as being suitable for the specific name part
		// if they appear in that position in any of the input names.
		entry.First = entry.First || p.First
		entry.Given = entry.Given || p.Given
		entry.Last = entry.Last || p.Last
		entry.Nick = entry.Nick || p.Nick

		// Nicknames are treated specially, and will never be returned in a
		// generated name unless :nick is explicitly specified. Because of this
		// we need to track if the name was EVER input as a normal (rather than
		// nick) name, and if so, allow it to be used as such OR as a nickname.
		entry.NotNick = entry.NotNick || !entry.Nick

		d[name] = entry
	} else {
		d[name] = p
	}

	return d
}

// The properties for a single name in the dictionary.
type Properties struct {
	First, Given, Last, Nick, NotNick bool
	Tags []Tag
}
