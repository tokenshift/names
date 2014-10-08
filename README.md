# Names

Random name generator.

## Concepts

Name files (ending with the extension ".names") contain lists of names and tags
that are applied to those names. New names are generated from these based on a
_name format_ that specifies tags to match. Multiple tags can be combined using
the `+` and `|` symbols. There are also filters, like `:first` and `:last`,
that will match any name part that was the first or last component respectively
of a full name.

## Use

These examples make use of the `Steven King.names` name file.

```
names '(Las Vegas + Male):first | Dog' '[:given]' 'Boulder:last'
```

This will generate a name consisting of a first name taken from all first names
tagged 'Las Vegas' and 'Male', OR a name tagged 'Dog'; an optional middle name,
which will use any given name from the name file (50/50 chance of it being
included); and a last name from all last names tagged 'Boulder'.

Note that each of the desired name components is provided as a separate command
line parameter.

## Filters

Supported filters include:

* `:first` 
  Given a name "John Richard Smith", would return only "John".
* `:given` 
  Given a name "John Richard Smith", would return "John" or "Richard".
* `:last`
  Given a name "John Richard Smith", would return only "Smith".
* `:nick` 
  Returns any name components marked as nicknames (e.g. the "Billy" in 'William
  "Billy" Starkey').

## Name Files

* Initials (the "D." in "Charles D. Campion") are ignored.
* Hyphenated names count as two different names. For example, the input name
  "Peter Goldsmith-Redman" could produce the name "Peter", "Goldsmith", or
  "Redman" (both "Goldsmith" and "Redman" would count as :last names).
* Nicknames (surrounded by quotes, as in 'William "Billy" Starkey') are never
  used, unless the `:nick` filter is specified.
