## Todo
### Editor:
    [x] Add basic word motions like: w, W, e, E
    [x] Add basic line motions like: $, ^, gg, G
    [x] Add commands for normal mode: delete, change, yank
    [x] Add commands for visual mode: delete, change, yank
    [x] Add history: undo, redo
    [ ] Add search in buffer
    [ ] Add multi-cursor support (maybe)
    [ ] Make keybindings configurable
    [ ] Fix inconsistency between word motions and cursor movement

### Problems:
    - Cursor disappears at end of line.
    - cw deletes current word + space + next word's first char.
        Got: <h>ello word" -> "<o>rld"
        Want: <h>ello word" -> "<> world"
    - cb acts strangely
        - Got: "hell<o> word" -> "hell<e>llo world"
        - Want: "hell<o> word" -> "<>o world"
    - dw current word + space + next word's first char.
        Got: "<h>ello word" -> "<o>rld"
        Want:  "<h>ello word" -> "<>o world"
    - db acts strangely
        - Got: "hell<o> word" -> "hell<e>llo world"
        - Want: "hell<o> word" -> "hell<e>llo world"
    - ce works as expected
    - cE works as expected
    - de works as expected
    - dE works as expected


