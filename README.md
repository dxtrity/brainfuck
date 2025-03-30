# Cerebrum Fornicate

An **overengineered** [Brainfuck](https://en.wikipedia.org/wiki/Brainfuck) interpreter that can do everything every other brainfuck interpreters can but also a little more.

## Core Features
- [x] Full Brainfuck Interpreter
- [x] Interpret by File
- [x] Interpret by Cmd Args

### Commands
`>` *increment* the data **pointer** <br>
`<` *decrement* the data **pointer** <br>
`+` *increment* the **byte** at the data pointer <br>
`-` *decrement* the **byte** at the data pointer <br>
`.` *output* the **byte** at the data pointer <br>
`,` accept one **byte** of *input* <br>
`[` *jump forward* past the **matching *]*** if the byte at the pointer is **0** <br>
`]` *jump back* to the **matching *[*** if the byte at the pointer is **nonzero** <br>
`#` print **debug information** <br>
`@` print a **memory snapshot** <br>

### Extra Features
This is an **overengineered** interpreter so it has things you don't need but are cool enough to make me spend 5h on the project.

#### Debug Information
<div align="center">
    <img src="./static/debug.png"/>
</div>

At anypoint you wish, add a **#** in the middle of your script and the interpreter will spew out data at you.

#### Memory Snapshots
<div align="center">
    <img src="./static/snapshot.png"/>
</div>

Like the debug, except putting **@** will spew out a snapshot of the current memory address at your pointer and ones around it in a table.

You can customise how much of the memory you see with the `-snapx` and `-snapy` flags. Each representing the amount of addresses above and below the pointer address.

You can also use the `-memory` flag to set the exact size of the memory array by bytes.

#### Images and Colors
<div align="center">
    <img src="./static/creepa.png"/>
</div>

The interpreter can spew out colors to the terminal. <br>
Depending on the number in the current address the interpreter will spew out a 2 character wide square.

`>32`:   Black <br>
`>64`:   Red <br>
`>96`:   Green <br>
`>128`:  Yellow <br>
`>160`:  Blue <br>
`>192`:  Magenta <br>
`>224`:  Cyan <br>
default: White <br>

The size at which the image wraps can be customised with `-w` flag. It's 16 by default