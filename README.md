# vi-minus-minus

It's like Vi/Vim, but worse.

## Installation

First, clone the repository:

```bash
git clone https://github.com/alex-laycalvert/vimm

cd vimm
```

Then, build the project:

```bash
go build
```


Prosper:

```bash
./vimm
```

## Usage

### Movements

Exiting:

- `<Ctrl-C>`: exit while in normal mode

Normal Vim keybindings:

- `h`: move left
- `j`: move down
- `k`: move up
- `l`: move right

Going to the top and bottom of a file:

- `g`: goto top
- `G`: goto bottom

### Editing

- `i`: enter insert mode
- `I`: enter insert mode and goto beginning of line
- `a`: enter insert mode and move right one
- `A`: enter insert mode and goto end of line
- `d`: delete current line and put in clipboard
- `p`: paste contents of clipboard at current col
- `o`: insert new line below current line
- `O`: insert new line above current line
- `<ESC>`: enter normal mode
