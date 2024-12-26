# Lets play Rummikub

## Rules

The rule set is defined [here](https://rummikub.com/wp-content/uploads/2019/12/2600-English-1.pdf).

## Version 1

When implementing the game Rummikub, my first thought was to take the rules and
playable moves and implement them directly in code. Starting at first with the
tiles, then moving to sets, then the overall gameplay loop.

I knew I wanted to have this be multiplayer, hence why I also created a server
package. But I wanted the actual game implementation details abstracted away.
Which resulted in me created a model package to retain the game specific objects
and behaviors. Maybe model should've been renamed to game or Rummikub but I didn't
want to let semantics get in the way of development and if I come up with a
better name I can always refactor later.

Creating the tiles or pieces first, I knew I wanted to interact with the underlying
struct through an interface. So the Piece interface was create to manipulate the
piece struct. Note the case sensitive naming which Golang utilizes to export
or unexport types, vars, consts, and functions.

The Piece interface has some getter functions to check if a Piece is a joker,
matches the color of a given piece, matches the value of a given piece, and to
return the uint8 value of the piece. The reason why I chose uint8 was because
none of the values will ever be negative and the bitsize of 8 was because none
of the pieces will ever exceed a value of 13.

Colors are also a uint8 value but they are fixed to 1, 2, 3, and 4 representing
the colors black, blue, red, and green respectively. This is done through an iota
with const variables to be used like an enum. Which is also why I've prefixed
the variable names with "Color" for helping me in autocompleting.

Not exported within the piece.go file is the isValidPiece function which returns
a boolean if the values are less than 14 and if the Colors are either black,
blue, red, or green.

Moving on to the Set interface, I started by unthinkingly implementing the allowed
set manipulations, omitting the joker moves which I will treat as a special case.
From the rules I extracted 4 main operations: inserting into a set, removing from
a set to create another, splitting a set by inserting a piece, and combining
pieces together by taking pieces from the set or player rack.

Additionally I included validation functions to check if a set is a group or run.
I've included these into the interface although in hindsight I should also treat
these the same way I treated the isValidPiece function for the pieces. I did the
same for isValidSet and created it to validate if a set is at least 3 pieces and
either a valid group or a valid run.

To manipulate the slice of Pieces in the set easier I created functions to clone
the tiles of a set, remove a piece by index, and insert a piece at index. Which
lead me to implement the insert piece and remove piece. The split piece operation
came just as quickly but I encountered an issue with the combine operation. This
was unique because I wasn't necessarily bound to pieces from the player's rack.
These combinations can be made by taking pieces from other sets on the board, so
I needed a way to pass in the Pieces I wanted either from a set or rack.

This lead me to come up with a variadic parameter where I validate if I've passed
in a set and index pair or just a plain piece. The heaviest part of the combine
operation was to validate the parameters, it was straightforward then to build
the set by taking these Pieces in order and checking if the sets I would be
making is valid and if the removed pieces didn't make existing sets invalid.

After all the operations were completed, I was able to begin building the game
state by initializing all the pieces, shuffling them, and then dealing them. I
needed to create a way for the user to play the game so I made a simple command
line interface where the use can tell which operation they wanted to invoke and
passed in parameters to these operations. Overall this was a good first attempt
but I have a few new more ideas in mind to make the code better. These will be
expanded upon in the next version.

## Version 2

In this next version, I realized that the operations for insert, splitting,
and combining pieces all had only one correct way of doing it. For instance, you
are allowed to insert into a group only if the piece is matching the value and
fulfills the missing piece. Inserting into a run is only possible at the end or
beginning. Pieces can only be removed if the remaining tiles are a valid set.
Splitting can only happen when you introduce an additional piece and both sets
from the split are valid sets.

Thus with this in mind, I went ahead and simplified the insert, remove, split,
and combine operations to take Pieces as parameters. Performing the validation
on the set after the operation, returning an error otherwise. In the cases of
splitting and combining where a new Set is created, we also return the created
Set to add to the board.

This made the player commands a lot more simpler where before they needed to
provide positional indexes to the command. Now the player can use a generic
piece selection method and pass the selected pieces into the command.

In the terminal interface, performing an operation such as insert is simply:
select the set from the board by index, select the piece from the rack by index.
The removal from the set and rack is hidden from the player. Splitting a set is
also the same process of selecting the run to split and selecting the piece from
the rack. Combining pieces into a set is also as easy as selecting pieces from
sets via the s#,# notation or selecting a piece from the rack via the r# notation.

Although this limits the player into performing strictly valid operations. Consider
the multiple split example from the rule set. The player can do this by combining
pieces from the sets together but when removing the pieces from the sets, the
intermediary sets are considered invalid.

In this example, a set denotes pieces encapsulated by square brackets. A piece
value and it's color are denoted by the number followed by W, X, Y, Z as color
(e.g. 1W, 9Y, 3Z, etc.). Pieces in the player's rack are surrounded by parentheses.

This is the start of a player's turn:
Player: (10Y, 5Z)
Board: [5W, 6W, 7W] [5X, 6X, 7X] [5Y, 6Y, 7Y, 8Y, 9Y]

But at the end they want the state of the game to be:
Player: ()
Board: [5W, 5X, 5Y, 5Z] [6W, 6X, 6Y] [7W, 7X, 7Y] [8Y, 9Y, 10Y]

In reaching this end state, a player might create this intermediary state:
Player: (10Y, 5Z)
Board: [5W, 5X, 5Y] [6W, 7W] [6X, 7X] [6Y, 7Y, 8Y, 9Y]

Because we are performing set validations during each operation, there is only
one way to reach this end state. Thus limiting the creativity and overall fun of
the game.

Player inserts 10Y into [5Y, 6Y, 7Y, 8Y, 9Y] resulting in [5Y, 6Y, 7Y, 8Y, 9Y, 10Y]
Player combines 5W, 5X, 5Y resulting in [5W, 5X, 5Y]
Player inserts 5Z into [5W, 5X, 5Y] resulting in [5W, 5X, 5Y, 5Z]
Player combines 6W, 6X, 6Y resulting in [6W, 6X, 6Y]
Player combines 7W, 7X, 7Y resulting in [7W, 7X, 7Y]
Player replaces the set [5W, 6W, 7W] with [5W, 5X, 5Y, 5Z]
Player replaces the set [5X, 6X, 7X] with [6W, 6X, 6Y]
Player replaces the set [5Y, 6Y, 7Y, 8Y, 9Y] with [7W, 7X, 7Y]
Player adds the set [8Y, 9Y, 10Y]

Although this order of operations keeps all the sets in a valid state, it is not
intuitive or human friendly as the juggling of valid set operations is required.

In the next version I want to address this and allow the player to perform these
invalid moves, undo/redo them, and apply them once the player is satisfied.
